package models //controler akan mengambil model
import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name        string `form:"name" json:"name" validate:"required"`
	Username    string `form:"username" json:"username" validate:"required"`
	Email       string `form:"email" json:"email" validate:"required"`
	Password    string `form:"password" json:"password" validate:"required"`
	Cart        Cart
	Transaction []Transaction
}

func CreateAccount(db *gorm.DB, newAccount *Account) (err error) {
	//db.AutoMigrate(&Account{})
	err = db.Create(newAccount).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadAccounts(db *gorm.DB, accounts *[]Account) (err error) {
	err = db.Find(accounts).Error
	if err != nil {
		return err
	}
	return nil
}

func FindUserByUsername(db *gorm.DB, account *Account, username string) (err error) {
	err = db.Where(&Account{Username: username}).First(account).Error
	if err != nil {
		return err
	}
	return nil
}
