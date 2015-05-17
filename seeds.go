package main

import (
	"fmt"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/infrastructure"
	"github.com/Solher/auth-scaffold/ressources"
	"github.com/Solher/auth-scaffold/utils"
)

func migrateDatabase() {
	store := infrastructure.NewGormStore()

	err := store.Connect("sqlite3", "database.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer store.Close()

	fmt.Println("Migrating database...")

	err = store.MigrateTables(domain.ModelDirectory.Models)
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

	err = store.ReinitTables(domain.ModelDirectory.Models)
	if err != nil {
		panic("Could not reinit database.")
	}

	err = seedDatabase(store)
	if err != nil {
		panic("Could not seed database.")
	}

	fmt.Println("Done.")
}

func seedDatabase(store *infrastructure.GormStore) error {
	fmt.Println("Seeding database...")

	accounts := []domain.Account{
		{
			Users: []domain.User{
				{
					FirstName: "Fabien",
					LastName:  "Herfray",
					Password:  utils.QuickHashPassword("qwertyuiop"),
					Email:     "fabien.herfray@me.com",
				},
			},
			IsAdmin: true,
		},
	}

	accountRepository := ressources.NewAccountRepo(store)
	_, err := accountRepository.Create(accounts)
	if err != nil {
		return err
	}

	return nil
}
