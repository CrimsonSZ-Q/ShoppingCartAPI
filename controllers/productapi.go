package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"shidqi/shoppingcartapi/database"
	"shidqi/shoppingcartapi/models"
	"strconv"
)

type ProductAPIController struct {
	// declare variables
	Db *gorm.DB
}

func InitProductAPIController() *ProductAPIController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Product{})

	return &ProductAPIController{Db: db}
}

func (controller *ProductAPIController) Greeting(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Selamat Datang...",
	})
}

// routing
// GET /products
func (controller *ProductAPIController) GetAllProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(products)
}

// POST /products/create
func (controller *ProductAPIController) CreateProduct(c *fiber.Ctx) error {

	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.SendStatus(400)
	}
	// save product
	err := models.CreateProduct(controller.Db, &product)
	if err != nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(product)
}

// GET /products/productdetail?id=xxx
func (controller *ProductAPIController) GetDetailProduct(c *fiber.Ctx) error {
	id := c.Query("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.FindProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(product)
}

// GET /products/detail/xxx
func (controller *ProductAPIController) GetDetailProduct2(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.FindProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(product)
}

// / PUT products/xx
func (controller *ProductAPIController) EditProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.FindProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var updateProduct models.Product

	if err := c.BodyParser(&updateProduct); err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	product.Name = updateProduct.Name
	product.Quantity = updateProduct.Quantity
	product.Price = updateProduct.Price

	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			product.Image = fmt.Sprintf("public/upload/%s", file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("public/upload/%s", file.Filename)); err != nil {
				return err
			}
		}
	}
	// save product
	models.UpdateProduct(controller.Db, &product)

	return c.JSON(fiber.Map{
		"message": "Berhasil Mengubah Product",
	})

}

// / DELETE /products/:id
func (controller *ProductAPIController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)
	return c.JSON(fiber.Map{
		"message": "Data was deleted",
	})

}
