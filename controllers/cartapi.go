package controllers

import (
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

// GET /shoppingcart/:cartid
func (controller *CartController) GetCart(c *fiber.Ctx) error {
	userId := c.Query("userid")
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Login first",
		})
	}

	var cart []models.Cart
	err = models.ViewCart(controller.Db, &cart, user_id)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(cart)
}

func (controller *CartController) AddtoCart(c *fiber.Ctx) error {
	userId := c.Query("userid")
	productId := c.Query("productid")

	product_id, _ := strconv.Atoi(productId)
	user_id, _ := strconv.Atoi(userId)

	var cart models.Cart
	var product models.Product
	var find models.Cart

	err2 := models.ReadProductById(controller.Db, &product, product_id)

	if err2 != nil {
		c.JSON(fiber.Map{
			"message": "Product Error",
			"status":  500,
		})
	}

	err4 := models.FindCart(controller.Db, &find, product_id, user_id)
	if err4 != nil {
		c.JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  500,
		})
	}

	if find.Id != 0 {
		find.Quantity = find.Quantity + 1
		find.Total = find.Total + product.Price

		models.UpdateCart(controller.Db, &find)

		return c.JSON(find)
	} else {
		cart.ProductId = product_id
		cart.UserId = user_id
		cart.Quantity = 1
		cart.Total = float64(cart.Quantity) * product.Price

		err3 := models.AddtoCart(controller.Db, &cart)

		if err3 != nil {
			c.JSON(fiber.Map{
				"message": "Internal Server Error",
				"status":  500,
			})
		}

		return c.JSON(cart)
	}

}
