package models //controler akan mengambil model
import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Id       int    `form:"id" json:"id" validate:"required"`
	Name     string `form:"name" json:"name" validate:"required"`
	Email    string `form:"email" json:"email" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	hPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	u.Password = string(hPass)
	return nil
}

func (u *Account) BeforeSave(tx *gorm.DB) (err error) {
	hPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	u.Password = string(hPass)
	return nil

}

func CreateAccount(db *gorm.DB, newAccount *Account) (err error) {
	db.AutoMigrate(&Account{})
	err = db.Create(newAccount).Error
	if err != nil {
		return err
	}
	return nil
}

func FindUserByEmail(db *gorm.DB, account *Account, email string) (err error) {
	err = db.Where("email=?", email).First(account).Error
	if err != nil {
		return err
	}
	return nil
}
