package models

import (
	"errors"

	"github.com/InSystem/lenslocked/rand"

	"github.com/InSystem/lenslocked/hash"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

var (
	ErrorNotFound          = errors.New("models: resourses not found")
	ErrorInvalidID         = errors.New("models: invalid ID")
	ErrorPasswordIncorrect = errors.New("models: password incorrect")
)

const ugerPwPepper = "some-rando-string"
const hmacSecretKey = "secret-hmac-key"

// ugerDB is uged to interact with the users database
type UserDB interface {
	//Query for single user
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Methods for altering uger
	Create(user *User) error
	Delete(id uint) error
	Update(user *User) error

	// uged to close DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

//UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	//Authenticate will verify the provided email and password are correct. If they
	//are correct, the user corresponding to that email will return. Otherwise You
	//will receive either: ErrNotFound, ErrPasswordIncorrect or other error if something
	// goes wrong.
	Authenticate(email, password string) (*User, error)

	UserDB
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// If the userGorm stops mathing the inteface of UserDb ide starts complaing
// So we can sure that userGorm match that interface
var _ UserDB = &userGorm{}
var _ UserDB = &userValidator{}
var _ UserService = &userService{}


// uger type
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;uniqueIndex"`
}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

// NewUserService create connection to the databse
func NewUserGorm(connectionInfo string) (*userGorm, error) {

	// Logger can be Warn Info Error Silent
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)

	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

// ByID will look up by the ID provided
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err

}

//ByEmail will looks up a user with a given email address and return that uger
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up  a uger with a given remember token and returns that uger
// This method will handle hashing for ug.
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

//Authenticate can be uged to authenticate a uger with the email and password
func (us *userService) Authenticate(email, password string) (*User, error) {
	founduger, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(founduger.PasswordHash), []byte(password+ugerPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrorPasswordIncorrect
		default:
			return nil, err
		}
	}

	return founduger, nil
}

// Create will create the provided uger
func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + ugerPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token
	}

	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

//Update will update the provided uger with all the data
//in provided uger object
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

//Delete will delete the uger with the provided ID
func (ug *userGorm) Delete(id uint) error {
	if id != 0 {
		user := User{Model: gorm.Model{ID: id}}
		return ug.db.Delete(&user).Error
	}
	return ErrorInvalidID
}

// Closes userGorms database connection
func (ug *userGorm) Close() error {
	db, err := ug.db.DB()
	if err != nil {
		panic(err)
	}
	return db.Close()
}

//DestructiveReset drops the all tables and rebuild theme
func (ug *userGorm) DestructiveReset() error {

	err := ug.db.Migrator().DropTable(&User{})
	if err != nil {
		return err
	}

	return ug.AutoMigrate()
}

// AutoMigrate automatically atomatically userGorm
func (ug *userGorm) AutoMigrate() error {

	err := ug.db.AutoMigrate(&User{})
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
