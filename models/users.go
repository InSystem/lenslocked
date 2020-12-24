package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

var (
	ErrorNotFound  = errors.New("models: resourses not found")
	ErrorInvalidID = errors.New("models: invalid ID")
	ErrorPasswordIncorrect = errors.New("models: password incorrect")
)

const userPwPepper = "some-rando-string"

type UserService struct {
	db *gorm.DB
}

// User type
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
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

//Authenticate can be used to authenticate a user with the email and password
func (us *UserService) Authenticate(email,password string) (*User, error){
	foundUser, err := us.ByEmail(email)
	if err != nil{
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash),[]byte(password + userPwPepper))
	if err != nil{
		switch err {
			case bcrypt.ErrMismatchedHashAndPassword: 
				return nil, ErrorPasswordIncorrect
			default:
				return nil, err
		}
	} 

	return foundUser, nil
}

// Create will create the provided user
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""
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

	return us.AutoMigrate()
}

// AutoMigrate automatically atomatically UserService
func (us *UserService) AutoMigrate() error {

	err := us.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}
