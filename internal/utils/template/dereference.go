package template

type DereferenceFunc func(interface{}) interface{}

func Dereference(backoff string) DereferenceFunc {
	return func(i interface{}) interface{} {
		if v, ok := i.(*int32); ok {
			if v != nil {
				return *v
			}
		} else if v, ok := i.(*int64); ok {
			if v != nil {
				return *v
			}
		} else if v, ok := i.(*float64); ok {
			if v != nil {
				return *v
			}
		} else if v, ok := i.(*string); ok {
			if v != nil {
				return *v
			}
		}
		return backoff
	}
}
