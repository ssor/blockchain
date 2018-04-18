package main

import (
	"github.com/gin-gonic/gin"
	"flag"
	"fmt"
	"blockchain/miner"
)

var (
	flagWalletPort = flag.Int("port", 5001, "listening port")
	portWalletMin  = 5000
	portWalletMax  = 5010
)

func main() {
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	port := *flagWalletPort

	if port > portWalletMin && port < portWalletMax {
		r := gin.Default()
		miner.InitRouter(r)
		r.Run(fmt.Sprintf(":%d", port))
	} else {
		return
	}
}
