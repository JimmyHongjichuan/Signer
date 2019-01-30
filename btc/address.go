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
const (
	P2PKH int = iota
	P2SH
)

func BitCoinHashToAddress(pkscript string, script_type int) (string, error) {
	hash, err := hex.DecodeString(pkscript)
	if err != nil {
		log.Fatal(err)
	}
	var prefix byte
	if script_type == P2PKH {
		prefix = P2PKH_PREFIX
		fmt.Println("P2PKH")
	} else {
		prefix = P2SH_PREFIX
		fmt.Println("P2SH")
	}
	pf := append([]byte{prefix}, hash...)
	b := append(pf, checkSum(pf)...)

	address := base58.Encode(b)
	return address, nil
}

func checkSum(publicKeyHash []byte) []byte {
	first := sha256.Sum256(publicKeyHash)
	second := sha256.Sum256(first[:])
	return second[0:4]
}


func GetInputAddress(addres string){
	chainParams :=  &chaincfg.MainNetParams;
	PubKey , err := hex.DecodeString(addres)
	if err != nil {
		panic(err)
	}
	addressPubKey, err := btcutil.NewAddressPubKey(PubKey, chainParams)
	if err != nil {
		log.Fatal(err)
	}
	address := addressPubKey.EncodeAddress()
	fmt.Println (address)
}

func GetInputAddrP2SH(addres string){
	chainParams :=  &chaincfg.MainNetParams;
	script , err := hex.DecodeString(addres)
	if err != nil {
		log.Fatal(err)
	}
	addressScript, err := btcutil.NewAddressScriptHash(script, chainParams)
	if err != nil {
		log.Fatal(err)
	}
	address := addressScript.EncodeAddress()
	fmt.Println (address)
}
