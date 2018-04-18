package main

import (
	"github.com/gin-gonic/gin"
	"flag"
	"fmt"
	"blockchain/miner"
)

var (
	flagCentralPort = flag.Int("port", 5001, "listening port")
	portCentralMin  = 5000
	portCentralMax  = 5010
)

func main() {
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	port := *flagCentralPort

	if port > portCentralMin && port < portCentralMax {
		r := gin.Default()
		miner.InitRouter(r)
		r.Run(fmt.Sprintf(":%d", port))
	} else {
		return
	}
}
