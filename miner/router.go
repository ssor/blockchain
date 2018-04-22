package miner

import (
    "github.com/gin-gonic/gin"
    "fmt"
)

func initRouter(r *gin.Engine) {
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
}

func NewMinerNode(portMin, portMax int) *MinerNode {
    m := &MinerNode{
        portMin: portMin,
        portMax: portMax,
    }
    return m
}

type MinerNode struct {
    portMin int
    portMax int
}

func (miner *MinerNode) Start() {
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", miner.portMin))
}
