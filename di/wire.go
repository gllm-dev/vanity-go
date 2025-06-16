//go:build wireinject
// +build wireinject

package di

import (
	"errors"
	"github.com/google/wire"
	"os"

	"go.gllm.dev/vanity-go/internal/adapters/handlers/rest"
	"go.gllm.dev/vanity-go/internal/services/gosvc"
)

type Domain string
type Repository string

func ProvideDomain() (Domain, error) {
	domain := os.Getenv("VANITY_DOMAIN")
	if domain == "" {
		return "", errors.New("VANITY_DOMAIN environment variable not set")
	}
	return Domain(domain), nil
}

func ProvideRepository() (Repository, error) {
	repository := os.Getenv("VANITY_REPOSITORY")
	if repository == "" {
		return "", errors.New("VANITY_REPOSITORY environment variable not set")
	}
	return Repository(repository), nil
}

func ProvideService(domain Domain, repository Repository) *gosvc.Service {
	return gosvc.New(string(domain), string(repository))
}

var serviceSet = wire.NewSet(
	ProvideDomain,
	ProvideRepository,
	ProvideService,
)

func ProvideRestServer() (*rest.Server, error) {
	wire.Build(
		rest.New,
		rest.LoadConfig,
		serviceSet,
	)

	return nil, nil
}
