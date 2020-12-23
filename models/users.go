package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

var (
	ErrorNotFound  = errors.New("models: resourses not found")
	ErrorInvalidID = errors.New("models: invalid ID")
)

type UserService struct {
	db *gorm.DB
}

// User type
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
	// Orders []Order
}

// NewUserService create connection to the databse
func NewUserService(connectionInfo string) (*UserService, error) {

	// Logger can be Warn Info Error Silent
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

// ByID will look up by the ID provided
func (us *UserService) ByID(id int) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err

}

//ByEmail will looks up a user with a given email address and return that user
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

// Create will create the provided user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

//Update will update the provided user with all the data
//in provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id != 0 {
		user := User{Model: gorm.Model{ID: id}}
		return us.db.Delete(&user).Error
	}
	return ErrorInvalidID
}

// Closes UserServices database connection
func (us *UserService) Close() error {
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

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}
