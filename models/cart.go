package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	Id       uint     `json:"id" validate:"required"`
	ProductId     int  `form:"productid" json:"productid" validate:"required"`
	UserId     uint  `form:"userid" json:"userid" validate:"required"`
	Quantity int     `form:"quantity" json:"quantity" validate:"required"`
	Total    float64 `form:"total" json:"total" validate:"required"`
	Status     string  `form:"status" json:"status" gorm:"default:process"`
	Product Product `gorm:"foreignkey:ProductId;references:Id"`
	User Account `gorm:"foreignkey:UserId;references:Id"`
}

unc ViewCart(db *gorm.DB, cart *[]Cart, id int) (err error) {
	err = db.Where(&Cart{UserId: id, Status: "process"}).Preload("User").Preload("Product").Find(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCart(db *gorm.DB, cart *Cart, product int, user int) (err error) {
	err = db.Where(&Cart{ProductId: product, UserId: user, Status: "process"}).Preload("User").Preload("Product").First(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCart(db *gorm.DB, cart *Cart) (err error) {
	db.Save(cart)
	
	return nil
}

func AddtoCart(db *gorm.DB, data *Cart) (err error) {
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCart(db *gorm.DB, data *Cart, id int) (err error) {
	err = db.Where("id=?", id).Delete(data).Error
	if err != nil {
		return err
	}
	return nil
}