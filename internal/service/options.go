package service

type option string

const whereOption option = "where"

type Option interface {
	name() string
	value() interface{}
}

type iOption struct {
	Name  string
	Value interface{}
}

func (impl *iOption) name() string {
	return impl.Name
}

func (impl *iOption) value() interface{} {
	return impl.Value
}

// OptionWhere
func OptionWhere(query string, arg interface{}) Option {
	return &iOption{
		Name:  string(whereOption),
		Value: []interface{}{query, arg},
	}
}

func GetOptionWhere(opts []Option) (string, interface{}) {
	for _, opt := range opts {
		if opt != nil && opt.name() == string(whereOption) {
			v := opt.value().([]interface{})
			return v[0].(string), v[1]
		}
	}

	return "", ""
}
