package client

type Client interface {
	Invoke()
}

var DefaultClient = New()

func New() *defaultClient {
	return &defaultClient{}
}

type defaultClient struct {
	beforeHandle []interface{}
	afterHandle  []interface{}
}

func (defaultClient *defaultClient) Invoke() {

}
