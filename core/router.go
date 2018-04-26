package core

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "blockchain/corelib"
    "github.com/sirupsen/logrus"
)

func initRouter(r *gin.Engine) {
    r.GET("/", defaultPage)
    r.GET("/ping", ping)
    r.GET("blocks", blocks)
    r.GET("version", version)

    r.LoadHTMLGlob("core/templates/*")
    r.Static("/static", "static")
}

func NewCoreNode(portFrom, portTo int) *CoreNode {
    core := &CoreNode{
        portFrom: portFrom,
        portTo:   portTo,
    }
    return core
}

type CoreNode struct {
    portFrom   int
    portTo     int
    blockchain *corelib.Blockchain
}

func (core *CoreNode) Start() {
    var err error
    if corelib.BlockchainInitialized() == true {
        core.blockchain, err = corelib.NewBlockchain("")
        if err != nil {
            logrus.Errorf("init blockchain error: %s", err)
            return
        }
    } else {
        //create genesis wallet
        w, err := corelib.NewWallet()
        if err != nil {
            logrus.Error("create wallet error: ", err)
            return
        }
        err = w.SaveToFile()
        if err != nil {
            logrus.Error("save wallet to file error: ", err)
            return
        }
        address, err := w.Address()
        if err != nil {
            logrus.Error("get wallet address error: ", err)
            return
        }
        core.blockchain, err = corelib.NewBlockchain(string(address))
        if err != nil {
            logrus.Error("init blockchain error: ", err)
            return
        }
    }
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", core.portFrom))
}
