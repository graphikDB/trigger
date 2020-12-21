package eval

// Mapper returns it's values as a map[string]interface{}
type Mapper interface {
	AsMap() map[string]interface{}
}

// MapperFunc implements Mapper
type MapperFunc func() map[string]interface{}

func (m MapperFunc) AsMap() map[string]interface{} {
	return m()
}
