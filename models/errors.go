package models

import (
	"strings"
)

const (
	//ErrorNotFound is returned when a resource cannot be found
	//in the database
	ErrorNotFound modelErroror = "models: resource not found"

	//ErrorPasswordIncorrect is returned when an invalid password
	//is used when attempting to authenticate a user
	ErrorPasswordIncorrect modelErroror = "models: invalid password provided"

	//ErrorEmailRequired is returned when email address is not provided
	//when create a user
	ErrorEmailRequired modelErroror = "models: email address is required"

	//ErrorEmailInvalid is returned when an invalid format email address
	//is provided when create a user
	ErrorEmailInvalid modelErroror = "models: email address is not valid"

	//ErrorEmailTaken is returned when an email address provided was taken
	//by another user on update and create a user
	ErrorEmailTaken modelErroror = "models: email address is already taken"

	//ErrorPasswordTooShort is returned when an update and create attempted
	//with a user password that is less than 8 characters.
	ErrorPasswordTooShort modelErroror = "models: password must be at least 8 characters long"

	//ErrorPasswordRequired is returned when an update and create attempted
	//with a user password that is empty
	ErrorPasswordRequired modelErroror = "models: password is required"

	//ErrorGalleryTitleRequired return when Title of Gallery is empty
	ErrorGalleryTitleRequired modelErroror = "models: Title is required"

	//ErrorPasswordHashRequired is return when an update and create without
	//password hash
	ErrorPasswordHashRequired modelErroror = "models: password hash is required"

	//ErrorIDInvalid is returned when a invalid ID is provided
	//to a method like Delete
	ErrorIDInvalid modelErroror = "models: ID provided was invalid"

	//ErrorRememberTooShort is return when Remember token string convert to len of bytes
	//at least 32
	ErrorRememberTooShort modelErroror = "models: remember token must be at least 32 bytes"

	//ErrorRememberHashRequired is return when Remember Hash is empty
	ErrorRememberHashRequired modelErroror = "models: remember hash is required"

	//ErrorRememberHashRequired is return when Remember is empty
	ErrorRememberRequired modelErroror = "models: remember is required"

	//ErrororUserIDRequired return when UserID of Gallery is zero
	ErrorUserIDRequired modelErroror = "models: UserID is required"
)

type modelErroror string

func (e modelErroror) Error() string {
	return string(e)
}
func (e modelErroror) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])

	return strings.Join(split, " ")
}

// type privateErroror string

// func (e privateErroror) Error() string {
// 	return string(e)
// }
