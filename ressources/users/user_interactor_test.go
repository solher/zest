// Generated by: main
// TypeWriter: interactor_test
// Directive: +gen on User

package users

import (
	"testing"

	"github.com/Solher/auth-scaffold/interfaces"
	. "github.com/smartystreets/goconvey/convey"
)

type stubRepository struct{}

func (r *stubRepository) Create(users []User) ([]User, error) {
	return users, nil
}

func (r *stubRepository) Find(filter *interfaces.Filter) ([]User, error) {
	return []User{}, nil
}

func (r *stubRepository) FindByID(id int, filter *interfaces.Filter) (*User, error) {
	return &User{}, nil
}

func (r *stubRepository) Upsert(users []User) ([]User, error) {
	return users, nil
}

func (r *stubRepository) DeleteAll(filter *interfaces.Filter) error {
	return nil
}

func (r *stubRepository) DeleteByID(id int) error {
	return nil
}

func TestInteractor(t *testing.T) {
	repo := &stubRepository{}
	interactor := NewInteractor(repo)

	Convey("Testing users interactor...", t, func() {
		Convey("Should be able to create users.", func() {
			users := []User{}
			_, err := interactor.Create(users)

			So(err, ShouldBeNil)
		})

		Convey("Should be able to find users.", func() {
			_, err := interactor.Find(nil)

			So(err, ShouldBeNil)
		})

		Convey("Should be able to find a user by id.", func() {
			_, err := interactor.FindByID(1, nil)

			So(err, ShouldBeNil)
		})

		Convey("Should be able to upsert users.", func() {
			users := []User{}
			_, err := interactor.Upsert(users)

			So(err, ShouldBeNil)
		})

		Convey("Should be able to delete users.", func() {
			err := interactor.DeleteAll(nil)
			So(err, ShouldBeNil)
		})

		Convey("Should be able to delete a user by id.", func() {
			err := interactor.DeleteByID(1)
			So(err, ShouldBeNil)
		})
	})
}
