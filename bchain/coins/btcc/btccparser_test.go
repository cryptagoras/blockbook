// +build unittest

package btcc

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestBtccAddressEncodeAddress(t *testing.T) {
	addr1, err := newBtccAddress("13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji", Legacy)
	if err != nil {
		t.Errorf("newBtccAddress() error = %v", err)
		return
	}
	got1 := addr1.String()
	if got1 != "13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji" {
		t.Errorf("String() got1 = %v, want %v", got1, "13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji")
	}
}

func TestBtccAddressAreEqual(t *testing.T) {
	addr1, err := newBtccAddress("13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji", Legacy)
	if err != nil {
		t.Errorf("newBtccAddress() error = %v", err)
		return
	}
	got1 := addr1.AreEqual("13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji")
	if got1 != true {
		t.Errorf("AreEqual() got1 = %v, want %v", got1, true)
	}
	got3 := addr1.AreEqual("1HoKgKQh7ZNomWURmS9Tk3z8JM2MWm7S1w")
	if got3 != false {
		t.Errorf("AreEqual() got3 = %v, want %v", got3, false)
	}
}

func TestBtccAddressInSlice(t *testing.T) {
	addr1, err := newBtccAddress("13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji", Legacy)
	if err != nil {
		t.Errorf("newBtccAddress() error = %v", err)
		return
	}
	got1 := addr1.InSlice([]string{"13zMwGC5bxRn9ckJ1mgxf7UR8qbbNe2iji", "1HoKgKQh7ZNomWURmS9Tk3z8JM2MWm7S1w"})
	if got1 != true {
		t.Errorf("InSlice() got1 = %v, want %v", got1, true)
	}
	got3 := addr1.InSlice([]string{"1HoKgKQh7ZNomWURmS9Tk3z8JM2MWm7S1w", "1E6Np6dUPYpBSdLMLuwBF8sRQ3cngdaRRY"})
	if got3 != false {
		t.Errorf("InSlice() got3 = %v, want %v", got3, false)
	}
}

func TestAddressToOutputScript(t *testing.T) {
	parser, err := NewBtccParser(GetChainParams("test"), &btc.Configuration{AddressFormat: "legacy"})
	if err != nil {
		t.Errorf("NewBtccParser() error = %v", err)
		return
	}
	want, err := hex.DecodeString("76a9144fa927fd3bcf57d4e3c582c3d2eb2bd3df8df47c88ac")
	if err != nil {
		panic(err)
	}
	got1, err := parser.AddressToOutputScript("mnnAKPTSrWjgoi3uEYaQkHA1QEC5btFeBr")
	if err != nil {
		t.Errorf("AddressToOutputScript() error = %v", err)
		return
	}
	if !bytes.Equal(got1, want) {
		t.Errorf("AddressToOutputScript() got1 = %v, want %v", got1, want)
	}
}

var (
	testTx1, testTx2 bchain.Tx

	testTxPacked1 = "0001e2408ba8d7af5401000000017f9a22c9cbf54bd902400df746f138f37bcf5b4d93eb755820e974ba43ed5f42040000006a4730440220037f4ed5427cde81d55b9b6a2fd08c8a25090c2c2fff3a75c1a57625ca8a7118022076c702fe55969fa08137f71afd4851c48e31082dd3c40c919c92cdbc826758d30121029f6da5623c9f9b68a9baf9c1bc7511df88fa34c6c2f71f7c62f2f03ff48dca80feffffff019c9700000000000017a9146144d57c8aff48492c9dfb914e120b20bad72d6f8773d00700"
	testTxPacked2 = "0007c91a899ab7da6a010000000001019d64f0c72a0d206001decbffaa722eb1044534c74eee7a5df8318e42a4323ec10000000017160014550da1f5d25a9dae2eafd6902b4194c4c6500af6ffffffff02809698000000000017a914cd668d781ece600efa4b2404dc91fd26b8b8aed8870553d7360000000017a914246655bdbd54c7e477d0ea2375e86e0db2b8f80a8702473044022076aba4ad559616905fa51d4ddd357fc1fdb428d40cb388e042cdd1da4a1b7357022011916f90c712ead9a66d5f058252efd280439ad8956a967e95d437d246710bc9012102a80a5964c5612bb769ef73147b2cf3c149bc0fd4ecb02f8097629c94ab013ffd00000000"
)

func init() {
	var (
		addr1, addr2, addr3 bchain.Address
		err                 error
	)
	addr1, err = newBtccAddress("3AZKvpKhSh1o8t1QrX3UeXG9d2BhCRnbcK", Legacy)
	if err == nil {
		addr2, err = newBtccAddress("2NByHN6A8QYkBATzxf4pRGbCSHD5CEN2TRu", Legacy)
	}
	if err == nil {
		addr3, err = newBtccAddress("2MvZguYaGjM7JihBgNqgLF2Ca2Enb76Hj9D", Legacy)
	}
	if err != nil {
		panic(err)
	}

	testTx1 = bchain.Tx{
		Hex:       "01000000017f9a22c9cbf54bd902400df746f138f37bcf5b4d93eb755820e974ba43ed5f42040000006a4730440220037f4ed5427cde81d55b9b6a2fd08c8a25090c2c2fff3a75c1a57625ca8a7118022076c702fe55969fa08137f71afd4851c48e31082dd3c40c919c92cdbc826758d30121029f6da5623c9f9b68a9baf9c1bc7511df88fa34c6c2f71f7c62f2f03ff48dca80feffffff019c9700000000000017a9146144d57c8aff48492c9dfb914e120b20bad72d6f8773d00700",
		Blocktime: 1519053802,
		Txid:      "056e3d82e5ffd0e915fb9b62797d76263508c34fe3e5dbed30dd3e943930f204",
		LockTime:  512115,
		Version:   1,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "4730440220037f4ed5427cde81d55b9b6a2fd08c8a25090c2c2fff3a75c1a57625ca8a7118022076c702fe55969fa08137f71afd4851c48e31082dd3c40c919c92cdbc826758d30121029f6da5623c9f9b68a9baf9c1bc7511df88fa34c6c2f71f7c62f2f03ff48dca80",
				},
				Txid:     "425fed43ba74e9205875eb934d5bcf7bf338f146f70d4002d94bf5cbc9229a7f",
				Vout:     4,
				Sequence: 4294967294,
			},
		},
		Vout: []bchain.Vout{
			{
				Value: 0.00038812,
				N:     0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a9146144d57c8aff48492c9dfb914e120b20bad72d6f87",
					Addresses: []string{
						"3AZKvpKhSh1o8t1QrX3UeXG9d2BhCRnbcK",
					},
				},
				Address: addr1,
			},
		},
	}

	testTx2 = bchain.Tx{
		Hex:       "010000000001019d64f0c72a0d206001decbffaa722eb1044534c74eee7a5df8318e42a4323ec10000000017160014550da1f5d25a9dae2eafd6902b4194c4c6500af6ffffffff02809698000000000017a914cd668d781ece600efa4b2404dc91fd26b8b8aed8870553d7360000000017a914246655bdbd54c7e477d0ea2375e86e0db2b8f80a8702473044022076aba4ad559616905fa51d4ddd357fc1fdb428d40cb388e042cdd1da4a1b7357022011916f90c712ead9a66d5f058252efd280439ad8956a967e95d437d246710bc9012102a80a5964c5612bb769ef73147b2cf3c149bc0fd4ecb02f8097629c94ab013ffd00000000",
		Blocktime: 1235678901,
		Txid:      "474e6795760ebe81cb4023dc227e5a0efe340e1771c89a0035276361ed733de7",
		LockTime:  0,
		Version:   1,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "160014550da1f5d25a9dae2eafd6902b4194c4c6500af6",
				},
				Txid:     "c13e32a4428e31f85d7aee4ec7344504b12e72aaffcbde0160200d2ac7f0649d",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				Value: .1,
				N:     0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a914cd668d781ece600efa4b2404dc91fd26b8b8aed887",
					Addresses: []string{
						"3LR5JME6o6EpxgNQywCYeeDB4rs2M5pqaW",
					},
				},
				Address: addr2,
			},
			{
				Value: 9.20081157,
				N:     1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a914246655bdbd54c7e477d0ea2375e86e0db2b8f80a87",
					Addresses: []string{
						"351UqoeF7tbxWuZ8hi4Td5DJotaRPAGRQ9",
					},
				},
				Address: addr3,
			},
		},
	}
}

func Test_UnpackTx(t *testing.T) {
	parser1, err := NewBtccParser(GetChainParams("main"), &btc.Configuration{AddressFormat: "legacy"})
	if err != nil {
		t.Errorf("NewBtccParser() error = %v", err)
		return
	}
	parser2, err := NewBtccParser(GetChainParams("test"), &btc.Configuration{AddressFormat: "legacy"})
	if err != nil {
		t.Errorf("NewBtccParser() error = %v", err)
		return
	}

	type args struct {
		packedTx string
		parser   *BtccParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "btc-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   parser1,
			},
			want:    &testTx1,
			want1:   123456,
			wantErr: false,
		},
		{
			name: "testnet-1",
			args: args{
				packedTx: testTxPacked2,
				parser:   parser2,
			},
			want:    &testTx2,
			want1:   510234,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
