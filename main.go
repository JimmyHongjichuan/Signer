package main

import (
	"fmt"
	"github.com/JimmyHongjichuan/Signer/btc"
	"github.com/btcsuite/btcd/chaincfg"
)




func main() {
	fmt.Println("Starting the application...")
	wif, _ := btc.Coinnetwork["btc"].CreatePrivateKey()
	Genaddress, _ := btc.Coinnetwork["btc"].GetAddress(wif)
	fmt.Printf("%s - %s\n", wif.String(), Genaddress.EncodeAddress())
	//transaction, err := btc.CreateTransaction("cSkELxYraVBYBeU1QvoasNYzdWJkXoS5x1LK7PMLE1q74TZTYMZG", "n1yJ5g9k5zSdU9iLGyjLhuF8RYvmVp5TR3", 499900000, "42a8cc0c246783d1d0c4d382938e6f47667bd7d108ab9bcb804710075399f827",&chaincfg.RegressionNetParams, true)
	////transaction, err := btc.CreateTransaction("cS5LWK2aUKgP9LmvViG3m9HkfwjaEJpGVbrFHuGZKvW2ae3W9aUe", "mrdKfqWEkwferzEQus5NpgK2Dtpq7Qcgif", 499900000, "12e0d25258ec29fadf75a3f569fccaeeb8ca4af5d2d34e9a48ab5a6fdc0efc1e",&chaincfg.TestNet3Params)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//data, _ := json.Marshal(transaction)
	//fmt.Println(string(data))

	address, err := btc.BitCoinHashToAddress("160bcebc8f48a635720638ecb8e6a11e8079b25a", btc.P2SH_PREFIX)
	//address, err := btc.BitCoinHashToAddress("78ff60a652028c5a898aeb32ba8cc7c57e039554", btc.REGTEST_PREFIX)
	//address, err := btc.BitCoinHashToAddress("160bcebc8f48a635720638ecb8e6a11e8079b25a", btc.P2PKH_PREFIX)
	if err != nil {
		panic(err)
	}
	fmt.Println("BITCOIN ADDRESS: ", address)
	address_input := "0020fa28dc1e5eb222055e90f8cade9bcd13ca9ddab7a5ed029e27d41a736f7455ce" //31w3iWUN5EMJMW2YRCc5m4RFqm3zN61xK2
	btc.GetInputAddress(address_input, &chaincfg.MainNetParams, btc.ScriptHash)

	address_input = "160bcebc8f48a635720638ecb8e6a11e8079b25a"//
	btc.GetInputAddress(address_input, &chaincfg.MainNetParams, btc.ScriptHash)

	address_input = "fa28dc1e5eb222055e90f8cade9bcd13ca9ddab7a5ed029e27d41a736f7455ce"
	btc.GetInputAddress(address_input, &chaincfg.MainNetParams, btc.WitnessScriptHash)

	address_input = "2c86f6f95f2fbcf5d7fe9e2e87f9860af9041e5f"
	btc.GetInputAddress(address_input, &chaincfg.RegressionNetParams, btc.WitnessPubKeyHash)

	txoutInfos := make([]btc.TxOutAddressInfo, 0)
	txoutInfos = append(txoutInfos, btc.TxOutAddressInfo{"n1yJ5g9k5zSdU9iLGyjLhuF8RYvmVp5TR3", 800000000})
	txoutInfos = append(txoutInfos, btc.TxOutAddressInfo{"mgbMRCieM4b2owd8K6bDd7fjhzNPMirqb7", 998500000})
	//btc.GenTx(&chaincfg.MainNetParams)
	btc.GenTx(&chaincfg.RegressionNetParams, "cRFq9JJGVqjGe9RDpT49BxCPbGJoy1N8jGk4Muit8iP1JtpJE4XA",
		1799000000,txoutInfos,
		"d3bb20262b32baef98ab3facbdb89323f07b703d0814cc13f76c2ef39e88cc60", 0,true)
}
