package blogServer

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (srv *BlogServer) SetRoutes() {
	srv.GinEngine.GET("/", srv.GetHome)
}

func (srv *BlogServer) GetHome(ctx *gin.Context) {
	token, role, err := srv.GetTokenAndRoleFromHeader(ctx)

	// No Token Error
	if err != nil {
		fmt.Println(err)
		ctx.Data(401, "text/html; charset=utf-8", make([]byte, 0))
		return
	}

	fmt.Printf("Verified ID Token: %v\n", token)
	fmt.Println("User's role: ", role)

	ctx.Data(200, "text/html; charset=utf-8", make([]byte, 0))
}

func (srv *BlogServer) GetBlogPostById(ctx *gin.Context) {
	_, _, err := srv.GetTokenAndRoleFromHeader(ctx)

	// No Token Error
	if err != nil {
		fmt.Println(err)
		ctx.Data(401, "text/html; charset=utf-8", make([]byte, 0))
		return
	}

	ctx.Data(200, "text/html; charset=utf-8", make([]byte, 0))
}

func (srv *BlogServer) GetBlogPostBySlug(ctx *gin.Context) {}

func (srv *BlogServer) PostAddBlogPost(ctx *gin.Context) {
	_, _, err := srv.GetTokenAndRoleFromHeader(ctx)

	// No Token Error
	if err != nil {
		fmt.Println(err)
		ctx.Data(401, "text/html; charset=utf-8", make([]byte, 0))
		return
	}

	ctx.Data(200, "text/html; charset=utf-8", make([]byte, 0))
}

func (srv *BlogServer) PostEditBlogPost(ctx *gin.Context) {}

func (srv *BlogServer) PostDeleteBlogPost(ctx *gin.Context) {}
