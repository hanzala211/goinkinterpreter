package evaluator

type Environment struct {
	values map[string]any
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
		parent: nil,
	}
}

func NewEnvironmentWithParent(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]any),
		parent: parent,
	}
}

func (e *Environment) Set(name string, value any) {
	if _, ok := e.values[name]; ok {
		panic("variable already exists")
	}
	e.values[name] = value
}

func (e *Environment) Assign(name string, newValue any) {
	if _, ok := e.values[name]; ok {
		e.values[name] = newValue
	}
	if e.parent != nil {
		e.parent.Assign(name, newValue)
	}
	panic("variable does not exist")
}

func (e *Environment) Get(name string) any {
	if v, ok := e.values[name]; ok {
		return v
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil
}
