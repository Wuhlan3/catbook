package main

import (
	"gin/controller"
	"gin/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine)* gin.Engine {
	r.Use(middleware.CORSMidddleware())
	r.LoadHTMLFiles("./front_end/upload.html","./front_end/home.html",
		"./front_end/login.html","./front_end/regis.html")
	//页面获取
	//r.GET("/home.html", func(c *gin.Context) {
	//	c.HTML(200, "home.html", nil)
	//})
	//r.GET("/upload.html", func(c *gin.Context) {
	//	c.HTML(200, "upload.html", nil)
	//})
	//r.GET("/login.html", func(c *gin.Context) {
	//	c.HTML(200, "login.html", nil)
	//})
	//r.GET("/regis.html", func(c *gin.Context) {
	//	c.HTML(200, "regis.html", nil)
	//})



	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info",middleware.AuthMiddleware(), controller.Info)

	//categoryRoutes := r.Group("/categories")
	//categoryController := controller.NewCategoryController()
	//categoryRoutes.POST("", categoryController.Create)
	//categoryRoutes.PUT("/:id", categoryController.Update)
	//categoryRoutes.GET("/:id", categoryController.Show)
	//categoryRoutes.DELETE("/:id", categoryController.Delete)
	r.POST("/file", controller.GetFile)
	postRoutes := r.Group("/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	postController := controller.NewPostController()
	postRoutes.POST("", postController.Create)
	postRoutes.POST("/update", postController.Update)
	postRoutes.GET("/show", postController.Show)
	postRoutes.POST("/delete", postController.Delete)
	postRoutes.POST("/page/list", postController.PageList)
	return r
}
