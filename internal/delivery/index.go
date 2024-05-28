package delivery

type Delivery interface {
	Listen()
}

type Validator interface {
	Validate(i interface{}) error
}
