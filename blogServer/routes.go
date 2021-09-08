package blogServer

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (srv *BlogServer) SetRoutes() {
	srv.GinEngine.GET("/", srv.GetHome)
}

func (srv *BlogServer) GetHome(ctx *gin.Context) {
	var header AuthorizationHeader

	// No Token Error
	if headerErr := ctx.ShouldBindHeader(&header); headerErr != nil {
		fmt.Println(headerErr)
		ctx.Data(401, "text/html; charset=utf-8", make([]byte, 0))
		return
	}

	srv.ValidateIdToken(header)

	ctx.Data(200, "text/html; charset=utf-8", make([]byte, 0))
}
