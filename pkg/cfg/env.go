package cfg

import (
	"os"

	"github.com/nitesh237/go-server-template/pkg/errors"
)

type Environment string

const (
	Test   Environment = "test"
	Dev    Environment = "dev"
	Docker Environment = "docker"
	QA     Environment = "qa"
	Prod   Environment = "prod"
)

func _getEnvironmentFromString(env string) (Environment, error) {
	switch env {
	case "test":
		return Test, nil
	case "dev":
		return Dev, nil
	case "qa":
		return QA, nil
	case "prod":
		return Prod, nil
	case "docker":
		return Docker, nil
	default:
		return "", errors.InvalidEnvironmentErrFn(env)
	}
}

// reads `ENVIRONMENT` env var
func GetEnvironment() (Environment, error) {
	envName, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return "", errors.ErrEnvironmentNotSet
	}

	return _getEnvironmentFromString(envName)
}

func IsLocalEnv() bool {
	env, err := GetEnvironment()
	if err != nil {
		return false
	}

	return env == Test || env == Dev || env == Docker
}
