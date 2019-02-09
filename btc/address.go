package btc

import (
	"crypto/sha256"
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

const P2PKH_PREFIX = byte(0x00)
const P2SH_PREFIX = byte(0x05)
const TESTNET_SCRIPT_PREFIX = byte(0xc4)
const REGTEST_PREFIX = byte(0xc4)

const (
	PubKeyHash int = iota
	ScriptHash
	WitnessPubKeyHash
	WitnessScriptHash
)

type Network struct {
	name        string
	symbol      string
	xpubkey     byte
	xprivatekey byte
}

var Coinnetwork = map[string]Network{
	"rdd": {name: "reddcoin", symbol: "rdd", xpubkey: 0x3d, xprivatekey: 0xbd},
	"dgb": {name: "digibyte", symbol: "dgb", xpubkey: 0x1e, xprivatekey: 0x80},
	"btc": {name: "bitcoin", symbol: "btc", xpubkey: 0x00, xprivatekey: 0x80},
	"ltc": {name: "litecoin", symbol: "ltc", xpubkey: 0x30, xprivatekey: 0xb0},
}


func (network Network) GetNetworkParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.PubKeyHashAddrID = network.xpubkey
	networkParams.PrivateKeyID = network.xprivatekey
	return networkParams
}

func (network Network) CreatePrivateKey() (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, network.GetNetworkParams(), true)
}


func (network Network) GetAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), network.GetNetworkParams())
}

func (network Network) ImportWIF(wifStr string) (*btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(network.GetNetworkParams()) {
		return nil, errors.New("The WIF string is not valid for the `" + network.name + "` network")
	}
	return wif, nil
}


func BitCoinHashToAddress(pkscript string, script_type byte) (string, error) {
	hash, err := hex.DecodeString(pkscript)
	if err != nil {
		log.Fatal(err)
	}

	pf := append([]byte{script_type}, hash...)
	b := append(pf, checkSum(pf)...)

	address := base58.Encode(b)
	return address, nil
}

func checkSum(publicKeyHash []byte) []byte {
	first := sha256.Sum256(publicKeyHash)
	second := sha256.Sum256(first[:])
	return second[0:4]
}


func GetInputAddress(script string, net *chaincfg.Params, hash_type int) {
	chainParams := net;
	PubKey, err := hex.DecodeString(script)
	if err != nil {
		panic(err)
	}

	var addressPubKey *btcutil.AddressPubKey
	var addressScript *btcutil.AddressScriptHash
	var addressWitnessPubKey *btcutil.AddressWitnessPubKeyHash
	var addressWitnessScriptKey *btcutil.AddressWitnessScriptHash
	if (hash_type == PubKeyHash) {

		addressPubKey, err = btcutil.NewAddressPubKey(PubKey, chainParams)
		if err != nil {
			log.Fatal(err)
		}
		address := addressPubKey.EncodeAddress()
		fmt.Println(address)
	} else if (hash_type == ScriptHash) {
		addressScript, err = btcutil.NewAddressScriptHash(PubKey, chainParams)
		if err != nil {
			log.Fatal(err)
		}
		address := addressScript.EncodeAddress()
		fmt.Println(address)
	} else if (hash_type == WitnessPubKeyHash){
		addressWitnessPubKey, err = btcutil.NewAddressWitnessPubKeyHash(PubKey, chainParams)
		if err != nil {
			log.Fatal(err)
		}
		address := 	addressWitnessPubKey.EncodeAddress()
		fmt.Println(address)
	} else if(hash_type == WitnessScriptHash){
		addressWitnessScriptKey, err = btcutil.NewAddressWitnessScriptHash(PubKey, chainParams)
		if err != nil {
			log.Fatal(err)
		}
		address := 	addressWitnessScriptKey.EncodeAddress()
		fmt.Println(address)
	}

}





