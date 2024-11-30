package cfg

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nitesh237/go-server-template/pkg/errors"

	"github.com/spf13/viper"
)

type ConfigType string

const (
	ConfigTypeYaml ConfigType = "yml"
)

const (
	defaultConfigDirectory     = "config"
	configParamsFileNameFormat = "%s-params"
	configEnvFileNameFormat    = "%s-%s"
)

func Load[T any](configDir, configName string, configType ConfigType) (*T, error) {
	env, err := GetEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get environment")
	}

	// loads config from file
	fmt.Printf("Initializing Environment [%s] and ConfigDir [%s]\n", env, configDir)

	// new viper instance per service to resolve conflict
	vp := viper.New()

	replacer := strings.NewReplacer(".", "_")
	vp.SetEnvKeyReplacer(replacer)

	// Set the path to look for the configurations file
	vp.AddConfigPath(configDir)

	// Set cfg file type
	vp.SetConfigType(string(configType))

	// Get params file name as per the convention
	fileName := GetParamFileName(configName)

	vp.SetConfigName(fileName)

	fmt.Printf("Reading from config file: %s\n", fileName)
	if err = vp.ReadInConfig(); err != nil {
		// ignoring the error here as params files are optional
		fmt.Printf("Unable to read params file %v, will proceed reading env specific config\n", err)
	}

	// Get env file name as per convention
	fileName = GetEnvFileName(env, configName)

	// Set the file name of the configurations file
	vp.SetConfigName(fileName)

	// Enable viper to read Environment Variables
	vp.AutomaticEnv()

	fmt.Printf("Reading from config file: %s\n", fileName)

	if err = vp.MergeInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to merge configs file")
	}

	config := new(T)

	if err = vp.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config from viper")
	}

	return config, nil
}

// return config params file name e.g., "frontend-params" as per convention "service-params.yml"
func GetParamFileName(configName string) string {
	return fmt.Sprintf(configParamsFileNameFormat, configName)
}

// return config env file name e.g., "frontend-qa" as per convention "service-env.yml"
func GetEnvFileName(env Environment, configName string) string {
	return fmt.Sprintf(configEnvFileNameFormat, configName, env)
}

// reads `CONFIG_DIR` env var
// If the env var is not set, we default to 'cwd/config'
func GetConfigDir() (string, error) {
	configDir, ok := os.LookupEnv("CONFIG_DIR")
	if !ok {
		currDir, err := os.Getwd()
		if err != nil {
			return "", errors.New("CONFIG_DIR not found")
		}
		configDir = filepath.Join(currDir, defaultConfigDirectory)
	}

	return configDir, nil
}

type Endpoint struct {
	Host     string
	Port     int
	IsSecure bool
}

func (e *Endpoint) GetURL() *url.URL {
	if e.IsSecure {
		return &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:%d", e.Host, e.Port),
		}
	}

	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", e.Host, e.Port),
	}
}
