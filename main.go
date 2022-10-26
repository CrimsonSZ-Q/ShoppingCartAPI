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
	cartApiController := controllers.InitCartController(store)
	transactionApiController := controllers.InitTransactionController(store)

	//grouping for controller
	p := app.Group("/products")
	c := app.Group("/cart")
	t := app.Group("/transactions")

	app.Get("/accounts", accountApiController.GetAllAccount)
	app.Post("/accounts/create", accountApiController.CreateAccount)
	app.Post("/login", accountApiController.LoginUser)
	app.Post("/logout", accountApiController.Logout)

	p.Get("/", productApiController.GetAllProduct)
	p.Post("/", productApiController.CreateProduct)
	p.Get("/detail/:id", productApiController.GetDetailProduct)
	p.Put("/:id", productApiController.EditProduct)
	p.Delete("/:id", productApiController.DeleteProduct)
	p.Get("/addtocart/:cartid/products/:productid", cartApiController.AddToCart)

	c.Get("/:cartid", cartApiController.GetDetailCart)

	t.Get("/out/:accountid", transactionApiController.InsertToTransaction)
	t.Get("/list/:accountid", transactionApiController.GetTransaction)
	t.Get("/detail/:transactionid", transactionApiController.DetailTransaction)

	app.Listen(":3000")
}
