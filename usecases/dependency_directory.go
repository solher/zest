package usecases

type dependencyDirectory struct {
	injector *Injector
}

func newDependencyDirectory() *dependencyDirectory {
	dependencyDir := &dependencyDirectory{injector: NewInjector()}

	return dependencyDir
}

func (dd *dependencyDirectory) Register(dependency interface{}) {
	dd.injector.Register(dependency)
}

func (dd *dependencyDirectory) RegisterMultiple(dependencies []interface{}) {
	dd.injector.RegisterMultiple(dependencies)
}

func (dd *dependencyDirectory) Populate() error {
	// for _, dep := range dd.injector.deps {
	// 	utils.Dump(reflect.ValueOf(dep))
	// }

	err := dd.injector.Populate()
	return err
}

func (dd *dependencyDirectory) Get(dependencies interface{}) error {
	err := dd.injector.Get(dependencies)
	return err
}

var DependencyDirectory *dependencyDirectory = newDependencyDirectory()
