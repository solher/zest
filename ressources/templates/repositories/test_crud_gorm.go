package repositories

import (
	"github.com/Solher/zest/ressources/templates"
	"github.com/clipperhouse/typewriter"
)

func init() {
	imports := []typewriter.ImportSpec{
		typewriter.ImportSpec{
			Name: ".",
			Path: "github.com/smartystreets/goconvey/convey",
		},
		typewriter.ImportSpec{
			Name: "_",
			Path: "github.com/mattn/go-sqlite3",
		},
	}

	err := typewriter.Register(templates.NewWrite("repository_test", testSlice, imports))
	if err != nil {
		panic(err)
	}
}

var testSlice = typewriter.TemplateSlice{
	repositoryTest,
}

var repositoryTest = &typewriter.Template{
	Name: "Repository_test",
	Text: `
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
			gormFilter.Where = ""
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

		err = db.AutoMigrate(&{{.Type}}{}).Error
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

		Convey("Testing {{.Name}}s repository...", t, func() {
			Convey("Should be able to create {{.Name}}s.", func() {
				{{.Name}}s := []{{.Type}}{
					{},
				}

				{{.Name}}s, err = repo.Create({{.Name}}s)

				So(err, ShouldBeNil)
				So({{.Name}}s[0].ID, ShouldEqual, 1)
			})

			Convey("Should be able to find {{.Name}}s.", func() {
				{{.Name}}s, err := repo.Find(nil)
				So(err, ShouldBeNil)
				So(users[0].ID, ShouldEqual, 1)
			})

			Convey("Should be able to find {{.Name}} by id.", func() {
				{{.Name}}, err := repo.FindByID(1, nil)
				So(err, ShouldBeNil)
				So(user.ID, ShouldEqual, 1)

				{{.Name}}, err = repo.FindByID(10, nil)
				So(err, ShouldNotBeNil)
			})

			Convey("Should be able to upsert {{.Name}}s.", func() {
				{{.Name}}s := []{{.Type}}{
					{},
				}

				{{.Name}}s, err = repo.Upsert({{.Name}}s)
				So(err, ShouldBeNil)
				So({{.Name}}s[0].ID, ShouldEqual, 2)

				{{.Name}}s = []{{.Type}}{
					{},
				}

				{{.Name}}s, err = repo.Upsert({{.Name}}s)
				So(err, ShouldBeNil)

				{{.Name}}, err := repo.FindByID(2, nil)
				So(err, ShouldBeNil)
				So(user.ID, ShouldEqual, 2)
			})

			Convey("Should be able to filter results.", func() {
				filter := &interfaces.Filter{
					Fields: []string{},
					Limit:  0,
					Offset: 0,
					Order:  "id asc",
				}

				{{.Name}}s, err := repo.Find(filter)
				So(err, ShouldBeNil)

				filter.Limit = 1
				{{.Name}}s, err = repo.Find(filter)
				So(err, ShouldBeNil)
				So(len({{.Name}}s), ShouldEqual, 1)

				filter.Offset = 1
				{{.Name}}s, err = repo.Find(filter)
				So(err, ShouldBeNil)

				filter.Order = "id desc"
				{{.Name}}s, err = repo.Find(filter)
				So(err, ShouldBeNil)

				filter.Limit = 0
				filter.Offset = 0

				{{.Name}}s = []{{.Type}}{

				}
				{{.Name}}s, err = repo.Upsert({{.Name}}s)

				{{.Name}}s, err = repo.Find(filter)
				So(len({{.Name}}s), ShouldEqual, 3)
			})

			Convey("Should be able to delete {{.Name}} by id.", func() {
				err := repo.DeleteByID(2)
				So(err, ShouldBeNil)

				_, err = repo.FindByID(2, nil)
				So(err.Error(), ShouldEqual, "record not found")
			})

			Convey("Should be able to delete {{.Name}}s.", func() {
				err := repo.DeleteAll(nil)
				So(err, ShouldBeNil)

				{{.Name}}s := []{{.Type}}{}

				{{.Name}}s, err = repo.Find(nil)
				So(err, ShouldBeNil)
				So({{.Name}}s, ShouldBeEmpty)
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
`}
