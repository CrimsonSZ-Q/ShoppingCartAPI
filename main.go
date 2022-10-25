package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"shidqi/shoppingcartapi/controllers"
)

func main() {
	store := session.New()
	app := fiber.New()

	// controllers
	productApiController := controllers.InitProductAPIController()
	accountApiController := controllers.InitAccountAPIController(store)
	cartApiController := controllers.InitCartController()
	transactionApiController := controllers.InitTransactionController(store)

	p := app.Group("/products")
	c := app.Group("/Cart")
	t := app.Group("/Transactions")

	p.Get("/", productApiController.GetAllProduct)
	p.Post("/", productApiController.CreateProduct)
	p.Get("/productdetail", productApiController.GetDetailProduct)
	p.Get("/detail/:id", productApiController.GetDetailProduct2)
	p.Put("/:id", productApiController.EditProduct)
	p.Delete("/:id", productApiController.DeleteProduct)
	p.Get("/addtocart/:cartid/products/:productid", cartApiController.AddToCart)

	app.Get("/accounts", accountApiController.GetAllAccount)
	app.Post("/accounts/create", accountApiController.CreateAccount)
	app.Post("/login", accountApiController.LoginUser)

	c.Get("/:userid", cartApiController.GetDetailCart)

	t.Get("/:userid", transactionApiController.InsertToTransaction)

	app.Listen(":3000")
}
