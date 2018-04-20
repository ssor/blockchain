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

func StartNode(port int) {
    r := gin.Default()
    initRouter(r)
    r.Run(fmt.Sprintf(":%d", port))
}
