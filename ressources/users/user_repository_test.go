// Generated by: main
// TypeWriter: repository_test
// Directive: +gen on User

package users

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/ressources/emails"
	"github.com/Solher/auth-scaffold/utils"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

type stubGormStore struct {
	db *gorm.DB
}

func (st *stubGormStore) Connect(adapter, url string) error {
	return nil
}

func (st *stubGormStore) Close() error {
	return nil
}

func (st *stubGormStore) GetDB() *gorm.DB {
	return st.db
}

func (st *stubGormStore) MigrateTables(tables []interface{}) error {
	return nil
}

func (st *stubGormStore) ReinitTables(tables []interface{}) error {
	return nil
}

func (st *stubGormStore) BuildQuery(filter *interfaces.Filter) (*gorm.DB, error) {
	query := st.db

	gormFilter := &interfaces.GormFilter{}

	if filter != nil {
		dbNamedFields := []string{}
		fields := filter.Fields

		for _, field := range fields {
			dbNamedFields = append(dbNamedFields, utils.ToDBName(field))
		}

		if filter.Order != "" {
			split := strings.Split(filter.Order, " ")
			filter.Order = utils.ToDBName(split[0]) + " " + split[1]
		}

		gormFilter.Fields = dbNamedFields
		gormFilter.Limit = filter.Limit
		gormFilter.Offset = filter.Offset
		gormFilter.Order = filter.Order
		gormFilter.Where = "first_name IN ('Robert', 'Fabien') OR (last_name = 'Dupont')"
	}

	if len(gormFilter.Fields) != 0 {
		query = query.Select(gormFilter.Fields)
	}

	if gormFilter.Offset != 0 {
		query = query.Offset(gormFilter.Offset)
	}

	if gormFilter.Limit != 0 {
		query = query.Limit(gormFilter.Limit)
	}

	if gormFilter.Order != "" {
		query = query.Order(gormFilter.Order)
	}

	if gormFilter.Where != "" {
		query = query.Where(gormFilter.Where)
	}

	return query, nil
}

func initDatabase() (interfaces.GormStore, error) {
	store := &stubGormStore{}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return store, err
	}

	err = db.AutoMigrate(&User{}, &emails.Email{}).Error
	if err != nil {
		return store, err
	}

	store.db = &db

	return store, nil
}

func TestRepository(t *testing.T) {
	store, err := initDatabase()
	if err != nil {
		panic(fmt.Sprintf("Error initializing database: %v", err))
	}
	repo := NewRepository(store)

	Convey("Testing users repository...", t, func() {
		Convey("Should be able to create users.", func() {
			users := []User{
				{
					FirstName: "Fabien",
					LastName:  "Herfray",
					Password:  "qwertyuiop",
					Emails: []emails.Email{
						{Email: "fabien.herfray@me.com"},
					},
				},
			}

			users, err = repo.Create(users)

			So(err, ShouldBeNil)
			So(users[0].ID, ShouldEqual, 1)
		})

		Convey("Should be able to find users.", func() {
			users, err := repo.Find(nil)
			So(err, ShouldBeNil)
			So(users[0].FirstName, ShouldEqual, "Fabien")
		})

		Convey("Should be able to find user by id.", func() {
			user, err := repo.FindByID(1)
			So(err, ShouldBeNil)
			So(user.FirstName, ShouldEqual, "Fabien")

			user, err = repo.FindByID(10)
			So(err, ShouldNotBeNil)
		})

		Convey("Should be able to upsert users.", func() {
			users := []User{
				{
					FirstName: "Thomas",
					LastName:  "Hourlier",
					Password:  "1234",
					Emails: []emails.Email{
						{Email: "thomas.hourlier@cnode.fr"},
					},
				},
			}

			users, err = repo.Upsert(users)
			So(err, ShouldBeNil)
			So(users[0].ID, ShouldEqual, 2)

			users = []User{
				{
					ID:        2,
					FirstName: "Fabien",
					Emails: []emails.Email{
						{Email: "hourliert@gmail.com"},
					},
				},
			}

			users, err = repo.Upsert(users)
			So(err, ShouldBeNil)

			user, err := repo.FindByID(2)
			So(err, ShouldBeNil)
			So(user.FirstName, ShouldEqual, "Fabien")
			So(user.LastName, ShouldEqual, "Hourlier")
		})

		Convey("Should be able to filter results.", func() {
			filter := &interfaces.Filter{
				Fields: []string{"lastName"},
				Limit:  0,
				Offset: 0,
				Order:  "id asc",
			}

			users, err := repo.Find(filter)
			So(err, ShouldBeNil)
			So(users[0].FirstName, ShouldEqual, "")
			So(users[0].LastName, ShouldEqual, "Herfray")

			filter.Limit = 1
			users, err = repo.Find(filter)
			So(err, ShouldBeNil)
			So(len(users), ShouldEqual, 1)

			filter.Offset = 1
			users, err = repo.Find(filter)
			So(err, ShouldBeNil)
			So(users[0].LastName, ShouldEqual, "Hourlier")

			filter.Order = "id desc"
			users, err = repo.Find(filter)
			So(err, ShouldBeNil)
			So(users[0].LastName, ShouldEqual, "Herfray")

			filter.Limit = 0
			filter.Offset = 0

			users = []User{
				{
					ID:        2,
					FirstName: "Thomas",
				},
			}
			users, err = repo.Upsert(users)

			users, err = repo.Find(filter)
			So(len(users), ShouldEqual, 1)
		})

		Convey("Should be able to delete user by id.", func() {
			err := repo.DeleteByID(2)
			So(err, ShouldBeNil)

			_, err = repo.FindByID(2)
			So(err.Error(), ShouldEqual, "record not found")
		})

		Convey("Should be able to delete users.", func() {
			err := repo.DeleteAll(nil)
			So(err, ShouldBeNil)

			users := []User{}

			users, err = repo.Find(nil)
			So(err, ShouldBeNil)
			So(users, ShouldBeEmpty)
		})
	})

	err = store.Close()
	if err != nil {
		panic("Error closing database.")
	}

	err = os.Remove("test.db")
	if err != nil {
		panic("Error removing database.")
	}
}
