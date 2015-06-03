package main

import (
	"fmt"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/ressources"
	"github.com/Solher/zest/utils"
)

func migrateDatabase() {
	store := infrastructure.NewGormStore()

	err := connectDB(store)
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

func resetDatabase() {
	store := infrastructure.NewGormStore()

	err := connectDB(store)
	if err != nil {
		panic("Could not connect to database: " + err.Error())
	}
	defer store.Close()

	fmt.Println("Resetting database...")

	err = store.ResetTables(domain.ModelDirectory.Models)
	if err != nil {
		panic("Could not reset database: " + err.Error())
	}

	err = seedDatabase(store)
	if err != nil {
		panic("Could not seed database: " + err.Error())
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
		},
	}

	accountRepository := ressources.NewAccountRepo(store)
	accounts, err := accountRepository.Create(accounts)
	if err != nil {
		return err
	}

	roles := []domain.Role{
		{Name: "Admin"},
		{Name: "Authenticated"},
		{Name: "Owner"},
		{Name: "Guest"},
		{Name: "Anyone"},
	}

	roleRepository := ressources.NewRoleRepo(store)
	roles, err = roleRepository.Create(roles)
	if err != nil {
		return err
	}

	roleMappings := []domain.RoleMapping{
		{AccountID: 1, RoleID: 1},
	}

	roleMappingRepository := ressources.NewRoleMappingRepo(store)
	roleMappings, err = roleMappingRepository.Create(roleMappings)
	if err != nil {
		return err
	}

	aclMappings := []domain.AclMapping{
		{RoleID: 5, AclID: 1},
		{RoleID: 5, AclID: 2},
		{RoleID: 5, AclID: 3},
		{RoleID: 5, AclID: 4},
		{RoleID: 3, AclID: 5},
		{RoleID: 3, AclID: 6},
		{RoleID: 3, AclID: 7},
		{RoleID: 3, AclID: 8},
		{RoleID: 3, AclID: 9},
		{RoleID: 3, AclID: 10},
		{RoleID: 3, AclID: 11},
		{RoleID: 1, AclID: 1},
		{RoleID: 1, AclID: 2},
		{RoleID: 1, AclID: 3},
		{RoleID: 1, AclID: 4},
		{RoleID: 1, AclID: 5},
		{RoleID: 1, AclID: 6},
		{RoleID: 1, AclID: 7},
		{RoleID: 1, AclID: 8},
		{RoleID: 1, AclID: 9},
		{RoleID: 1, AclID: 10},
		{RoleID: 1, AclID: 11},
		{RoleID: 1, AclID: 12},
		{RoleID: 1, AclID: 13},
		{RoleID: 1, AclID: 14},
		{RoleID: 1, AclID: 15},
		{RoleID: 1, AclID: 16},
		{RoleID: 1, AclID: 17},
		{RoleID: 1, AclID: 18},
		{RoleID: 1, AclID: 19},
		{RoleID: 1, AclID: 20},
		{RoleID: 1, AclID: 21},
		{RoleID: 1, AclID: 22},
		{RoleID: 1, AclID: 23},
		{RoleID: 1, AclID: 24},
		{RoleID: 1, AclID: 25},
	}

	aclMappingRepository := ressources.NewAclMappingRepo(store)
	aclMappings, err = aclMappingRepository.Create(aclMappings)
	if err != nil {
		return err
	}

	return nil
}
