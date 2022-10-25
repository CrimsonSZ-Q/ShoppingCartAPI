package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId   int
	Products []*Product `gorm:"many2many : cart_products;"`
}

func CreateCart(db *gorm.DB, newCart *Cart, userId int) (err error) {
	newCart.UserId = userId
	err = db.Create(newCart).Error
	if err != nil {
		return err
	}
	return nil
}

func ViewCart(db *gorm.DB, cart *[]Cart, id int) (err error) {
	err = db.Where(&Cart{UserId: id}).Preload("Products").Find(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCart(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where(&Cart{UserId: id}).Preload("User").Preload("Product").First(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCartById(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where(&Cart{UserId: id}).First(cart).Error
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

func DeleteProductInChart(db *gorm.DB, products []*Product, newCart *Cart, userId int) (err error) {
	db.Model(&newCart).Association("Products").Delete(products)

	return nil
}
