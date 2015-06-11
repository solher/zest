package infrastructure

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codegangsta/inject"
)

type Injector struct {
	injector inject.Injector
	deps     []interface{}
}

func NewInjector() *Injector {
	return &Injector{injector: inject.New()}
}

func (inj *Injector) Register(dependency interface{}) {
	inj.deps = append(inj.deps, dependency)
}

func (inj *Injector) RegisterMultiple(dependencies []interface{}) {
	inj.deps = append(inj.deps, dependencies...)
}

func (inj *Injector) GetByType(obj interface{}) interface{} {
	for _, dep := range inj.deps {
		if reflect.TypeOf(obj) == reflect.TypeOf(dep) {
			return dep
		}
	}

	return nil
}

func (inj *Injector) Get(obj interface{}) error {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("Invalid param: is not a pointer to a struct")
	}

	objValue := reflect.Indirect(reflect.ValueOf(obj))
	numField := objValue.NumField()

	for i := 0; i < numField; i++ {
		for _, dep := range inj.deps {
			if objValue.Field(i).Type() == reflect.TypeOf(dep) {
				objValue.Field(i).Set(reflect.ValueOf(dep))
			}
		}
	}

	obj = objValue.Interface()

	return nil
}

func (inj *Injector) Populate() error {
	injector := inj.injector
	failedDeps := dependencies{}
	values := []interface{}{}

	for _, obj := range inj.deps {
		failedDeps = append(failedDeps, dependency{Object: obj})
	}

	lastLen := [2]int{len(failedDeps) + 1, len(failedDeps) + 2}

	for len(failedDeps) > 0 {
		if lastLen[0] <= len(failedDeps) && lastLen[1] <= lastLen[0] {
			return errors.New(fmt.Sprintf("Dependencies not found: %v", failedDeps.GetMissing()))
		}
		lastLen[1] = lastLen[0]
		lastLen[0] = len(failedDeps)

		for _, dep := range failedDeps {
			obj := dep.Object
			kind := reflect.ValueOf(obj).Kind()

			switch kind {
			case reflect.Func:
				vals, err := injector.Invoke(obj)

				if err != nil {
				} else {
					failedDeps.Remove(dep)

					for _, val := range vals {
						injector.Map(val.Interface())
						values = append(values, val.Interface())
					}
				}
			case reflect.Struct, reflect.Ptr:
				failedDeps.Remove(dep)
				injector.Map(obj)
				values = append(values, obj)
			}
		}
	}

	inj.deps = values

	return nil
}

type dependency struct {
	Object interface{}
}

type dependencies []dependency

func (slc *dependencies) GetMissing() []reflect.Type {
	s := *slc
	missing := []reflect.Type{}

	for _, dep := range s {
		missing = append(missing, reflect.TypeOf(dep.Object))
	}

	return missing
}

func (slc *dependencies) Add(dep dependency) {
	s := *slc
	s = append(s, dep)
	*slc = s
}

func (slc *dependencies) Remove(dep dependency) {
	s := *slc

	for i, d := range s {
		if reflect.ValueOf(d.Object) == reflect.ValueOf(dep.Object) {
			s = append(s[:i], s[i+1:]...)
			*slc = s
			return
		}
	}
}
