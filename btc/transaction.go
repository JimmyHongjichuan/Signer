package btc

import (
	"errors"
	"github.com/btcsuite/btcd/btcec"
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


type Transaction struct {
	TxId               string `json:"txid"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Amount             int64  `json:"amount"`
	UnsignedTx         string `json:"unsignedtx"`
	SignedTx           string `json:"signedtx"`
}

type TxOutAddressInfo struct {
	DestinationAddress string `json:"destination_address"`
	PayAmount             int64  `json:"pay_amount"`
}

func CreateTransaction(secret string, destination string, amount int64, txHash string, net *chaincfg.Params,WIF_compress bool) (Transaction, error) {
	var transaction Transaction
	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return Transaction{}, err
	}
	var addresspubkey *btcutil.AddressPubKey;
	if (WIF_compress == true) {
		addresspubkey, _ = btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), net)
	} else {
		addresspubkey, _ = btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), net)
	}
	sourceTx := wire.NewMsgTx(wire.TxVersion)
	sourceUtxoHash, _ := chainhash.NewHashFromStr(txHash)
	sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 1)
	sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
	destinationAddress, err := btcutil.DecodeAddress(destination, net)
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), net)
	if err != nil {
		return Transaction{}, err
	}
	fmt.Printf("Source Address: %s\n", sourceAddress)
	destinationPkScript, _ := txscript.PayToAddrScript(destinationAddress)
	sourcePkScript, _ := txscript.PayToAddrScript(sourceAddress)
	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)
	sourceTxHash := sourceTx.TxHash()
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&sourceTxHash, 1)
	redeemTxIn := wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, WIF_compress)
	if err != nil {
		return Transaction{}, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return Transaction{}, err
	}
	if err := vm.Execute(); err != nil {
		return Transaction{}, err
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	return transaction, nil
}
type addressToKey struct {
	key        *btcec.PrivateKey
	compressed bool
}

func mkGetKey(keys map[string]addressToKey) txscript.KeyDB {
	if keys == nil {
		return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey,
			bool, error) {
			return nil, false, errors.New("nope")
		})
	}
	return txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey,
		bool, error) {
		a2k, ok := keys[addr.EncodeAddress()]
		if !ok {
			return nil, false, errors.New("nope")
		}
		return a2k.key, a2k.compressed, nil
	})
}

func checkScripts(msg string, tx *wire.MsgTx, idx int, inputAmt int64, sigScript, pkScript []byte) error {
	tx.TxIn[idx].SignatureScript = sigScript
	vm, err := txscript.NewEngine(pkScript, tx, idx,
		txscript.ScriptBip16|txscript.ScriptVerifyDERSignatures, nil, nil, inputAmt)
	if err != nil {
		return fmt.Errorf("failed to make script engine for %s: %v",
			msg, err)
	}

	err = vm.Execute()
	if err != nil {
		return fmt.Errorf("invalid script signature for %s: %v", msg,
			err)
	}

	return nil
}

func signAndCheck(net *chaincfg.Params, msg string, tx *wire.MsgTx, idx int, inputAmt int64, pkScript []byte,
	hashType txscript.SigHashType, kdb txscript.KeyDB, sdb txscript.ScriptDB,
	previousScript []byte) error {

	sigScript, err := txscript.SignTxOutput(net, tx, idx,
		pkScript, hashType, kdb, sdb, nil)
	if err != nil {
		return fmt.Errorf("failed to sign output %s: %v", msg, err)
	}

	return checkScripts(msg, tx, idx, inputAmt, sigScript, pkScript)
}

func mkGetScript(scripts map[string][]byte) txscript.ScriptDB {
	if scripts == nil {
		return txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
			return nil, errors.New("nope")
		})
	}
	return txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
		script, ok := scripts[addr.EncodeAddress()]
		if !ok {
			return nil, errors.New("nope")
		}
		return script, nil
	})
}


func GenMultiSigTx(net *chaincfg.Params,privWif []string, txHash string, amount int64,  txoutInfos []TxOutAddressInfo, txinIndex int32,  WIF_compress bool){
	var addressPubKeys []*btcutil.AddressPubKey = make([]*btcutil.AddressPubKey ,0)

	var multisiginfos =  make(map[string]addressToKey)
	for _, wif := range privWif {
		decodedWif, err := btcutil.DecodeWIF(wif)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Decoded WIF: %v\n", decodedWif)

		var addressPubKey *btcutil.AddressPubKey
		if (WIF_compress == true) {
			addressPubKey, err = btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), net)

			if err != nil {
				log.Fatal(err)
			}
			addressPubKeys = append(addressPubKeys,addressPubKey)
		} else {
			addressPubKey, err = btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeUncompressed(), net)
			if err != nil {
				log.Fatal(err)
			}
			addressPubKeys = append(addressPubKeys,addressPubKey)
		}
		multisiginfos[addressPubKey.EncodeAddress()] = addressToKey{decodedWif.PrivKey, true}
	}
	pkScript, err := txscript.MultiSigScript(
		addressPubKeys,
		2)
	if err != nil {
		log.Fatal(err)
	}

	scriptAddr, err := btcutil.NewAddressScriptHash(
		pkScript, net)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Source Address: %s\n", scriptAddr)
	scriptPkScript, err := txscript.PayToAddrScript(scriptAddr)
	if err != nil {
		log.Fatal(err)
	}

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	sourceUTXOHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UTXO hash: %s\n", sourceUTXOHash) // utxo hash: 12e0d25258ec29fadf75a3f569fccaeeb8ca4af5d2d34e9a48ab5a6fdc0efc1e
	sourceUTXOIndex := uint32(txinIndex)
	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)
	redeemTx.AddTxIn(sourceTxIn)

	for _, txoutInfo := range txoutInfos {
		destinationAddress, err := btcutil.DecodeAddress(txoutInfo.DestinationAddress, net)
		if err != nil {
			log.Fatal(err)
		}
		destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
		if err != nil {
			log.Fatal(err)
		}
		redeemTxOut := wire.NewTxOut(txoutInfo.PayAmount, destinationPkScript)
		redeemTx.AddTxOut(redeemTxOut)
	}

	msg := fmt.Sprintf("%d:%d", txscript.SigHashAll, txinIndex)
	if err := signAndCheck(net, msg, redeemTx, 0, amount,
		scriptPkScript, txscript.SigHashAll,
		mkGetKey(multisiginfos), mkGetScript(map[string][]byte{
			scriptAddr.EncodeAddress(): pkScript,
		}),
		nil); err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, redeemTx.SerializeSize()))
	redeemTx.Serialize(buf)

	fmt.Printf("Redeem multisig Tx: %v\n", hex.EncodeToString(buf.Bytes()))
}

func GenTx(net *chaincfg.Params,privWif string,amount int64, txoutInfos []TxOutAddressInfo, txHash string , txinIndex int32,WIF_compress bool) {
	//privWif := "cSkELxYraVBYBeU1QvoasNYzdWJkXoS5x1LK7PMLE1q74TZTYMZG"
	//txHash := "42a8cc0c246783d1d0c4d382938e6f47667bd7d108ab9bcb804710075399f827"
	//destination := "n1yJ5g9k5zSdU9iLGyjLhuF8RYvmVp5TR3"
	//amount := int64(1800000000)
	//txFee := int64(500000)
	sourceUTXOIndex := uint32(txinIndex)
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
	var addressPubKey *btcutil.AddressPubKey;
	if (WIF_compress == true) {
		addressPubKey, err = btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), net)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		addressPubKey, err = btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeUncompressed(), net)
		if err != nil {
			log.Fatal(err)
		}
	}
	//addressPubKey, err := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeCompressed(), chainParams)//decodedWif.PrivKey.PubKey().SerializeUncompressed()


	sourceUTXOHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UTXO hash: %s\n", sourceUTXOHash) // utxo hash: 12e0d25258ec29fadf75a3f569fccaeeb8ca4af5d2d34e9a48ab5a6fdc0efc1e

	sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
	sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)

	sourceAddress, err := btcutil.DecodeAddress(addressPubKey.EncodeAddress(), chainParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Source Address: %s\n", sourceAddress) // Source Address: mgjHgKi1g6qLFBM1gQwuMjjVBGMJdrs9pP



	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		log.Fatal(err)
	}

	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	redeemTx.AddTxIn(sourceTxIn)

	for _, txoutInfo := range txoutInfos {
		destinationAddress, err := btcutil.DecodeAddress(txoutInfo.DestinationAddress, chainParams)
		if err != nil {
			log.Fatal(err)
		}
		destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
		if err != nil {
			log.Fatal(err)
		}
		redeemTxOut := wire.NewTxOut(txoutInfo.PayAmount, destinationPkScript)
		redeemTx.AddTxOut(redeemTxOut)
	}
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTxOut.PkScript, txscript.SigHashAll, decodedWif.PrivKey, WIF_compress)
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
