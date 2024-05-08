package adapters

type IdAdapter interface {
	GenId() (string, error)
}
