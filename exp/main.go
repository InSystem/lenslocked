package main

import (
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
	Email string `gorm:not null;uniqueIndex;`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// db, err := gorm.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }


	sqlDB, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db, err := sqlDB.DB()

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	sqlDB.Migrator().DropTable(&User{})
	// sqlDB.AutoMigrate(&User{})
}
