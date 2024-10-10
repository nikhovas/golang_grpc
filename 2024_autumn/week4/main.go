package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gortm:"size:100;not null"`
	Email  string  `gorm:"unique;not null"`
	Orders []Order `gorm:"foreignKey:UserID"`
}

type OrderStatus string

const (
	Accepted   OrderStatus = "accepted"
	InProgress OrderStatus = "in_progress"
	Ready      OrderStatus = "ready"
)

type Order struct {
	ID     uint64      `gorm:"primarykey"`
	Text   string      `gorm:"type:text;not null"`
	Status OrderStatus `gorm:"type:varchar(20);not null"`
	UserID uint        `gorm:"not null"`
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
	if count > 0 {
		fmt.Println("Не нужно вставлять никаких новых данных")
	} else {
		// Create sample users
		users := []User{
			{Name: "Alice Smith", Email: "alice@example.com"},
			{Name: "Bob Johnson", Email: "bob@example.com"},
			{Name: "Charlie Davis", Email: "charlie@example.com"},
		}

		if err := db.Create(&users).Error; err != nil {
			log.Fatalf("Failed to insert users: %v", err)
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
			log.Fatalf("Failed to insert orders: %v", err)
		}

		fmt.Println("Данные созданы")
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
	result = db.Preload("Orders", "status = ?", Ready).Find(&usersWithFinishedOrders)
	// Закомментирован вариант для предзагрузки без условия
	// result := db.Preload("Orders").Find(&usersWithFinishedOrders)
	if result.Error != nil {
		panic(result.Error)
	}

	for _, user := range usersWithFinishedOrders {
		fmt.Println(user.Name)
		for _, order := range user.Orders {
			fmt.Printf("- order_id: %d; text: %s; status: %s\n", order.ID, order.Text, order.Status)
		}
	}
}
