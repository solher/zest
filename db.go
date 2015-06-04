package main

import (
	"fmt"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/ressources"
	"github.com/Solher/zest/usecases"
	"github.com/Solher/zest/utils"
)

func updateDatabase(routes map[usecases.DirectoryKey]usecases.Route) {
	store := infrastructure.NewGormStore()

	err := connectDB(store)
	if err != nil {
		panic("Could not connect to database.")
	}
	defer store.Close()

	fmt.Println("Updating database...")

	err = store.MigrateTables(domain.ModelDirectory.Models)
	if err != nil {
		panic("Could not migrate database.")
	}

	aclRepo := ressources.NewAclRepo(store)
	ressources.NewAclInter(aclRepo).RefreshFromRoutes(routes)
	if err != nil {
		panic("Could not refresh ACLs from routes.")
	}

	fmt.Println("Done.")
}

func resetDatabase(routes map[usecases.DirectoryKey]usecases.Route) {
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

	err = seedDatabase(store, routes)
	if err != nil {
		panic("Could not seed database: " + err.Error())
	}

	fmt.Println("Done.")
}

func seedDatabase(store *infrastructure.GormStore, routes map[usecases.DirectoryKey]usecases.Route) error {
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

	accountRepo := ressources.NewAccountRepo(store)
	accounts, err := accountRepo.Create(accounts)
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

	roleRepo := ressources.NewRoleRepo(store)
	roles, err = roleRepo.Create(roles)
	if err != nil {
		return err
	}

	roleMappings := []domain.RoleMapping{
		{AccountID: 1, RoleID: 1},
	}

	roleMappingRepo := ressources.NewRoleMappingRepo(store)
	roleMappings, err = roleMappingRepo.Create(roleMappings)
	if err != nil {
		return err
	}

	aclRepo := ressources.NewAclRepo(store)
	ressources.NewAclInter(aclRepo).RefreshFromRoutes(routes)
	acls, err := aclRepo.Find(nil, nil)
	if err != nil {
		return err
	}

	aclMappings := make([]domain.AclMapping, len(acls))
	for i, acl := range acls {
		aclMappings[i].AclID = acl.ID
		aclMappings[i].RoleID = 1

		if acl.Ressource == "accounts" {
			aclMappings = append(aclMappings, domain.AclMapping{AclID: acl.ID, RoleID: 5})
		}
	}

	aclMappingRepo := ressources.NewAclMappingRepo(store)
	aclMappings, err = aclMappingRepo.Create(aclMappings)
	if err != nil {
		return err
	}

	return nil
}
