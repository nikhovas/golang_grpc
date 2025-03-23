package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"size:100;not null"`
	Email  string  `gorm:"unique;not null"`
	Orders []Order `gorm:"foreignKey:UserID"`
}

type OrderStatus string

const (
	Accepted   OrderStatus = "ACCEPTED"
	InProgress OrderStatus = "IN_PROGRESS"
	Ready      OrderStatus = "READY"
)

type Order struct {
	ID     uint        `gorm:"primaryKey"`
	Status OrderStatus `gorm:"type:varchar(20);not null"`
	Text   string
	UserID uint `gorm:"not null"`
}

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost dbname=postgres port=5432 sslmode=disable TimeZone=UTC"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}, &Order{})
	if err != nil {
		panic(err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		users := []User{
			{Name: "Alice Smith", Email: "a@a.ru"},
			{Name: "Alice2 Smith", Email: "a2@a.ru"},
			{Name: "Alice3 Smith", Email: "a3@a.ru"},
		}

		if err := db.Create(&users).Error; err != nil {
			panic(err)
		}

		// Create sample orders
		initOrders := []Order{
			{Text: "Order 1 for Alice", Status: Accepted, UserID: users[0].ID},
			{Text: "Order 2 for Alice", Status: InProgress, UserID: users[0].ID},
			{Text: "Order 3 for Alice", Status: Ready, UserID: users[0].ID},
			{Text: "Order 1 for Bob", Status: Ready, UserID: users[1].ID},
			{Text: "Order 2 for Bob", Status: Accepted, UserID: users[1].ID},
			{Text: "Order 1 for Charlie", Status: Ready, UserID: users[2].ID},
		}

		if err := db.Create(&initOrders).Error; err != nil {
			panic(err)
		}
	}

	// Достаем заказы
	var orders []Order
	// Синтаксис на самом деле одинаковый с Create, просто там мы сразу обратились к Error
	result := db.Where("user_id = ?", 1).Find(&orders)
	if result.Error != nil {
		panic("Не получилось достать заказы пользователя 1")
	}

	fmt.Printf("Orders for User ID %d:\n", 1)
	for _, order := range orders {
		fmt.Printf("Order ID: %d, Text: %s, Status: %s\n", order.ID, order.Text, order.Status)
	}

	var usersWithFinishedOrders []User
	// Обычная (полная) предзагрузка
	// err = db.Preload("Orders").Find(&usersWithFinishedOrders).Error
	// Условная
	err = db.Preload("Orders", "status = ?", Ready).Find(&usersWithFinishedOrders).Error
	if err != nil {
		panic(result.Error)
	}

	for _, user := range usersWithFinishedOrders {
		fmt.Println(user.Name)
		for _, order := range user.Orders {
			fmt.Printf("- order_id: %d; text: %s; status: %s\n", order.ID, order.Text, order.Status)
		}
	}
}
