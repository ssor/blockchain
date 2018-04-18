package main

import (
	"github.com/gin-gonic/gin"
	"flag"
	"fmt"
	"blockchain/miner"
)

var (
	flagMinerPort = flag.Int("port", 5001, "listening port")
	portMinerMin  = 5000
	portMinerMax  = 5010
)

func main() {
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	port := *flagMinerPort

	if port > portMinerMin && port < portMinerMax {
		r := gin.Default()
		miner.InitRouter(r)
		r.Run(fmt.Sprintf(":%d", port))
	} else {
		return
	}
}
