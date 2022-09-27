package kafkautil

type Config struct {
	BootstrapServers []string          `yaml:"bootstrap_servers,flow"`
	Topic            map[string]string `yaml:"topic"`
	Consumer         *ConsumerConfig   `yaml:"consumer"`
}

type ConsumerConfig struct {
	GroupId          string `yaml:"group_id"`
	EnableAutoCommit bool   `yaml:"enable_auto_commit"`
}
