package btc

import (
	"crypto/sha256"
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





