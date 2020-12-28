package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/InSystem/lenslocked/rand"

	"github.com/InSystem/lenslocked/hash"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

var (
	ErrorNotFound             = errors.New("models: resourses not found")
	ErrorIDInvalid            = errors.New("models: invalid ID")
	ErrorPasswordIncorrect    = errors.New("models: password incorrect")
	ErrorEmailRequired        = errors.New("models: email is required")
	ErroroEmailInvalid        = errors.New("models: email is invalid")
	ErrorEmailTaken           = errors.New("models: email is already taken")
	ErrorPasswordRequired     = errors.New("models: password is reauired")
	ErrorPasswordTooShort     = errors.New("models: password should be at least 8 characters long")
	ErrorPasswordHashRequired = errors.New("models: password hash is required")
	ErroroRememberTooShort    = errors.New("models: remember token must be at least 32 bytes")
	ErrorRememberRequired     = errors.New("models: remember is required")
)

const userPwPepper = "some-random-string"
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
	db *gorm.DB
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

// NewUserService create connection to the databse
func NewUserGorm(connectionInfo string) (*userGorm, error) {

	// Logger can be Warn Info Error Silent
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &userGorm{
		db: db,
	}, nil
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac, userPwPepper)

	return &userService{
		UserDB: uv,
	}, nil
}

// ByID will look up by the ID provided
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err

}

//ByEmail will looks up a user with a given email address and return that user.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up  a user with a given remember token and returns that user.
// This  method expects the token that already been hashed
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
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
	err = bcrypt.CompareHashAndPassword([]byte(founduger.PasswordHash), []byte(password+userPwPepper))
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
	return ug.db.Create(user).Error
}

//Update will update the provided uger with all the data
//in provided uger object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//Delete will delete the uger with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
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

func newUserValidator(udb UserDB, hmac hash.HMAC, pepper string) *userValidator {
	return &userValidator{
		UserDB:      udb,
		hmac:        hmac,
		emailRegexp: regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
		pepper:      pepper,
	}
}

type userValidator struct {
	UserDB
	hmac        hash.HMAC
	emailRegexp *regexp.Regexp
	pepper      string
}

// ByRemember hashes the given token and then call ByRemember
// on subsequent layer UserDB.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValidatorFunction(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// Create will create the provided uger
func (uv *userValidator) Create(user *User) error {
	err := runUserValidatorFunction(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {
	err := runUserValidatorFunction(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Delete will delete the uger with the provided ID
func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	err := runUserValidatorFunction(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// ByEmail will normalize the email address before calling
// ByEmail on the DB layer
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValidatorFunction(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

type userValidatorFunction func(*User) error

func runUserValidatorFunction(user *User, fns ...userValidatorFunction) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash a user's password with a predefined
// pepper (userPepper) and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}

	if len(user.Password) < 8 {
		return ErrorPasswordTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrorRememberRequired
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrorPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrorPasswordHashRequired
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		return nil
	}

	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErroroRememberTooShort
	}

	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValidatorFunction {
	return userValidatorFunction(func(user *User) error {
		if user.ID <= n {
			return ErrorIDInvalid
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrorEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email != "" && !uv.emailRegexp.MatchString(user.Email) {
		return ErroroEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrorNotFound {
		// Email adress is not taken
		return nil
	}
	if err != nil {
		return err
	}

	if user.ID != existing.ID {
		return ErrorEmailTaken
	}

	return nil
}
