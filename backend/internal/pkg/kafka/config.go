package kafka

type Config struct {
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password"`
	ProducerAck  string   `yaml:"producerAck"`
	CompressType string   `yaml:"compressType"`
	Addr         []string `yaml:"addr"`
}
