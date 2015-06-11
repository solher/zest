package infrastructure

//
// import (
// 	"reflect"
// 	"strings"
//
// 	"github.com/solher/zest/utils"
// 	"github.com/codegangsta/inject"
// )
//
// type missingDependency struct {
// 	Object          interface{}
// 	Missing, Future reflect.Type
// 	FieldsToInject  []string
// }
//
// type missingDependencies []missingDependency
//
// func (slc *missingDependencies) GetMissing() []reflect.Type {
// 	s := *slc
// 	missing := []reflect.Type{}
//
// 	for _, dep := range s {
// 		skip := false
//
// 		for _, dep2 := range s {
// 			futureImplementsMissing := false
// 			if dep.Missing.Kind() == reflect.Interface {
// 				futureImplementsMissing = dep2.Future.Implements(dep.Missing)
// 			}
//
// 			missingImplementsFuture := false
// 			if dep2.Future.Kind() == reflect.Interface {
// 				missingImplementsFuture = dep.Missing.Implements(dep2.Future)
// 			}
//
// 			if dep.Missing == dep2.Future || futureImplementsMissing || missingImplementsFuture {
// 				skip = true
// 			}
// 		}
//
// 		if skip {
// 			continue
// 		}
//
// 		missing = append(missing, dep.Missing)
// 	}
//
// 	return missing
// }
//
// func (slc *missingDependencies) Add(missingDep *missingDependency) {
// 	s := *slc
// 	s = append(s, *missingDep)
// 	*slc = s
// }
//
// func (slc *missingDependencies) Find(missing, future reflect.Type) *missingDependency {
// 	s := *slc
//
// 	for _, dep := range s {
// 		// future1ImplementsFuture2 := false
// 		// if future.Kind() == reflect.Interface {
// 		// 	future1ImplementsFuture2 = dep.Future.Implements(future)
// 		// }
// 		//
// 		// future2ImplementsFuture1 := false
// 		// if dep.Future.Kind() == reflect.Interface {
// 		// 	future2ImplementsFuture1 = future.Implements(dep.Future)
// 		// }
// 		//
// 		// missing1ImplementsMissing2 := false
// 		// if missing.Kind() == reflect.Interface {
// 		// 	missing1ImplementsMissing2 = dep.Missing.Implements(missing)
// 		// }
// 		//
// 		// missing2ImplementsMissing1 := false
// 		// if dep.Missing.Kind() == reflect.Interface {
// 		// 	missing2ImplementsMissing1 = missing.Implements(dep.Missing)
// 		// }
//
// 		// if (dep.Future == future || future1ImplementsFuture2 || future2ImplementsFuture1) &&
// 		// 	(dep.Missing == missing || missing1ImplementsMissing2 || missing2ImplementsMissing1) {
// 		// 	return &dep
// 		// }
// 		if dep.Future == future && dep.Missing == missing {
// 			return &dep
// 		}
// 	}
//
// 	return nil
// }
//
// func (slc *missingDependencies) FindByObject(object interface{}) *missingDependency {
// 	s := *slc
//
// 	utils.Dump("VALUE1")
// 	utils.Dump(reflect.ValueOf(object))
//
// 	utils.Dump(len(s))
//
// 	for _, dep := range s {
// 		utils.Dump("VALUE2")
// 		utils.Dump(reflect.ValueOf(dep.Object))
//
// 		if reflect.ValueOf(dep.Object) == reflect.ValueOf(object) {
// 			utils.Dump("EQUALS")
// 			return &dep
// 		}
// 	}
//
// 	return nil
// }
//
// func (slc *missingDependencies) Remove(obj interface{}) {
// 	s := *slc
//
// 	for i, dep := range s {
// 		if reflect.DeepEqual(dep.Object, obj) {
// 			s = append(s[:i], s[i+1:]...)
// 			*slc = s
// 			return
// 		}
// 	}
// }
//
// type Injector struct {
// 	injector inject.Injector
// 	deps     []interface{}
// }
//
// func NewInjector() *Injector {
// 	return &Injector{injector: inject.New()}
// }
//
// func (inj *Injector) Register(dependency interface{}) {
// 	inj.deps = append(inj.deps, dependency)
// }
//
// func (inj *Injector) RegisterMultiple(dependencies []interface{}) {
// 	inj.deps = append(inj.deps, dependencies...)
// }
//
// func (inj *Injector) Populate() ([]interface{}, error) {
// 	// Missing are missing dependencies, Future are future dependencies that will be available if the missing ones are found.
// 	missingDeps := missingDependencies{}
// 	// Keys are partially injected dependencies indexes, values are dependency fields to manually inject.
// 	toManuallyInject := make(map[int]missingDependency)
//
// 	dependencies := inj.deps
// 	// populated := []interface{}{}
//
// 	injector := inj.injector
// 	values := []interface{}{}
// 	firstTry := true
// 	lastLen := 1
//
// 	for firstTry || len(missingDeps) > 0 {
// 		if lastLen == len(missingDeps) {
// 			// return nil, errors.New(fmt.Sprintf("Dependencies not found: %v", missingDeps.GetMissing()))
// 		}
// 		lastLen = len(missingDeps)
//
// 		for i := range dependencies {
// 			if !firstTry && missingDeps.FindByObject(dependencies[i]) == nil {
// 				continue
// 			}
//
// 			if reflect.ValueOf(dependencies[i]).Kind() == reflect.Func {
// 				utils.Dump("Trying to invoke constructor")
// 				vals, err := injector.Invoke(dependencies[i])
// 				if err == nil {
// 					missingDeps.Remove(dependencies[i])
//
// 					for _, val := range vals {
// 						interVal := val.Interface()
//
// 						injector.Map(interVal)
// 						values = append(values, interVal)
//
// 						utils.Dump("Success:")
// 						utils.Dump(val)
// 					}
// 				} else if strings.Contains(err.Error(), "Value not found for type") {
// 					missingType, fieldsToInject := inj.getMissingType(dependencies[i])
//
// 					utils.Dump("Error: missing type")
// 					utils.Dump(missingType)
//
// 					future := reflect.TypeOf(dependencies[i])
// 					if reflect.ValueOf(dependencies[i]).Kind() == reflect.Func {
// 						future = future.Out(0)
// 					}
//
// 					missingDep := &missingDependency{Object: dependencies[i], Missing: missingType, Future: future, FieldsToInject: fieldsToInject}
//
// 					if missingDeps.Find(missingDep.Future, missingDep.Missing) != nil {
// 						utils.Dump("Dep to manually inject found")
// 						utils.Dump(dependencies[i])
// 						toManuallyInject[i] = *missingDep
//
// 						var injectorReflect reflect.Value
// 						reflect.Copy(injectorReflect, reflect.ValueOf(injector))
// 						injectorTmp := injectorReflect.Interface().(inject.Injector)
//
// 						injectorTmp.Map(reflect.New(missingType))
//
// 						vals, err := injectorTmp.Invoke(dependencies[i])
// 						if err == nil {
// 							missingDeps.Remove(dependencies[i])
//
// 							for _, val := range vals {
// 								interVal := val.Interface()
//
// 								injector.Map(interVal)
// 								values = append(values, interVal)
// 							}
// 						}
// 					} else {
// 						if missingDeps.FindByObject(dependencies[i]) == nil {
// 							utils.Dump("Adding missing dep")
// 							utils.Dump(missingDep)
// 							missingDeps.Add(missingDep)
// 						}
// 					}
// 				} else {
// 					return nil, err
// 				}
// 			} else {
// 				injector.Map(dependencies[i])
// 				values = append(values, dependencies[i])
// 			}
// 		}
//
// 		firstTry = false
// 	}
//
// 	for index, missingDependency := range toManuallyInject {
// 		fieldNames := missingDependency.FieldsToInject
// 		var injectValue reflect.Value
//
// 		for _, value := range values {
// 			if reflect.TypeOf(value) == missingDependency.Missing {
// 				injectValue = reflect.ValueOf(value)
// 			}
// 		}
//
// 		for _, name := range fieldNames {
// 			reflect.ValueOf(dependencies[index]).FieldByName(name).Set(injectValue)
// 		}
// 	}
//
// 	return values, nil
// }
//
// func (inj *Injector) getMissingType(f interface{}) (reflect.Type, []string) {
// 	t := reflect.TypeOf(f)
//
// 	var missingType reflect.Type
//
// 	for i := 0; i < t.NumIn(); i++ {
// 		argType := t.In(i)
// 		val := inj.injector.Get(argType)
// 		if !val.IsValid() {
// 			missingType = argType
// 			break
// 		}
// 	}
//
// 	fieldsToInject := []string{}
//
// 	if missingType != nil {
// 		if reflect.ValueOf(t).Kind() == reflect.Func {
// 			t = missingType.Out(0)
// 		}
//
// 		if reflect.ValueOf(t).Kind() == reflect.Struct {
// 			for i := 0; i < t.NumField(); i++ {
// 				field := t.Field(i)
// 				if field.Type == missingType {
// 					fieldsToInject = append(fieldsToInject, field.Name)
// 				}
// 			}
// 		}
// 	}
//
// 	return missingType, fieldsToInject
// }
