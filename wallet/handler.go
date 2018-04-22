package wallet

import "github.com/gin-gonic/gin"
import (
    "net/http"
    "blockchain/corelib"
    "github.com/sirupsen/logrus"
)

func version(c *gin.Context) {
    c.JSON(http.StatusOK, nil)
}

func defaultPage(c *gin.Context) {
    c.HTML(http.StatusOK, "default.html", nil)
}

func ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "pong",
    })
}

func wallets(c *gin.Context) {
    wallets, err := corelib.GetAllWallets()
    if err == nil {
        c.JSON(http.StatusOK, wallets)
    } else {
        logrus.Info("get wallets error: ", err)
        c.JSON(http.StatusOK, nil)
    }
}

func addWallet(c *gin.Context) {

    c.JSON(http.StatusOK, nil)

}
