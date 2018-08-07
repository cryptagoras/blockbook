// +build integration

package bch

import (
	"blockbook/bchain"
	"blockbook/bchain/tests/rpc"
	"encoding/json"
	"flag"
	"os"
	"testing"
)

func getRPCClient(chain string) func(json.RawMessage) (bchain.BlockChain, error) {
	return func(cfg json.RawMessage) (bchain.BlockChain, error) {
		c, err := NewBtccRPC(cfg, nil)
		if err != nil {
			return nil, err
		}
		cli := c.(*BtccRPC)
		cli.Parser, err = NewBtccParser(GetChainParams(chain), cli.ChainConfig)
		if err != nil {
			return nil, err
		}
		cli.Mempool = bchain.NewUTXOMempool(cli, cli.ChainConfig.MempoolWorkers, cli.ChainConfig.MempoolSubWorkers)
		return cli, nil
	}
}

var tests struct {
	mainnet *rpc.Test
	testnet *rpc.Test
}

func TestMain(m *testing.M) {
	flag.Parse()

	t, err := rpc.NewTest("Bitcoin Core", getRPCClient("main"))
	if err != nil {
		panic(err)
	}

	tests.mainnet = t

	t, err = rpc.NewTest("Bitcoin Core Testnet", getRPCClient("test"))
	if err != nil {
		panic(err)
	}

	tests.testnet = t

	os.Exit(m.Run())
}

func TestBtccRPC_GetBlockHash(t *testing.T) {
	tests.mainnet.TestGetBlockHash(t)
}

func TestBtccRPC_GetBlock(t *testing.T) {
	tests.mainnet.TestGetBlock(t)
}

func TestBtccRPC_GetTransaction(t *testing.T) {
	tests.mainnet.TestGetTransaction(t)
}

func TestBtccRPC_GetTransactionForMempool(t *testing.T) {
	tests.mainnet.TestGetTransactionForMempool(t)
}

func TestBtccRPC_MempoolSync(t *testing.T) {
	tests.mainnet.TestMempoolSync(t)
}

func TestBtccRPC_GetMempoolEntry(t *testing.T) {
	tests.mainnet.TestGetMempoolEntry(t)
}

func TestBtccRPC_EstimateSmartFee(t *testing.T) {
	tests.mainnet.TestEstimateSmartFee(t)
}

func TestBtccRPC_EstimateFee(t *testing.T) {
	tests.mainnet.TestEstimateFee(t)
}

func TestBtccRPC_GetBestBlockHash(t *testing.T) {
	tests.mainnet.TestGetBestBlockHash(t)
}

func TestBtccRPC_GetBestBlockHeight(t *testing.T) {
	tests.mainnet.TestGetBestBlockHeight(t)
}

func TestBtccRPC_GetBlockHeader(t *testing.T) {
	tests.mainnet.TestGetBlockHeader(t)
}

func TestBtccTestnetRPC_GetBlockHash(t *testing.T) {
	tests.testnet.TestGetBlockHash(t)
}

func TestBtccTestnetRPC_GetBlock(t *testing.T) {
	tests.testnet.TestGetBlock(t)
}

func TestBtccTestnetRPC_GetTransaction(t *testing.T) {
	tests.testnet.TestGetTransaction(t)
}

func TestBtccTestnetRPC_GetTransactionForMempool(t *testing.T) {
	tests.testnet.TestGetTransactionForMempool(t)
}

func TestBtccTestnetRPC_MempoolSync(t *testing.T) {
	tests.testnet.TestMempoolSync(t)
}

func TestBtccTestnetRPC_GetMempoolEntry(t *testing.T) {
	tests.testnet.TestGetMempoolEntry(t)
}

func TestBtccTestnetRPC_EstimateSmartFee(t *testing.T) {
	tests.testnet.TestEstimateSmartFee(t)
}

func TestBtccTestnetRPC_EstimateFee(t *testing.T) {
	tests.testnet.TestEstimateFee(t)
}

func TestBtccTestnetRPC_GetBestBlockHash(t *testing.T) {
	tests.testnet.TestGetBestBlockHash(t)
}

func TestBtccTestnetRPC_GetBestBlockHeight(t *testing.T) {
	tests.testnet.TestGetBestBlockHeight(t)
}

func TestBtccTestnetRPC_GetBlockHeader(t *testing.T) {
	tests.testnet.TestGetBlockHeader(t)
}
