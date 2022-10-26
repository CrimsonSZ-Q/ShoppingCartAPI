package models

import (
	"gorm.io/gorm"
)

// CRUD
type Cart struct {
	gorm.Model
	AccountID uint
	Products  []*Product `gorm:"many2many:cart_products;"`
}

func CreateCart(db *gorm.DB, newCart *Cart, accountId uint) (err error) {
	newCart.AccountID = accountId
	err = db.Create(newCart).Error
	if err != nil {
		return err
	}
	return nil
}

func ViewCart(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where("account_id=?", id).Preload("Products").Find(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCartById(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where("account_id=?", id).First(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCart(db *gorm.DB, cart *Cart) (err error) {
	db.Save(cart)

	return nil
}

func AddtoCart(db *gorm.DB, data *Cart, product *Product) (err error) {
	data.Products = append(data.Products, product)
	err = db.Save(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteProductInChart(db *gorm.DB, products []*Product, newCart *Cart, accountId uint) (err error) {
	db.Model(&newCart).Association("Products").Delete(products)

	return nil
}
