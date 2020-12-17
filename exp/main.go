package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sveta"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// _, err = db.Exec(`
	// INSERT INTO lenslocked_dev.users(name, email)
	// VALUES($1, $2)`, "sveta shchukina", "shchukina@gmail.com")

	// var id int
	// err = db.QueryRow(`
	// INSERT INTO lenslocked_dev.users(name, email)
	// VALUES($1, $2)
	// RETURNING ID`,
	//  "sveta shchukina", "shchukina@gmail.com").Scan(&id)

	//  fmt.Println(id)

	// var id int
	// var name, email string
	// err = db.QueryRow(`
	// SELECT id, name, email
	// FROM lenslocked_dev.users
	// WHERE id=$1`,
	// 	4).Scan(&id, &name, &email)

	// fmt.Println(name)

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("no rows")
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	// rows, err := db.Query(`
	// SELECT id, name, email
	// FROM lenslocked_dev.users
	// WHERE id=$1`,
	// 	4)

	// rows, err := db.Query(`
	// SELECT id, name, email
	// FROM lenslocked_dev.users`)

	// if err != nil {
	// 	panic(err)
	// }

	// defer rows.Close()
	// for rows.Next() {
	// 	err = rows.Scan(&id, &name, &email)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(name)
	// }

	// type User struct {
	// 	ID int
	// 	Name string
	// 	Email string
	// }

	// var  users []User
	// rows, err := db.Query(`
	// SELECT id, name, email
	// FROM lenslocked_dev.users`)

	// if err != nil {
	// 	panic(err)
	// }

	// defer rows.Close()
	// for rows.Next() {
	// 	var user User
	// 	err = rows.Scan(&user.ID, &user.Name, &user.Email)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	users = append(users, user)

	// 	fmt.Println(users)
	// }

	// if rows.Err != nil {
	// 	// handle error
	// }

	// for i := 0; i < 6; i++ {
	// 	description := fmt.Sprintf("A-Adapter-x%d", i)
	// 	_, err = db.Exec(`INSERT INTO lenslocked_dev.orders (user_id, amount, description)
	// 	VALUES($1, $2, $3)`, i, i*2, description)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	rows, err := db.Query(`SELECT * FROM lenslocked_dev.users
	INNER JOIN lenslocked_dev.orders ON lenslocked_dev.users.id = lenslocked_dev.orders.user_id`)

	for rows.Next() {

		var userID, orderID, amount int
		var email, name, desc string

		if err := rows.Scan(&userID, &name, &email, &orderID, &userID, &amount, &desc); err != nil {
			fmt.Println("userID:", userID, "name:", name, "email:", email, "orderId:", orderID, "amount:", amount, "desc:", desc)
		}
	
	}

	if rows.Err() != nil {
		panic(err)
	}

}
