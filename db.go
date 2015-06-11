package zest

import (
	"errors"
	"fmt"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/infrastructure"
	"github.com/Solher/zest/ressources"
	"github.com/Solher/zest/usecases"
)

func updateDatabase(z *Zest) error {
	type dependencies struct {
		Store    *infrastructure.GormStore
		AclInter *ressources.AclInter
		RouteDir *usecases.RouteDirectory
	}

	d := &dependencies{}
	err := z.injector.Get(d)
	if err != nil {
		return err
	}

	fmt.Println("Updating database...")

	err = d.Store.MigrateTables(domain.ModelDirectory.Models)
	if err != nil {
		return errors.New("Could not migrate database.")
	}

	d.AclInter.RefreshFromRoutes(d.RouteDir.Routes())
	if err != nil {
		return errors.New("Could not refresh ACLs from routes.")
	}

	fmt.Println("Done.")

	return nil
}

func resetDatabase(z *Zest) error {
	type dependencies struct {
		Store *infrastructure.GormStore
	}

	d := &dependencies{}
	err := z.injector.Get(d)
	if err != nil {
		return err
	}

	fmt.Println("Resetting database...")

	err = d.Store.ResetTables(domain.ModelDirectory.Models)
	if err != nil {
		return errors.New("Could not reset database: " + err.Error())
	}

	fmt.Println("Done.")

	return nil
}

func seedDatabase(z *Zest) error {
	type dependencies struct {
		Store            *infrastructure.GormStore
		AccountInter     *ressources.AccountInter
		RoleInter        *ressources.RoleInter
		RoleMappingInter *ressources.RoleMappingInter
		AclInter         *ressources.AclInter
		AclMappingInter  *ressources.AclMappingInter
		RouteDir         *usecases.RouteDirectory
	}

	d := &dependencies{}
	err := z.injector.Get(d)
	if err != nil {
		return err
	}

	fmt.Println("Seeding database...")

	user := &domain.User{
		FirstName: "Fabien",
		LastName:  "Herfray",
		Password:  "qwertyuiop",
		Email:     "fabien.herfray@me.com",
	}

	account, err := d.AccountInter.Signup(user)
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

	roles, err = d.RoleInter.Create(roles)
	if err != nil {
		return err
	}

	roleMappings := []domain.RoleMapping{
		{AccountID: account.ID, RoleID: roles[0].ID},
	}

	roleMappings, err = d.RoleMappingInter.Create(roleMappings)
	if err != nil {
		return err
	}

	d.AclInter.RefreshFromRoutes(d.RouteDir.Routes())
	acls, err := d.AclInter.Find(usecases.QueryContext{})
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

	aclMappings, err = d.AclMappingInter.Create(aclMappings)
	if err != nil {
		return err
	}

	return nil
}
