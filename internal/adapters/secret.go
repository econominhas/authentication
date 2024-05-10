package adapters

type SecretAdapter interface {
	GenSecret(length int) (string, error)
}
