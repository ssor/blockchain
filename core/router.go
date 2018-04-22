package core

import (
    "github.com/gin-gonic/gin"
    "fmt"
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
    portFrom int
    portTo   int
}

func (core *CoreNode) Start() {
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", core.portFrom))
}
