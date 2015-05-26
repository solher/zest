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

func reinitDatabase() {
	store := infrastructure.NewGormStore()

	err := connectDB(store)
	if err != nil {
		panic("Could not connect to database: " + err.Error())
	}
	defer store.Close()

	fmt.Println("Reinitializing database...")

	err = store.ReinitTables(domain.ModelDirectory.Models)
	if err != nil {
		panic("Could not reinit database: " + err.Error())
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

	acls := []domain.Acl{
		{Ressource: "accounts", Method: "Signin"},
		{Ressource: "accounts", Method: "Signup"},
		{Ressource: "accounts", Method: "Signout"},
		{Ressource: "accounts", Method: "Current"},
		{Ressource: "users", Method: "Create"},
		{Ressource: "users", Method: "Find"},
		{Ressource: "users", Method: "FindByID"},
		{Ressource: "users", Method: "Upsert"},
		{Ressource: "users", Method: "DeleteAll"},
		{Ressource: "users", Method: "DeleteByID"},
		{Ressource: "sessions", Method: "Create"},
		{Ressource: "sessions", Method: "Find"},
		{Ressource: "sessions", Method: "FindByID"},
		{Ressource: "sessions", Method: "Upsert"},
		{Ressource: "sessions", Method: "DeleteAll"},
		{Ressource: "sessions", Method: "DeleteByID"},
		{Ressource: "roleMappings", Method: "Create"},
		{Ressource: "roleMappings", Method: "Find"},
		{Ressource: "roleMappings", Method: "FindByID"},
		{Ressource: "roleMappings", Method: "Upsert"},
		{Ressource: "roleMappings", Method: "DeleteAll"},
		{Ressource: "roleMappings", Method: "DeleteByID"},
	}

	aclRepository := ressources.NewAclRepo(store)
	acls, err = aclRepository.Create(acls)
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
	}

	aclMappingRepository := ressources.NewAclMappingRepo(store)
	aclMappings, err = aclMappingRepository.Create(aclMappings)
	if err != nil {
		return err
	}

	return nil
}
