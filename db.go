package zest

import (
	"errors"
	"fmt"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/infrastructure"
	"github.com/solher/zest/ressources"
	"github.com/solher/zest/usecases"
)

func updateDatabase(z *Zest) error {
	type dependencies struct {
		Store    *infrastructure.GormStore
		AclInter *ressources.AclInter
		RouteDir *usecases.RouteDirectory
	}

	d := &dependencies{}
	err := z.Injector.Get(d)
	if err != nil {
		return err
	}

	if z.DatabaseURL != "" {
		err = d.Store.Connect("postgres", z.DatabaseURL)
	} else {
		err = d.Store.Connect("sqlite3", "database.db")
	}
	if err != nil {
		return err
	}

	fmt.Println("Updating database...")

	err = d.Store.MigrateTables(domain.ModelDirectory.Models)
	if err != nil {
		return errors.New("Could not migrate database: " + err.Error())
	}

	err = d.AclInter.RefreshFromRoutes(d.RouteDir.Routes)
	if err != nil {
		return errors.New("Could not refresh ACLs from routes: " + err.Error())
	}

	d.Store.Close()

	fmt.Println("Done.")

	return nil
}

func reinitDatabase(z *Zest) error {
	type dependencies struct {
		Store *infrastructure.GormStore
	}

	d := &dependencies{}
	err := z.Injector.Get(d)
	if err != nil {
		return err
	}

	if z.DatabaseURL != "" {
		err = d.Store.Connect("postgres", z.DatabaseURL)
	} else {
		err = d.Store.Connect("sqlite3", "database.db")
	}
	if err != nil {
		return err
	}

	fmt.Println("Reinitializing database...")

	err = d.Store.ResetTables(domain.ModelDirectory.Models)
	if err != nil {
		return errors.New("Could not reinit database: " + err.Error())
	}

	d.Store.Close()

	fmt.Println("Done.")

	return nil
}

func seedDatabase(z *Zest) error {
	type dependencies struct {
		Store             *infrastructure.GormStore
		AccountGuestInter *ressources.AccountGuestInter
		RoleRepo          *ressources.RoleRepo
		RoleMappingRepo   *ressources.RoleMappingRepo
		AclInter          *ressources.AclInter
		AclMappingRepo    *ressources.AclMappingRepo
		RouteDir          *usecases.RouteDirectory
	}

	d := &dependencies{}
	err := z.Injector.Get(d)
	if err != nil {
		return err
	}

	if z.DatabaseURL != "" {
		err = d.Store.Connect("postgres", z.DatabaseURL)
	} else {
		err = d.Store.Connect("sqlite3", "database.db")
	}
	if err != nil {
		return err
	}

	fmt.Println("Seeding database...")

	user := &domain.User{
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "admin",
		Email:     "admin",
	}

	account, err := d.AccountGuestInter.Signup(user)
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

	roles, err = d.RoleRepo.Create(roles)
	if err != nil {
		return err
	}

	roleMappings := []domain.RoleMapping{
		{AccountID: account.ID, RoleID: roles[0].ID},
	}

	roleMappings, err = d.RoleMappingRepo.Create(roleMappings)
	if err != nil {
		return err
	}

	err = d.AclInter.RefreshFromRoutes(d.RouteDir.Routes)
	if err != nil {
		return err
	}

	d.Store.Close()

	return nil
}
