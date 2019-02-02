package main

import (
	"fmt"
	"github.com/JimmyHongjichuan/Signer/btc"
	"github.com/btcsuite/btcd/chaincfg"
)






func main() {
	address, err := btc.BitCoinHashToAddress("be55f02c63d6cb855c953e44cb86f2e57d01e3e5", btc.P2SH_PREFIX)
	//address, err := btc.BitCoinHashToAddress("d626f681ab699644fc0120642fb71f590a4b86fd", btc.REGTEST_PREFIX)

	if err != nil {
		panic(err)
	}
	fmt.Println("BITCOIN ADDRESS: ", address)
	address_input := "0020fa28dc1e5eb222055e90f8cade9bcd13ca9ddab7a5ed029e27d41a736f7455ce" //31w3iWUN5EMJMW2YRCc5m4RFqm3zN61xK2
	btc.GetInputAddress(address_input, &chaincfg.MainNetParams, btc.ScriptHash)
	address_input = "0020fa28dc1e5eb222055e90f8cade9bcd13ca9ddab7a5ed029e27d41a736f7455ce"//
	param := chaincfg.RegressionNetParams;
	param.PubKeyHashAddrID = 0xc4
	btc.GetInputAddress(address_input, &param, btc.ScriptHash)
	address_input = "fa28dc1e5eb222055e90f8cade9bcd13ca9ddab7a5ed029e27d41a736f7455ce"
	btc.GetInputAddress(address_input, &chaincfg.MainNetParams, btc.WitnessScriptHash)
	address_input = "2c86f6f95f2fbcf5d7fe9e2e87f9860af9041e5f"
	btc.GetInputAddress(address_input, &param, btc.WitnessPubKeyHash)
	btc.GenTx(&chaincfg.MainNetParams)

}
