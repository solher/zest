package usecases

type dependencyDirectory struct {
	deps []interface{}
}

func newDependencyDirectory() *dependencyDirectory {
	return &dependencyDirectory{}
}

func (dd *dependencyDirectory) Register(dependency interface{}) {
	dd.deps = append(dd.deps, dependency)
}

func (dd *dependencyDirectory) RegisterMultiple(dependencies []interface{}) {
	dd.deps = append(dd.deps, dependencies...)
}

func (dd *dependencyDirectory) Get() []interface{} {
	return dd.deps
}

var DependencyDirectory *dependencyDirectory = newDependencyDirectory()
