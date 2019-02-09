package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)


func GenTx(net *chaincfg.Params) {
	privWif := "L5MFEsN1sXwtDsZvTRoZmzEK73GYJibkPvvtYLr3Eju3L4s92Xuv"//"KwLQYGA7nsvQqZeJP38qPQHgjcPZvE79jRQgPgu5tAHvFF8gWm3n"
	txHash := "a41e8a514ba537ff618528add6910e315193094df0ff2345f540f4206079c005"
	destination := "2NDF4ygHLwp6JRPGrmv4j6SMmYyq3QVWWZV"
	amount := int64(11650795)
	txFee := int64(500000)
	sourceUTXOIndex := uint32(1)
	chainParams := net

	decodedWif, err := btcutil.DecodeWIF(privWif)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded WIF: %v\n", decodedWif) // Decoded WIF: cS5LWK2aUKgP9LmvViG3m9HkfwjaEJpGVbrFHuGZKvW2ae3W9aUe
	//PubKey, err := hex.DecodeString("030b4bbfeca237a4bab81a3adeef76cc1cbcfa5e7cac5c22754e47ba42e1fe9579")
	//if err != nil {
	//	panic(err)
	//}

	addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeUncompressed(), chainParams)//decodedWif.PrivKey.PubKey().SerializeUncompressed()
	if err != nil {
		log.Fatal(err)
	}

	sourceUTXOHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UTXO hash: %s\n", sourceUTXOHash) // utxo hash: 12e0d25258ec29fadf75a3f569fccaeeb8ca4af5d2d34e9a48ab5a6fdc0efc1e

	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(destination, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source Address: %s\n", sourceAddress) // Source Address: mgjHgKi1g6qLFBM1gQwuMjjVBGMJdrs9pP

	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		log.Fatal(err)
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		log.Fatal(err)
	}

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	redeemTx.AddTxIn(sourceTxIn)
	redeemTxOut := wire.NewTxOut((amount - txFee), destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)

	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTxOut.PkScript, txscript.SigHashAll, decodedWif.PrivKey, false)
	if err != nil {
		log.Fatal(err)
	}

	redeemTx.TxIn[0].SignatureScript = sigScript
	fmt.Printf("Signature Script: %v\n", hex.EncodeToString(sigScript)) // Signature Script: 473...b67

	// validate signature
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTxOut.PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		log.Fatal(err)
	}

	if err := vm.Execute(); err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, redeemTx.SerializeSize()))
	redeemTx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes())) // redeem Tx: 01000000011efc...5bb88ac00000000
}
