package controllers //controler akan mengambil model
import (
	"github.com/go-playground/validator/v10"
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

type RegisterForm struct {
	Name     string `form:"name" json:"name" validate:"required"`
	Username string `form:"username" json:"username" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

var checker = validator.New()

func InitAccountAPIController(s *session.Store) *AccountAPIController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Account{})

	return &AccountAPIController{Db: db, store: s}
}

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
	var myForm RegisterForm
	var cart models.Cart

	if err := c.BodyParser(&myForm); err != nil {
		return c.SendStatus(400)
	}

	errChecker := checker.Struct(myForm)
	if errChecker != nil {
		return c.SendStatus(400)
	}

	errUsername := models.FindUserByUsername(controller.Db, &account, myForm.Username)
	if errUsername != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"message": "Username telah digunakan",
		})
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte(myForm.Password), 10)
	sHash := string(hashPass)
	account.Password = sHash
	account.Username = myForm.Username
	account.Name = myForm.Name
	account.Email = myForm.Email

	// save account
	err := models.CreateAccount(controller.Db, &account)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	//create cart for account
	errCart := models.CreateCart(controller.Db, &cart, account.Id)
	if errCart != nil {
		return c.SendStatus(500) // http 500 internal server error

	}

	// if succeed
	return c.JSON(account)
}

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
		sess.Set("username", login.Username)

		sess.Save()

		return c.JSON(fiber.Map{
			"message": "Selamat Datang",
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
