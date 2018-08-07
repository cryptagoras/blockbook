package btcc

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/bchutil"
	"github.com/schancel/cashaddr-converter/address"
)

type AddressFormat = uint8

const (
	Legacy AddressFormat = iota
)

const (
	MainNetPrefix = "bitcoincore:"
	TestNetPrefix = "btcctest:"
	RegTestPrefix = "btccreg:"
)

// BtccParser handle
type BtccParser struct {
	*btc.BitcoinParser
	AddressFormat AddressFormat
}

// NewBtccParser returns new BtccParser instance
func NewBtccParser(params *chaincfg.Params, c *btc.Configuration) (*BtccParser, error) {
	var format AddressFormat
	switch c.AddressFormat {
	case "":
		fallthrough
	case "legacy":
		format = Legacy
	default:
		return nil, fmt.Errorf("Unknown address format: %s", c.AddressFormat)
	}
	p := &BtccParser{
		BitcoinParser: &btc.BitcoinParser{
			BaseParser: &bchain.BaseParser{
				AddressFactory:       func(addr string) (bchain.Address, error) { return newBtccAddress(addr, format) },
				BlockAddressesToKeep: c.BlockAddressesToKeep,
			},
			Params: params,
			OutputScriptToAddressesFunc: outputScriptToAddresses,
		},
		AddressFormat: format,
	}
	return p, nil
}

// GetChainParams contains network parameters for the main Bitcoin Cash network,
// the regression test Bitcoin Cash network, the test Bitcoin Cash network and
// the simulation test Bitcoin Cash network, in this order
func GetChainParams(chain string) *chaincfg.Params {
	var params *chaincfg.Params
	switch chain {
	case "test":
		params = &chaincfg.TestNet3Params
		params.Net = bchutil.TestnetMagic
	case "regtest":
		params = &chaincfg.RegressionNetParams
		params.Net = bchutil.Regtestmagic
	default:
		params = &chaincfg.MainNetParams
		params.Net = bchutil.MainnetMagic
	}

	return params
}

// GetAddrIDFromAddress returns internal address representation of given address
func (p *BtccParser) GetAddrIDFromAddress(address string) ([]byte, error) {
	return p.AddressToOutputScript(address)
}

// AddressToOutputScript converts bitcoin address to ScriptPubKey
func (p *BtccParser) AddressToOutputScript(address string) ([]byte, error) {
	da, err := btcutil.DecodeAddress(address, p.Params)
	if err != nil {
		return nil, err
	}
	script, err := txscript.PayToAddrScript(da)
	if err != nil {
		return nil, err
	}
	return script, nil
}

// outputScriptToAddresses converts ScriptPubKey to bitcoin addresses
func outputScriptToAddresses(script []byte, params *chaincfg.Params) ([]string, error) {
	a, err := bchutil.ExtractPkScriptAddrs(script, params)
	if err != nil {
		return nil, err
	}
	return []string{a.EncodeAddress()}, nil
}

type btccAddress struct {
	addr string
}

func newBtccAddress(addr string, format AddressFormat) (*btccAddress, error) {
	da, err := address.NewFromString(addr)
	if err != nil {
		return nil, err
	}
	var ea string
	switch format {
	case Legacy:
		if a, err := da.Legacy(); err != nil {
			return nil, err
		} else {
			ea, err = a.Encode()
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("Unknown address format: %d", format)
	}
	return &btccAddress{addr: ea}, nil
}

func (a *btccAddress) String() string {
	return a.addr
}

func (a *btccAddress) AreEqual(addr string) bool {
	return a.String() == addr
}

func (a *btccAddress) InSlice(addrs []string) bool {
	ea := a.String()
	for _, addr := range addrs {
		if ea == addr {
			return true
		}
	}
	return false
}
