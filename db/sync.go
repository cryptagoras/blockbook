package db

import (
	"blockbook/bchain"
	"blockbook/common"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

// SyncWorker is handle to SyncWorker
type SyncWorker struct {
	db                     *RocksDB
	chain                  bchain.BlockChain
	syncWorkers, syncChunk int
	dryRun                 bool
	startHeight            uint32
	startHash              string
	chanOsSignal           chan os.Signal
	metrics                *common.Metrics
	is                     *common.InternalState
}

// NewSyncWorker creates new SyncWorker and returns its handle
func NewSyncWorker(db *RocksDB, chain bchain.BlockChain, syncWorkers, syncChunk int, minStartHeight int, dryRun bool, chanOsSignal chan os.Signal, metrics *common.Metrics, is *common.InternalState) (*SyncWorker, error) {
	if minStartHeight < 0 {
		minStartHeight = 0
	}
	return &SyncWorker{
		db:           db,
		chain:        chain,
		syncWorkers:  syncWorkers,
		syncChunk:    syncChunk,
		dryRun:       dryRun,
		startHeight:  uint32(minStartHeight),
		chanOsSignal: chanOsSignal,
		metrics:      metrics,
		is:           is,
	}, nil
}

var errSynced = errors.New("synced")

// ResyncIndex synchronizes index to the top of the blockchain
// onNewBlock is called when new block is connected, but not in initial parallel sync
func (w *SyncWorker) ResyncIndex(onNewBlock func(hash string)) error {
	start := time.Now()
	w.is.StartedSync()

	err := w.resyncIndex(onNewBlock)

	switch err {
	case nil:
		d := time.Since(start)
		glog.Info("resync: finished in ", d)
		w.metrics.IndexResyncDuration.Observe(float64(d) / 1e6) // in milliseconds
		w.metrics.IndexDBSize.Set(float64(w.db.DatabaseSizeOnDisk()))
		bh, _, err := w.db.GetBestBlock()
		if err == nil {
			w.is.FinishedSync(bh)
		}
		return nil
	case errSynced:
		// this is not actually error but flag that resync wasn't necessary
		w.is.FinishedSyncNoChange()
		return nil
	}

	w.metrics.IndexResyncErrors.With(common.Labels{"error": err.Error()}).Inc()

	return err
}

func (w *SyncWorker) resyncIndex(onNewBlock func(hash string)) error {
	remoteBestHash, err := w.chain.GetBestBlockHash()
	if err != nil {
		return err
	}
	localBestHeight, localBestHash, err := w.db.GetBestBlock()
	if err != nil {
		return err
	}
	// If the locally indexed block is the same as the best block on the network, we're done.
	if localBestHash == remoteBestHash {
		glog.Infof("resync: synced at %d %s", localBestHeight, localBestHash)
		return errSynced
	}
	if localBestHash != "" {
		remoteHash, err := w.chain.GetBlockHash(localBestHeight)
		// for some coins (eth) remote can be at lower best height after rollback
		if err != nil && err != bchain.ErrBlockNotFound {
			return err
		}
		if remoteHash != localBestHash {
			// forked - the remote hash differs from the local hash at the same height
			glog.Info("resync: local is forked at height ", localBestHeight, ", local hash ", localBestHash, ", remote hash", remoteHash)
			return w.handleFork(localBestHeight, localBestHash, onNewBlock)
		}
		glog.Info("resync: local at ", localBestHeight, " is behind")
		w.startHeight = localBestHeight + 1
	} else {
		// database is empty, start genesis
		glog.Info("resync: genesis from block ", w.startHeight)
	}
	w.startHash, err = w.chain.GetBlockHash(w.startHeight)
	if err != nil {
		return err
	}
	// if parallel operation is enabled and the number of blocks to be connected is large,
	// use parallel routine to load majority of blocks
	if w.syncWorkers > 1 {
		remoteBestHeight, err := w.chain.GetBestBlockHeight()
		if err != nil {
			return err
		}
		if remoteBestHeight < w.startHeight {
			glog.Error("resync: error - remote best height ", remoteBestHeight, " less than sync start height ", w.startHeight)
			return errors.New("resync: remote best height error")
		}
		if remoteBestHeight-w.startHeight > uint32(w.syncChunk) {
			glog.Infof("resync: parallel sync of blocks %d-%d, using %d workers", w.startHeight, remoteBestHeight, w.syncWorkers)
			err = w.ConnectBlocksParallel(w.startHeight, remoteBestHeight)
			if err != nil {
				return err
			}
			// after parallel load finish the sync using standard way,
			// new blocks may have been created in the meantime
			return w.resyncIndex(onNewBlock)
		}
	}
	return w.connectBlocks(onNewBlock)
}

func (w *SyncWorker) handleFork(localBestHeight uint32, localBestHash string, onNewBlock func(hash string)) error {
	// find forked blocks, disconnect them and then synchronize again
	var height uint32
	hashes := []string{localBestHash}
	for height = localBestHeight - 1; height >= 0; height-- {
		local, err := w.db.GetBlockHash(height)
		if err != nil {
			return err
		}
		if local == "" {
			break
		}
		remote, err := w.chain.GetBlockHash(height)
		// for some coins (eth) remote can be at lower best height after rollback
		if err != nil && err != bchain.ErrBlockNotFound {
			return err
		}
		if local == remote {
			break
		}
		hashes = append(hashes, local)
	}
	if err := w.DisconnectBlocks(height+1, localBestHeight, hashes); err != nil {
		return err
	}
	return w.resyncIndex(onNewBlock)
}

func (w *SyncWorker) connectBlocks(onNewBlock func(hash string)) error {
	bch := make(chan blockResult, 8)
	done := make(chan struct{})
	defer close(done)

	go w.getBlockChain(bch, done)

	var lastRes blockResult
	for res := range bch {
		lastRes = res
		if res.err != nil {
			return res.err
		}
		err := w.db.ConnectBlock(res.block)
		if err != nil {
			return err
		}
		if onNewBlock != nil {
			onNewBlock(res.block.Hash)
		}
		if res.block.Height > 0 && res.block.Height%1000 == 0 {
			glog.Info("connected block ", res.block.Height, " ", res.block.Hash)
		}
	}

	if lastRes.block != nil {
		glog.Infof("resync: synced at %d %s", lastRes.block.Height, lastRes.block.Hash)
	}

	return nil
}

// ConnectBlocksParallel uses parallel goroutines to get data from blockchain daemon
func (w *SyncWorker) ConnectBlocksParallel(lower, higher uint32) error {
	type hashHeight struct {
		hash   string
		height uint32
	}
	var err error
	var wg sync.WaitGroup
	bch := make(chan *bchain.Block, w.syncWorkers)
	hch := make(chan hashHeight, w.syncWorkers)
	hchClosed := atomic.Value{}
	hchClosed.Store(false)
	var getBlockMux sync.Mutex
	getBlockCond := sync.NewCond(&getBlockMux)
	lastConnectedBlock := lower - 1
	writeBlockDone := make(chan struct{})
	writeBlockWorker := func() {
		defer close(writeBlockDone)
		lastBlock := lower - 1
		for b := range bch {
			if lastBlock+1 != b.Height {
				glog.Error("writeBlockWorker skipped block, last connected block", lastBlock, ", new block ", b.Height)
			}
			err := w.db.ConnectBlock(b)
			if err != nil {
				glog.Error("writeBlockWorker ", b.Height, " ", b.Hash, " error ", err)
			}
			lastBlock = b.Height
		}
		glog.Info("WriteBlock exiting...")
	}
	getBlockWorker := func(i int) {
		defer wg.Done()
		var err error
		var block *bchain.Block
		for hh := range hch {
			for {
				block, err = w.chain.GetBlock(hh.hash, hh.height)
				if err != nil {
					// signal came while looping in the error loop
					if hchClosed.Load() == true {
						glog.Error("getBlockWorker ", i, " connect block error ", err, ". Exiting...")
						return
					}
					glog.Error("getBlockWorker ", i, " connect block error ", err, ". Retrying...")
					w.metrics.IndexResyncErrors.With(common.Labels{"error": err.Error()}).Inc()
					time.Sleep(time.Millisecond * 500)
				} else {
					break
				}
			}
			if w.dryRun {
				continue
			}
			getBlockMux.Lock()
			for {
				// we must make sure that the blocks are written to db in the correct order
				if lastConnectedBlock+1 == hh.height {
					// we have the right block, pass it to the writeBlockWorker
					lastConnectedBlock = hh.height
					bch <- block
					getBlockCond.Broadcast()
					break
				}
				// break the endless loop on OS signal
				if hchClosed.Load() == true {
					break
				}
				// wait for the time this block is top be passed to the writeBlockWorker
				getBlockCond.Wait()
			}
			getBlockMux.Unlock()
		}
		glog.Info("getBlockWorker ", i, " exiting...")
	}
	for i := 0; i < w.syncWorkers; i++ {
		wg.Add(1)
		go getBlockWorker(i)
	}
	go writeBlockWorker()
	var hash string
ConnectLoop:
	for h := lower; h <= higher; {
		select {
		case <-w.chanOsSignal:
			err = errors.Errorf("connectBlocksParallel interrupted at height %d", h)
			break ConnectLoop
		default:
			hash, err = w.chain.GetBlockHash(h)
			if err != nil {
				glog.Error("GetBlockHash error ", err)
				w.metrics.IndexResyncErrors.With(common.Labels{"error": err.Error()}).Inc()
				time.Sleep(time.Millisecond * 500)
				continue
			}
			hch <- hashHeight{hash, h}
			if h > 0 && h%1000 == 0 {
				glog.Info("connecting block ", h, " ", hash)
			}
			h++
		}
	}
	close(hch)
	// signal stop to workers that are in a loop
	hchClosed.Store(true)
	// broadcast syncWorkers times to unstuck all waiting getBlockWorkers
	for i := 0; i < w.syncWorkers; i++ {
		getBlockCond.Broadcast()
	}
	// first wait for the getBlockWorkers to finish and then close bch channel
	// so that the getBlockWorkers do not write to the closed channel
	wg.Wait()
	close(bch)
	<-writeBlockDone
	return err
}

type blockResult struct {
	block *bchain.Block
	err   error
}

func (w *SyncWorker) getBlockChain(out chan blockResult, done chan struct{}) {
	defer close(out)

	hash := w.startHash
	height := w.startHeight

	// some coins do not return Next hash
	// must loop until error
	for {
		select {
		case <-done:
			return
		default:
		}
		block, err := w.chain.GetBlock(hash, height)
		if err != nil {
			if err == bchain.ErrBlockNotFound {
				break
			}
			out <- blockResult{err: err}
			return
		}
		hash = block.Next
		height++
		out <- blockResult{block: block}
	}
}

// DisconnectBlocks removes all data belonging to blocks in range lower-higher,
// using block data from blockchain, if they are available,
// otherwise doing full scan
func (w *SyncWorker) DisconnectBlocks(lower uint32, higher uint32, hashes []string) error {
	glog.Infof("sync: disconnecting blocks %d-%d", lower, higher)
	// if the chain uses Block to Addresses mapping, always use DisconnectBlockRange
	if w.chain.GetChainParser().KeepBlockAddresses() > 0 {
		return w.db.DisconnectBlockRange(lower, higher)
	}
	blocks := make([]*bchain.Block, len(hashes))
	var err error
	// get all blocks first to see if we can avoid full scan
	for i, hash := range hashes {
		blocks[i], err = w.chain.GetBlock(hash, 0)
		if err != nil {
			// cannot get a block, we must do full range scan
			return w.db.DisconnectBlockRange(lower, higher)
		}
	}
	// then disconnect one after another
	for i, block := range blocks {
		glog.Info("Disconnecting block ", (int(higher) - i), " ", block.Hash)
		if err = w.db.DisconnectBlock(block); err != nil {
			return err
		}
	}
	return nil
}
