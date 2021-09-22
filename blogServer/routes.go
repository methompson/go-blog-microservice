package blogServer

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"methompson.com/blog-microservice/blogServer/dbController"
)

func (srv *BlogServer) SetRoutes() {
	srv.GinEngine.GET("/", srv.GetHome)

	srv.GinEngine.GET("/blog", srv.GetBlogPostsByFirstPage)
	srv.GinEngine.GET("/blog/page/:page", srv.GetBlogPostsByPage)
	srv.GinEngine.GET("/blog/id/:id", srv.GetBlogPostById)
	srv.GinEngine.GET("/blog/post/:slug", srv.GetBlogPostBySlug)

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

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", make([]byte, 0))
}

func (srv *BlogServer) GetBlogPostsByPage(ctx *gin.Context) {
	page := ctx.Param("page")

	// Not sure this will ever happen
	if len(page) == 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid page number",
			},
		)

		return
	}

	pageNum, pageNumErr := strconv.Atoi(page)

	if pageNumErr != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid page number",
			},
		)

		return
	}

	srv.GetBlogPosts(ctx, pageNum)
}

func (srv *BlogServer) GetBlogPostsByFirstPage(ctx *gin.Context) {
	srv.GetBlogPosts(ctx, 1)
}

func (srv *BlogServer) GetBlogPosts(ctx *gin.Context, page int) {
	pagination := ctx.Query("pagination")

	paginationNum, paginationNumErr := strconv.Atoi(pagination)
	if paginationNumErr != nil {
		paginationNum = -1
	}

	posts, getPostsErr := srv.BlogController.GetBlogPosts(page, paginationNum)

	if getPostsErr != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "error retrieving blog posts",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		posts,
	)
}

func (srv *BlogServer) GetBlogPostById(ctx *gin.Context) {
	id := ctx.Param("id")

	// Not sure this will ever happen
	if len(id) == 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid id",
			},
		)

		return
	}

	getBlog, getBlogErr := srv.BlogController.GetBlogPostById(id)

	if getBlogErr != nil {
		switch getBlogErr.(type) {
		case dbController.NoResultsError:
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "page does not exist"},
			)
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": getBlogErr.Error()},
			)
		}

		return
	}

	ctx.JSON(
		http.StatusOK,
		getBlog.GetMap(),
	)
}

func (srv *BlogServer) GetBlogPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	// Not sure this will ever happen
	if len(slug) == 0 {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid slug",
			},
		)

		return
	}

	getBlog, getBlogErr := srv.BlogController.GetBlogPostBySlug(slug)

	if getBlogErr != nil {
		switch getBlogErr.(type) {
		case dbController.NoResultsError:
			ctx.JSON(
				http.StatusNotFound,
				gin.H{"error": "page does not exist"},
			)
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": getBlogErr.Error()},
			)
		}

		return
	}

	ctx.JSON(
		http.StatusOK,
		getBlog.GetMap(),
	)
}

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

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"id": id,
		},
	)
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

	ctx.JSON(http.StatusOK, gin.H{})
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
		case dbController.InvalidInputError:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid id. blog does not exist. no blog post deleted"},
			)
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "error deleting blog"},
			)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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
