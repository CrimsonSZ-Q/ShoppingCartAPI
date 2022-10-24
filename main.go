package main

import (
	"github.com/gofiber/fiber/v2"
	"shidqi/shoppingcartapi/controllers"
)

func main() {

	app := fiber.New()

	// controllers
	productApiController := controllers.InitProductAPIController()
	// prodController := controllers.InitProductController()

	p := app.Group("/api")
	p.Get("/hello", productApiController.Greeting)
	p.Get("/products", productApiController.GetAllProduct)
	p.Post("/products", productApiController.CreateProduct)
	p.Get("/products/productdetail", productApiController.GetDetailProduct)
	p.Get("/products/detail/:id", productApiController.GetDetailProduct2)
	p.Put("/products/:id", productApiController.EditProduct)
	p.Delete("/products/:id", productApiController.DeleteProduct)

	app.Listen(":3000")
}
