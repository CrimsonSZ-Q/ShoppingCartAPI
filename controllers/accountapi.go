package controllers //controler akan mengambil model
import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"shidqi/shoppingcartapi/database"
	"shidqi/shoppingcartapi/models"
	//"strconv"
	//"fmt"
)

type AccountAPIController struct {
	// declare variables
	Db    *gorm.DB
	store *session.Store
}

type LoginForm struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

func InitAccountAPIController(s *session.Store) *AccountAPIController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Account{})

	return &AccountAPIController{Db: db, store: s}
}

//View all accounts
func (controller *AccountAPIController) GetAllAccount(c *fiber.Ctx) error {
	// load all products
	var accounts []models.Account
	err := models.ReadAccounts(controller.Db, &accounts)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(accounts)
}

// POST /accounts/create
func (controller *AccountAPIController) CreateAccount(c *fiber.Ctx) error {

	var account models.Account
	var cart models.Cart

	if err := c.BodyParser(&account); err != nil {
		return c.SendStatus(400)
	}

	errUsername := models.FindUserByUsername(controller.Db, &account, account.Username)
	if errUsername != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"message": "Username telah digunakan",
		})
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte(account.Password), 10)
	sHash := string(hashPass)
	account.Password = sHash

	// save account
	err := models.CreateAccount(controller.Db, &account)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	//create cart for account
	errCart := models.CreateCart(controller.Db, &cart, account.ID)
	if errCart != nil {
		return c.SendStatus(500) // http 500 internal server error

	}

	// if succeed
	return c.JSON(account)
}

//Login 
func (controller *AccountAPIController) LoginUser(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}

	var account models.Account
	var login LoginForm

	if err := c.BodyParser(&login); err != nil {
		return c.SendStatus(500)
	}

	// Find user
	err2 := models.FindUserByUsername(controller.Db, &account, login.Username)
	if err2 != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Compare the password
	comparePass := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(login.Password))
	status := comparePass == nil
	if status {
		sess.Set("username", account.Username)
		sess.Set("accountId", account.ID)

		sess.Save()

		return c.JSON(fiber.Map{
			"message":  "Selamat Datang",
			"username": account.Username,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Akun tidak ditemukan",
	})
}

func (controller *AccountAPIController) CheckSession(c *fiber.Ctx) error {
	sess, _ := controller.store.Get(c)
	name := sess.Get("username")
	id := sess.Get("id")

	return c.JSON(fiber.Map{
		"message": "Terimakasih",
		"Id":      id,
		"Nama":    name,
	})
}

//Logout
func (controller *AccountAPIController) Logout(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Destroy()

	return c.JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}
