package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "sveta"
		password = "postgres"
		dbname   = "lenslocked_test"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// TODO silent logger
	us, err := NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := us.DestructiveReset(); err != nil {
		panic(err)
	}

	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Khanh Minh",
		Email: "khanhminh@gmail.com",
	}
	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected >0. Recieved ID:%d", user.ID)
	}
	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected Created At to be recent. Recieved %s", user.CreatedAt)
	}

	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected Created At to be recent. Recieved %s", user.UpdatedAt)
	}

}
