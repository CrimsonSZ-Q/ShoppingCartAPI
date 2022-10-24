package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Id       int `form:"id" json: "id" validate:"required"`
	UserID   uint
	Products []*Product `gorm:"many2many:transaksi_products;"`
}

func CreateTransaction(db *gorm.DB, newTransaction *Transaction, userId uint, products []*Product) (err error) {
	newTransaction.UserID = userId
	newTransaction.Products = products
	err = db.Create(newTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertProductToTransaction(db *gorm.DB, insertedTransaction *Cart, product *Product) (err error) {
	insertedTransaction.Products = append(insertedTransaction.Products, product)
	err = db.Save(insertedTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadAllProductsInTransaction(db *gorm.DB, transaction *Transaction, id int) (err error) {
	err = db.Where("id=?", id).Preload("Products").Find(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadTransactionById(db *gorm.DB, transactions *[]Transaction, id int) (err error) {
	err = db.Where("user_id=?", id).Find(transactions).Error
	if err != nil {
		return err
	}
	return nil
}
