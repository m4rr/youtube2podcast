package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func runWebServer(tehPod Podcast) {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", tehPod)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
