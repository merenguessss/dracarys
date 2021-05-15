package plugin

type Factory struct {
}

var DefaultFactory = New()

var New = func() *Factory {
	return &Factory{}
}

func (f *Factory) Init() {

}
