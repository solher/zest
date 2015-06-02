package infrastructure

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type stubUser struct {
	ID        int
	FirstName string
	LastName  string
	Password  string
	Emails    []stubEmail
	CreatedAt time.Time
	UpdatedAt time.Time
}

type stubEmail struct {
	ID        int
	UserID    int `sql:"index"`
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TestGormStore(t *testing.T) {
	// store := NewGormStore()
	//
	Convey("Testing Gorm store...", t, func() {
		// 	Convey("Should be able to connect.", func() {
		// 		err := store.Connect("sqlite3", "test.db")
		// 		So(err, ShouldBeNil)
		// 	})
		//
		// 	Convey("Should be able to migrate tables.", func() {
		// 		err := store.MigrateTables([]interface{}{stubUser{}, stubEmail{}})
		// 		So(err, ShouldBeNil)
		// 	})
		//
		// 	Convey("Should be able to reinit tables.", func() {
		// 		err := store.ReinitTables([]interface{}{stubUser{}, stubEmail{}})
		// 		So(err, ShouldBeNil)
		// 	})
		//
		// 	Convey("Should be able to process filters.", func() {
		// 		filter := &interfaces.Filter{
		// 			Fields: []string{"firstName", "lastName"},
		// 			Limit:  2,
		// 			Offset: 1,
		// 			Order:  "firstName desc",
		// 			Where: map[string]interface{}{
		// 				"and": []interface{}{
		// 					map[string]interface{}{
		// 						"firstName":  map[string]interface{}{"like": "M%XY"},
		// 						"birthPlace": map[string]interface{}{"nlike": "%St%"},
		// 						"money":      200.5,
		// 					},
		// 				},
		// 				"or": []interface{}{
		// 					map[string]interface{}{"lastName": map[string]interface{}{"eq": "Fabien"}},
		// 					map[string]interface{}{"age": map[string]interface{}{"gt": 23}},
		// 					map[string]interface{}{"age": map[string]interface{}{"lt": 26}},
		// 				},
		// 				// "not": map[string]interface{}{"firstName": "Fabien"},
		// 				"not": map[string]interface{}{
		// 					"firstName": map[string]interface{}{"neq": "Fabien"},
		// 					"or": []interface{}{
		// 						map[string]interface{}{"lastName": "Herfray"},
		// 						map[string]interface{}{"money": map[string]interface{}{"gte": 0.0}},
		// 						map[string]interface{}{"money": map[string]interface{}{"lte": 1000.5}}},
		// 				},
		// 				"password":   "qwertyuiop",
		// 				"age":        22,
		// 				"graduated":  []int{2010, 2015},
		// 				"money":      3000.55,
		// 				"avg":        []float64{15.5, 13.24},
		// 				"birthPlace": []interface{}{"Chalon", "Macon"},
		// 			},
		// 			Include: []interface{}{
		// 				"Adresses",
		// 				map[string]interface{}{
		// 					"relation": "emAils",
		// 					"where": map[string]interface{}{
		// 						"graduated":  []int{2010, 2015},
		// 						"money":      3000.55,
		// 						"avg":        []float64{15.5, 13.24},
		// 						"birthPlace": []interface{}{"Chalon", "Macon"},
		// 					},
		// 					"include": []interface{}{"subjecT"},
		// 				},
		// 			},
		// 		}
		//
		// 		_, err := store.BuildQuery(filter)
		// 		So(err, ShouldBeNil)
		//
		// 		filter.Order = "filter desc asc"
		//
		// 		_, err = store.BuildQuery(filter)
		// 		So(err.Error(), ShouldEqual, "invalid order filter")
		// 	})
		//
		// 	Convey("Should be able to close connection.", func() {
		// 		err := store.Close()
		// 		So(err, ShouldBeNil)
		//
		// 		err = os.Remove("test.db")
		// 		So(err, ShouldBeNil)
		// 	})
	})
}
