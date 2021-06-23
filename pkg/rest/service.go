package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/btowers/blog-go/pkg/lister"
	"github.com/btowers/blog-go/pkg/rest/middleware"

	"github.com/btowers/blog-go/pkg/adder"
	"github.com/btowers/blog-go/pkg/auth"
	"github.com/btowers/blog-go/pkg/remover"
	"github.com/btowers/blog-go/pkg/updater"
	"github.com/gin-gonic/gin"
)

func NewRouter(aut auth.Service, add adder.Service, list lister.Service, remove remover.Service, update updater.Service) *gin.Engine {

	router := gin.Default()

	a := router.Group("/api")
	{
		// User Account Routes
		b := a.Group("/user")
		{
			b.POST("/login", middleware.JWT(aut).LoginHandler)
			b.POST("/register", func(c *gin.Context) {
				var user auth.User
				c.ShouldBindJSON(&user)
				aut.Register(user)
			})
			b.GET("/logout", middleware.JWT(aut).MiddlewareFunc(), middleware.JWT(aut).LogoutHandler)
			b.PUT("/update", middleware.JWT(aut).MiddlewareFunc(), func(c *gin.Context) {
				var user updater.User
				c.ShouldBindJSON(&user)
				update.UpdateUser(user.Email, user)
			})
			b.DELETE("/delete", middleware.JWT(aut).MiddlewareFunc(), func(c *gin.Context) {
				var user remover.User
				c.ShouldBindJSON(&user)
				remove.DeleteUser(user)
			})
		}

		// Blog Post Routes
		p := a.Group("/post", middleware.JWT(aut).MiddlewareFunc())
		{
			p.POST("/", func(c *gin.Context) {
				var post adder.Post
				var user adder.User

				userInterface, _ := c.Get("user")
				userBytes, _ := json.Marshal(userInterface)
				json.Unmarshal(userBytes, &user)
				c.ShouldBindJSON(&post)
				post.Author = adder.User{FirstName: user.FirstName, LastName: user.LastName, Id: user.Id}
				add.AddPost(post)
			})
			p.GET("", func(c *gin.Context) {

				var post lister.Post
				post.Id = c.Query("id")
				//c.ShouldBindJSON(&post)
				a, err := list.GetPost(post)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(a)
				c.JSON(http.StatusOK, a)
			})
			p.PUT("/", func(c *gin.Context) {
				var user updater.User
				c.ShouldBindJSON(&user)
				update.UpdateUser(user.Email, user)
			})
			p.DELETE("/", func(c *gin.Context) {
				var user remover.User
				c.ShouldBindJSON(&user)
				remove.DeleteUser(user)
			})
		}
	}
	/*
		router.GET("/auth/checkAuthentication", middleware.JWT(aut).MiddlewareFunc(), func(c *gin.Context) {

			userClaim, _ := c.Get("user")
			userClaimM, _ := json.Marshal(userClaim)

			var userr bson.M
			json.Unmarshal(userClaimM, &userr)

			c.JSON(http.StatusOK, gin.H{
				"authorized": true,
				"user":       userr,
			})
		})

		router.StaticFile("./favicon.ico", "../dist/favicon.ico")
		router.Static("/static", "../dist/static")
		router.LoadHTMLGlob("../dist/templates/*")

		router.NoRoute(func(c *gin.Context) {
			c.File("../dist/templates/index.html")
		})
	*/
	return router
}
