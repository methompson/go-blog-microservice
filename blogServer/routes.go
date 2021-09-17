package blogServer

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"methompson.com/blog-microservice/blogServer/dbController"
)

func (srv *BlogServer) SetRoutes() {
	srv.GinEngine.GET("/", srv.GetHome)

	srv.GinEngine.POST("/add-blog-post", srv.PostAddBlogPost)
	srv.GinEngine.POST("/edit-blog-post", srv.PostEditBlogPost)
	srv.GinEngine.POST("/delete-blog-post", srv.PostDeleteBlogPost)
}

func (srv *BlogServer) GetHome(ctx *gin.Context) {
	token, role, err := srv.GetTokenAndRoleFromHeader(ctx)

	// No Token Error
	if err != nil {
		fmt.Println(err)
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Not Authorized"},
		)
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
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Not Authorized"},
		)
		return
	}

	ctx.Data(200, "text/html; charset=utf-8", make([]byte, 0))
}

func (srv *BlogServer) GetBlogPostBySlug(ctx *gin.Context) {}

func (srv *BlogServer) PostAddBlogPost(ctx *gin.Context) {
	authErr := srv.standardAuthHandler(ctx)

	if authErr != nil {
		return
	}

	// Extract the body
	var body AddBlogBody

	if bindJsonErr := ctx.ShouldBindJSON(&body); bindJsonErr != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "missing required values"},
		)
		return
	}

	id, addBlogErr := srv.BlogController.AddBlogPost(body)

	if addBlogErr != nil {
		switch addBlogErr.(type) {
		case dbController.DuplicateEntryError:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Slug Already Exists"},
			)
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "error adding blog"},
			)
		}
		return
	}

	ctx.JSON(200, gin.H{
		"id": id,
	})
}

func (srv *BlogServer) PostEditBlogPost(ctx *gin.Context) {
	authErr := srv.standardAuthHandler(ctx)

	if authErr != nil {
		return
	}

	var body EditBlogBody

	if bindJsonErr := ctx.ShouldBindJSON(&body); bindJsonErr != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "missing required values"},
		)
		return
	}

	editBlogErr := srv.BlogController.EditBlogPost(body)

	if editBlogErr != nil {
		switch editBlogErr.(type) {
		case dbController.DuplicateEntryError:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Slug Already Exists"},
			)
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "error adding blog"},
			)
		}
		return
	}

	ctx.JSON(200, gin.H{})
}

func (srv *BlogServer) PostDeleteBlogPost(ctx *gin.Context) {
	authErr := srv.standardAuthHandler(ctx)

	if authErr != nil {
		return
	}

	var body DeleteBlogBody

	if bindJsonErr := ctx.ShouldBindJSON(&body); bindJsonErr != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "missing required values"},
		)
		return
	}

	deleteBlogErr := srv.BlogController.DeleteBlogPost(body)

	if deleteBlogErr != nil {
		switch deleteBlogErr.(type) {
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "error adding blog"},
			)
		}
		return
	}

	ctx.JSON(200, gin.H{})
}

func (srv *BlogServer) standardAuthHandler(ctx *gin.Context) error {
	_, role, getTokenErr := srv.GetTokenAndRoleFromHeader(ctx)

	// No Token Error
	if getTokenErr != nil {
		fmt.Println(getTokenErr)
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid token"},
		)
		return getTokenErr
	}

	// Role Error
	if !srv.CanEditBlog(role) {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "not authorized"},
		)
		return errors.New("not authorized")
	}

	return nil
}
