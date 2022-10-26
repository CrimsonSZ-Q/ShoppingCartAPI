package controllers

import (
	"shidqi/shoppingcartapi/database"
	"shidqi/shoppingcartapi/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type TransactionController struct {
	// Declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitTransactionController(s *session.Store) *TransactionController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Transaction{})

	return &TransactionController{Db: db, store: s}
}

// GET /out/:userid
func (controller *TransactionController) InsertToTransaction(c *fiber.Ctx) error {
	params := c.AllParams()

	intAccountId, _ := strconv.Atoi(params["accountid"])

	var transaction models.Transaction
	var cart models.Cart

	// Find cart by id
	err := models.FindCartById(controller.Db, &cart, intAccountId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	//see cart
	err1 := models.ViewCart(controller.Db, &cart, intAccountId)
	if err1 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// cart with 0 product
	if len(cart.Products) == 0 {
		return c.JSON(fiber.Map{
			"message": "Cart kosong, Silahkan tambahkan produk kedalam cart",
		})
	}

	err2 := models.CreateTransaction(controller.Db, &transaction, uint(intAccountId), cart.Products)
	if err2 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Delete products in cart
	err3 := models.DeleteProductInChart(controller.Db, cart.Products, &cart, uint(intAccountId))
	if err3 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"message": "Transaksi Sukses",
	})
}

// GET /transaction/list/:transactionid
func (controller *TransactionController) GetTransaction(c *fiber.Ctx) error {
	params := c.AllParams()

	intAccountId, _ := strconv.Atoi(params["accountid"])

	var transactions []models.Transaction
	err := models.ViewTransactionById(controller.Db, &transactions, intAccountId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"Title":        "History Transaksi",
		"Transactions": transactions,
	})

}

// GET /transaction/detail/:transaksiid
func (controller *TransactionController) DetailTransaction(c *fiber.Ctx) error {
	params := c.AllParams()

	intTransactionId, _ := strconv.Atoi(params["transactionid"])

	var transaction models.Transaction
	err := models.ViewTransaction(controller.Db, &transaction, intTransactionId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"Title":    "History Transaksi",
		"Products": transaction.Products,
	})
}
