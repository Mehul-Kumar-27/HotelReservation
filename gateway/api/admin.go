package api

import (


	"github.com/gin-gonic/gin"
)

var (
	Admin = NewRouter("admin/", adminRoutes)
)

var (
	adminRoutes = []Route{
		NewRoute("", "GET", helloAdmin),
	}
	
)

func helloAdmin(ctx *gin.Context) {
	
	ctx.IndentedJSON(200, gin.H{"body": "Hello admin"})
}
