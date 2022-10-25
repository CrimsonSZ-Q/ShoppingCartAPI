package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Id       int `form:"id" json: "id" validate:"required"`
	UserId   int
	Products []*Product `gorm:"many2many:transaction_products;"`
}

func CreateTransaction(db *gorm.DB, newTransaction *Transaction, userId int, products []*Product) (err error) {
	newTransaction.UserId = userId
	newTransaction.Products = products
	err = db.Create(newTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func AddProductToTransaction(db *gorm.DB, insertedTransaction *Cart, product *Product) (err error) {
	insertedTransaction.Products = append(insertedTransaction.Products, product)
	err = db.Save(insertedTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func ViewTransaction(db *gorm.DB, transaction *Transaction, id int) (err error) {
	err = db.Where("id=?", id).Preload("Products").Find(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func ViewTransactionById(db *gorm.DB, trans *[]Transaction, id int) (err error) {
	err = db.Where("user_id = ?", id).Find(trans).Error
	if err != nil {
		return err
	}
	return nil
}
