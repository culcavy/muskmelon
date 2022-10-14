package object

// NewEnvironment Environment 的构造函数
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Environment 存放上下文。本质是个 hashmap
type Environment struct {
	store map[string]Object
}

// Get 从 hashmap 中取数据
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set 存数据到 hashmap 中
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
