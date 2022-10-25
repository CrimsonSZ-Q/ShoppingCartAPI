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

// GET /checkout/:userid
func (controller *TransactionController) InsertToTransaction(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaction models.Transaction
	var cart models.Cart

	// Find the product first,
	err := models.FindCart(controller.Db, &cart, intUserId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.CreateTransaction(controller.Db, &transaction, int(intUserId), cart.Products)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Delete products in cart
	errss := models.DeleteProductInChart(controller.Db, cart.Products, &cart, int(intUserId))
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"message": "Transaksi Sukses",
	})
}

// GET /historytransaksi/:userid
func (controller *TransactionController) GetTransaction(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transactions []models.Transaction
	err := models.ViewTransactionById(controller.Db, &transactions, intUserId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"Title":      "History Transaksi",
		"Transaksis": transactions,
	})

}

// GET /history/detail/:transaksiid
func (controller *TransactionController) DetailTransaction(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

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
