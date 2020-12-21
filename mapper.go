package eval

type MapperFunc func() map[string]interface{}
type TriggerFunc func(input map[string]interface{}) error
