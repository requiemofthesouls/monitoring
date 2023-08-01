package def

import (
	"github.com/requiemofthesouls/container"
	"github.com/requiemofthesouls/monitoring"
)

const DIWrapper = "monitoring.wrapper"

type Wrapper = monitoring.Wrapper

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIWrapper,
			Build: func(container container.Container) (_ interface{}, err error) {
				return monitoring.New(), nil
			},
		})
	})
}
