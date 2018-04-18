package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"blockchain/miner"
	"fmt"
)

/*http://book.8btc.com/books/6/masterbitcoin2cn/_book/ch09.html
https://github.com/Jeiwan/
*/
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
