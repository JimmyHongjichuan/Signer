package main

import (
	"fmt"
	"github.com/JimmyHongjichuan/Signer/btc"
	"github.com/JimmyHongjichuan/btc_watcher/primitives"
)






func main() {
	address, err := btc.BitCoinHashToAddress("c19ba54b40598eab41f636b4b5c3fe6493dddd64", primitives.P2SH)
	//address, err := primitives.BitCoinHashToAddress("78c0c1bd8bf13af60f4b0371c2f5f9353de777c9", primitives.P2PKH)

	if err != nil {
		panic(err)
	}
	fmt.Println("BITCOIN ADDRESS: ", address)
	address_input := "03fcad68904be57b0a8133293ef4c3f49218df2da8886ed1257bc9f7e7d5a70e28"
	btc.GetInputAddress(address_input)
	script_input := "5221020fa7bed1b89df218a2ed2c94ebbf872a7bda0f48d231eb8cb6f16b87d9bb52112102d7e287092457f2bea226cd7537c5ee99af50cca923795a2ea65cf249f783c5d12102e8b48f3c0a7c452792fa96cdcf2fc6a23298f4d6512bd8aa9a25210b66a1d45053ae"
	btc.GetInputAddrP2SH(script_input)
}
