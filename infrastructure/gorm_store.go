package infrastructure

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/interfaces"
	"github.com/solher/zest/usecases"
	"github.com/solher/zest/utils"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore() *GormStore {
	return &GormStore{}
}

func (st *GormStore) Connect(adapter, url string) error {
	db, err := gorm.Open(adapter, url)
	db.LogMode(true)
	st.db = &db

	return err
}

func (st *GormStore) Close() error {
	err := st.db.Close()

	return err
}

func (st *GormStore) GetDB() *gorm.DB {
	return st.db
}

func (st *GormStore) MigrateTables(tables []interface{}) error {
	for _, table := range tables {
		upcastedTable := reflect.New(reflect.TypeOf(table)).Interface()

		err := st.db.AutoMigrate(upcastedTable).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *GormStore) ResetTables(tables []interface{}) error {
	for _, table := range tables {
		upcastedTable := reflect.New(reflect.TypeOf(table)).Interface()

		err := st.db.DropTableIfExists(upcastedTable).Error
		if err != nil {
			return err
		}

		err = st.db.CreateTable(upcastedTable).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *GormStore) BuildQuery(filter *usecases.Filter, ownerRelations []domain.DBRelation) (*gorm.DB, error) {
	query := st.db

	if ownerRelations != nil {
		for _, relation := range ownerRelations {
			relation.Resource = utils.ToDBName(relation.Resource)
			relation.Fk = utils.ToDBName(relation.Fk)
			relation.Related = utils.ToDBName(relation.Related)

			queryString := fmt.Sprintf("INNER JOIN %s ON %s.%s = %s.id", relation.Resource, relation.Resource, relation.Fk, relation.Related)
			query = query.Joins(queryString)
			query = query.Table(relation.Related)
		}
	}

	if filter != nil {
		gormFilter, err := processFilter(filter)
		if err != nil {
			return nil, err
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

		for _, include := range gormFilter.Include {
			if include.Relation == "" {
				break
			}

			if include.Where == "" {
				query = query.Preload(include.Relation)
			} else {
				query = query.Preload(include.Relation, include.Where)
			}
		}
	}

	return query, nil
}

const (
	orSql    = " OR "
	andSql   = " AND "
	gtSql    = " > "
	gteSql   = " >= "
	ltSql    = " < "
	lteSql   = " <= "
	eqSql    = " = "
	neqSql   = " <> "
	likeSql  = " LIKE "
	nlikeSql = " NOT LIKE "
)

func processFilter(filter *usecases.Filter) (*interfaces.GormFilter, error) {
	fields := filter.Fields
	dbNamedFields := make([]string, len(fields))

	for i, field := range fields {
		dbNamedFields[i] = utils.ToDBName(field)
	}

	if filter.Order != "" {
		order := strings.ToLower(filter.Order)
		matched, err := regexp.MatchString("\\A\\w+ (asc|desc)\\z", order)
		if err != nil || !matched {
			return nil, errors.New("invalid order filter")
		}

		split := strings.Split(filter.Order, " ")
		filter.Order = utils.ToDBName(split[0]) + " " + split[1]
	}

	processedFilter := &interfaces.GormFilter{Fields: dbNamedFields, Limit: filter.Limit, Offset: filter.Offset, Order: filter.Order}

	buffer := &bytes.Buffer{}
	err := processCondition(buffer, "", andSql, "", filter.Where)
	if err != nil {
		return nil, err
	}
	processedFilter.Where = buffer.String()

	gormIncludes, err := processInclude(filter.Include)
	if err != nil {
		return nil, err
	}
	processedFilter.Include = gormIncludes

	return processedFilter, nil
}

func processCondition(buffer *bytes.Buffer, attribute, operator, sign string, condition interface{}) error {
	switch condition.(type) {
	case map[string]interface{}:
		processUnaryCondition(buffer, attribute, operator, condition.(map[string]interface{}))

	case interface{}:
		if buffer.Len() != 0 {
			buffer.WriteString(operator)
		}
		processOperation(buffer, attribute, operator, sign, condition)
	}

	return nil
}

func processUnaryCondition(buffer *bytes.Buffer, attribute, operator string, condition map[string]interface{}) error {
	for key := range condition {
		lowerKey := strings.ToLower(key)

		switch lowerKey {
		case "gt":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", gtSql, condition[key])
			break

		case "gte":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", gteSql, condition[key])
			break

		case "lt":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", ltSql, condition[key])
			break

		case "lte":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", lteSql, condition[key])
			break

		case "eq":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", eqSql, condition[key])
			break

		case "neq":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", neqSql, condition[key])
			break

		case "like":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", likeSql, condition[key])
			break

		case "nlike":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, attribute, "", nlikeSql, condition[key])
			break

		case "not":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			newBuffer := &bytes.Buffer{}

			buffer.WriteString("NOT (")
			processCondition(newBuffer, "", andSql, eqSql, condition[key])

			buffer.Write(newBuffer.Bytes())
			buffer.WriteString(")")

		case "or":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, "", orSql, eqSql, condition[key].([]interface{}))

		case "and":
			if buffer.Len() != 0 {
				buffer.WriteString(operator)
			}
			processOperation(buffer, "", andSql, eqSql, condition[key].([]interface{}))

		default:
			processCondition(buffer, key, operator, eqSql, condition[key])
		}
	}

	return nil
}

func processOperation(buffer *bytes.Buffer, attribute, operator, sign string, condition interface{}) error {
	switch condition.(type) {
	case bool:
		if condition.(bool) {
			processSimpleOperationStr(buffer, attribute, sign, "1")
		} else {
			processSimpleOperationStr(buffer, attribute, sign, "0")
		}

	case string:
		processSimpleOperationStr(buffer, attribute, sign, condition.(string))

	case int:
		processSimpleOperation(buffer, attribute, sign, strconv.FormatInt(int64(condition.(int)), 10))

	case float64:
		processSimpleOperation(buffer, attribute, sign, strconv.FormatFloat(condition.(float64), 'f', -1, 64))

	case []int:
		intArray := condition.([]int)
		lenArray := len(intArray)

		buffer.WriteString(utils.ToDBName(attribute))
		buffer.WriteString(" IN (")

		for i, value := range intArray {
			buffer.WriteString(strconv.FormatInt(int64(value), 10))
			if i < lenArray-1 {
				buffer.WriteString(", ")
			}
		}

		buffer.WriteString(")")

	case []float64:
		floatArray := condition.([]float64)
		lenArray := len(floatArray)

		buffer.WriteString(utils.ToDBName(attribute))
		buffer.WriteString(" IN (")

		for i, value := range floatArray {
			buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
			if i < lenArray-1 {
				buffer.WriteString(", ")
			}
		}

		buffer.WriteString(")")

	case []interface{}:
		conditions := condition.([]interface{})

		arrStr := []string{}
		strType := reflect.TypeOf("")

		for _, condition := range conditions {
			if reflect.TypeOf(condition) == strType {
				arrStr = append(arrStr, condition.(string))
			}
		}

		if len(arrStr) == 0 {
			newBuffer := &bytes.Buffer{}

			buffer.WriteString("(")

			for _, condition := range conditions {
				processCondition(newBuffer, "", operator, sign, condition)
			}
			buffer.Write(newBuffer.Bytes())

			buffer.WriteString(")")
		} else {
			lenArray := len(arrStr)

			buffer.WriteString(utils.ToDBName(attribute))
			buffer.WriteString(" IN (")

			for i, value := range arrStr {
				buffer.WriteRune('\'')
				buffer.WriteString(value)
				buffer.WriteRune('\'')

				if i < lenArray-1 {
					buffer.WriteString(", ")
				}
			}

			buffer.WriteString(")")
		}
	}

	return nil
}

func processSimpleOperation(buffer *bytes.Buffer, attribute, sign, condition string) {
	buffer.WriteString(utils.ToDBName(attribute))
	buffer.WriteString(sign)
	buffer.WriteString(condition)
}

func processSimpleOperationStr(buffer *bytes.Buffer, attribute, sign, condition string) {
	buffer.WriteString(utils.ToDBName(attribute))
	buffer.WriteString(sign)
	buffer.WriteRune('\'')
	buffer.WriteString(condition)
	buffer.WriteRune('\'')
}

func processInclude(include []interface{}) ([]interfaces.GormInclude, error) {
	processedIncludes := []interfaces.GormInclude{}

	processedIncludes, err := processNestedInclude(include, processedIncludes, "")
	if err != nil {
		return nil, err
	}

	return processedIncludes, nil
}

func processNestedInclude(include interface{}, processedIncludes []interfaces.GormInclude, parentModel string) ([]interfaces.GormInclude, error) {
	switch include.(type) {
	case []interface{}:
		includeArr := include.([]interface{})

		for _, nestedInclude := range includeArr {
			var err error
			processedIncludes, err = processNestedInclude(nestedInclude, processedIncludes, parentModel)
			if err != nil {
				return nil, err
			}
		}

	case map[string]interface{}:
		includeMap := include.(map[string]interface{})
		processedInclude := interfaces.GormInclude{}

		value := includeMap["relation"]
		switch strValue := value.(type) {
		case string:
			processedInclude.Relation = parentModel + strings.Title(strValue)
		}

		value = includeMap["where"]
		buffer := &bytes.Buffer{}
		err := processCondition(buffer, "", andSql, "", value)
		if err != nil {
			return nil, err
		}
		processedInclude.Where = buffer.String()

		value = includeMap["include"]
		processedIncludes, err = processNestedInclude(value, processedIncludes, processedInclude.Relation+".")
		if err != nil {
			return nil, err
		}

		processedIncludes = append(processedIncludes, processedInclude)

	case string:
		relation := parentModel + strings.Title(include.(string))
		processedInclude := interfaces.GormInclude{Relation: relation}
		processedIncludes = append(processedIncludes, processedInclude)
	}

	return processedIncludes, nil
}
