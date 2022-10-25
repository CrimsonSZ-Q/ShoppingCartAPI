package controllers

import (
	//"github.com/gofiber/fiber/v2/middleware/session"
	"shidqi/shoppingcartapi/database"
	"shidqi/shoppingcartapi/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartController struct {
	Db *gorm.DB
}

func InitCartController() *CartController {
	db := database.InitDb()

	db.AutoMigrate(&models.Cart{})

	return &CartController{Db: db}
}

// GET /addtocart/:cartid/products/:productid
func (controller *CartController) AddToCart(c *fiber.Ctx) error {
	params := c.AllParams()

	intCartId, _ := strconv.Atoi(params["cartid"])
	intProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	// Find the product first,
	err := models.FindProductById(controller.Db, &product, intProductId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	err2 := models.FindCartById(controller.Db, &cart, intCartId)
	if err2 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	err3 := models.AddtoCart(controller.Db, &cart, &product)
	if err3 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"message": "Produk berhasil masuk keranjang",
		"Product": product,
	})
}

// GET /:cartid
func (controller *CartController) GetDetailCart(c *fiber.Ctx) error {
	params := c.AllParams()

	intCartId, _ := strconv.Atoi(params["cartid"])

	var cart models.Cart
	err := models.FindCart(controller.Db, &cart, intCartId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"Message": "Shopping Cart",
	})
}
