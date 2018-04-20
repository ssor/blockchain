package core

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func blocks(c *gin.Context) {
    c.JSON(http.StatusOK, nil)
}

func version(c *gin.Context) {
    c.JSON(http.StatusOK, nil)
}

func defaultPage(c *gin.Context) {
    c.HTML(http.StatusOK, "default.html", nil)
}

func ping(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
}
