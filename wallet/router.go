package wallet

import (
    "github.com/gin-gonic/gin"
    "fmt"
)

func initRouter(r *gin.Engine) {
    r.GET("/", defaultPage)
    r.GET("/ping", ping)
    r.GET("wallets", wallets)
    r.GET("version", version)
    r.GET("add_wallet", addWallet)

    r.LoadHTMLGlob("wallet/templates/*")
    r.Static("/static", "static")
}

func (wallet *WalletNode) Start() {
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", wallet.port))
}

func NewWalletNode(port int) *WalletNode {
    w := &WalletNode{
        port: port,
    }
    return w
}

type WalletNode struct {
    port int
}
