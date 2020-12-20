package main

import (
	// "bufio"
	// "os"
	// "strings"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sveta"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
	Orders []Order
}

type Order struct {
	UserID uint
	Amount int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// db, err := gorm.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	defer sqlDB.Close()

	// sqlDB.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{}, &Order{})

	// name, email := getInfo()

	// u := User{
	// 	Name: name,
	// 	Email: email,
	// }

	// if err = sqlDB.Create(&u).Error; err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%v\n", u)

	var u User 
	if err := db.Preload("Orders").First(&u).Error; err != nil {
		panic(err)
	}

	createOrder(db, u, 1001, "fake description #1")
	createOrder(db, u, 1002, "fake description #2")
	// db.Where("id = ?", 3).First(&u)
	// fmt.Println(u)
}

// func getInfo() (name, email string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("Enter name: ")
// 	name, _ = reader.ReadString('\n')

// 	fmt.Println("Enter email: ")
// 	email, _ = reader.ReadString('\n')

// 	name = strings.TrimSpace(name)
// 	email = strings.TrimSpace(email)
// 	return name, email
// }


func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserID: user.ID,
		Amount: amount, 
		Description: desc,
	}).Error

	if err != nil {
		panic(err)
	}
}