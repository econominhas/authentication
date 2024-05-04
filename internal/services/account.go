package services

import "github.com/econominhas/authentication/internal/adapters"

type AccountService struct {
	GoogleAdapter   adapters.SignInProviderAdapter
	FacebookAdapter adapters.SignInProviderAdapter
}
