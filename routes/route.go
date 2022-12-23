package routes

import (
	"blog/controller"
	"blog/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	var apiPrefix = app.Group("/api")
	apiPrefix.Post("/api/register", controller.Register)
	apiPrefix.Post("/api/login", controller.Login)
	apiPrefix.Use(middleware.IsAuthenticate)
	blogPrefix := apiPrefix.Group("/blog/")
	blogPrefix.Post("post", controller.CreatePost)
	blogPrefix.Get("allpost", controller.AllPost)
	blogPrefix.Get("allpost/:id", controller.DetailPost)
	blogPrefix.Put("updatepost/:id", controller.UpdatePost)
	blogPrefix.Get("uniquepost", controller.UniquePost)
	blogPrefix.Delete("deletepost/:id", controller.DeletePost)
	apiPrefix.Post("/upload-image", controller.Upload)
	apiPrefix.Static("/uploads", "./uploads")
}
