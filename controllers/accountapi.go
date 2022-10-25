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
	var convertpass LoginForm

	if err := c.BodyParser(&account); err != nil {
		return c.SendStatus(400)
	}

	convertpassword, _ := bcrypt.GenerateFromPassword([]byte(convertpass.Password), 10)
	sHash := string(convertpassword)
	account.Password = sHash

	// save account
	err := models.CreateAccount(controller.Db, &account)
	if err != nil {
		return c.SendStatus(500)
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
		return c.SendStatus(500) // Unsuccessful login (cannot find user)
	}

	// Compare password
	compare := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(login.Password))
	status := compare == nil
	if status { // compare == nil artinya hasil compare di atas true
		sess.Set("username", login.Username)
		sess.Set("id", account.Id)
		sess.Save()

		return c.JSON(fiber.Map{
			"message": "Selamat Datang",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data not found",
	})
}

func (controller *AccountAPIController) CheckSession(c *fiber.Ctx) error {
	sess, _ := controller.store.Get(c)
	name := sess.Get("name")
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
