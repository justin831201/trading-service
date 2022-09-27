package env

import (
	"github.com/justin831201/trading-service/pkg/kafkautil"
	"github.com/justin831201/trading-service/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Address string            `yaml:"address"`
	Logger  *logger.Config    `yaml:"logger"`
	Kafka   *kafkautil.Config `yaml:"kafka"`
}

var config *Config

// LoadConfig looks for config file from the command line and parse the yml file into config struct
func LoadConfig(configFile string) *Config {
	content, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// Load yml file into defined structure
	config = new(Config)
	_ = yaml.Unmarshal(content, config)
	return config
}

func GetConfig() *Config {
	return config
}
