package core

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

func StartNode(port int) {
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", port))
}
