package models

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrorNotFound = errors.New("models: resourses not found")
)

type UserService struct{
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
	// Orders []Order
}

// NewUserService create connection to the databse
func NewUserService(connectionInfo string) (*UserService, error){
db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

// ByID will look up by the ID provided
func (us *UserService) ByID(id int) (*User, error) {
	var u User
	err := us.db.Where("id = ?", id).First(&u).Error

	switch err {
	case nil:
		return &u, error(nil)
	case gorm.ErrRecordNotFound:
		return nil, ErrorNotFound
	default:
		return nil, err
	}

}

// Closes UserServices database connection
func(us *UserService) Close () error {
	db, err := us.db.DB()
	if err != nil {
		panic(err)
	}
	return db.Close()
}

//DestructiveReset drops the all tables and rebuild theme
func (us *UserService) DestructiveReset() error {
	
	err := us.db.Migrator().DropTable(&User{})
	if err != nil {
		return err
	}

	return us.db.AutoMigrate(&User{})
}