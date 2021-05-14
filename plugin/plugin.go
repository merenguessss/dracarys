package plugin

type factory struct {
}

var PluginsFactory = New()

var New = func() *factory {
	return &factory{}
}

func (f *factory) Init() {

}
