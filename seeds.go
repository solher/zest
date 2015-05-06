package main

import (
	"fmt"

	"github.com/Solher/auth-boilerplate-3/infrastructure"
	"github.com/Solher/auth-boilerplate-3/models"
)

func migrateDatabase() {
	store := infrastructure.NewGormStore()

	err := store.Connect("sqlite3", "database.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer store.Close()

	fmt.Println("Migrating database...")

	err = store.MigrateTables([]interface{}{models.User{}, models.Email{}})
	if err != nil {
		panic("Could not migrate database.")
	}

	fmt.Println("Done.")
}

func reinitDatabase() {
	store := infrastructure.NewGormStore()

	err := store.Connect("sqlite3", "database.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer store.Close()

	fmt.Println("Reinitializing database...")

	err = store.ReinitTables([]interface{}{models.User{}, models.Email{}})
	if err != nil {
		panic("Could not reinit database.")
	}

	fmt.Println("Done.")
}

// func seedDatabase() {
// 	fmt.Println("Seeding database...")
//
// 	users := []User{
// 		{
// 			FirstName: "Fabien",
// 			LastName:  "Herfray",
// 			Password:  "qwertyuiop",
// 			Emails: []Email{
// 				{Email: "fabien.herfray@me.com"},
// 			},
// 		},
// 		{
// 			FirstName: "Thomas",
// 			LastName:  "Hourlier",
// 			Password:  "1234",
// 			Emails: []Email{
// 				{Email: "thomas.hourlier@cnode.fr"},
// 				{Email: "hourliert@gmail.com"},
// 			},
// 		},
// 	}
//
// 	for _, user := range users {
// 		GlobalDatabase.Create(&user)
// 		fmt.Println("User created:\n", user)
// 	}
//
// 	fmt.Println("Done.")
// }
