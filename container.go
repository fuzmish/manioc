package manioc

type defaultContainer struct {
	defaultScope
}

func (c *defaultContainer) getRegisterContext() registerContext {
	return c.context
}

func newDefaultContainer() *defaultContainer {
	return &defaultContainer{
		defaultScope: defaultScope{
			context:     newDefaultContext(),
			childScopes: make([]Scope, 0),
		},
	}
}

//nolint:gochecknoglobals
var globalContainer Container = newDefaultContainer()
