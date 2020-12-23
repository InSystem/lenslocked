package main

import (
	"github.com/InSystem/lenslocked/models"
	"fmt"
	"net/http"

	"github.com/InSystem/lenslocked/controllers"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "sveta"
	password = "postgres"
	dbname   = "lenslocked_dev"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Page not found :(</h2>")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	// user := models.User{
	// 	Name: "alexey ryabov2",
	// 	Email: "alexey@gmail2.com",
	// }

	// if err := us.Create(&user); err != nil {
	// 	panic(err)
	// }

	// user.Email = "alesha@gmail2.com"
	
	if err := us.Delete(uint(6)); err != nil {
		panic(err)
	}


	// userByID, err := us.ByID(int(user.ID))
	// fmt.Println(userByID)

	// userByEmail, err := us.ByEmail("alesha@gmail2.com")
	// fmt.Println(userByEmail)

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(notFound)

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
