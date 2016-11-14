package sensu

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/check"
	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/client"
)

const defaultRabbitMQURI string = "amqp://guest:guest@localhost:5672/%2f"

type configFlagSet struct {
	configFile string
	verbose    bool
}

type Config struct {
	flagSet *configFlagSet
	config  *configPayload
}

type configPayload struct {
	Client      *client.Client `json:"client,omitempty"`
	Checks      []*check.Check `json:"checks,omitempty"`
	RabbitMQURI *string        `json:"rabbitmq_uri,omitempty"`
}

func NewConfigFromFlagSet(flagset *configFlagSet) (*Config, error) {
	var cfg = Config{flagset, &configPayload{}}

	if flagset != nil && flagset.configFile != "" {
		buf, err := ioutil.ReadFile(flagset.configFile)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(buf, &cfg.config); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c *Config) RabbitMQURI() string {
	if cfg := c.config; cfg != nil && cfg.RabbitMQURI != nil {
		return *cfg.RabbitMQURI
	} else if uri := os.Getenv("RABBITMQ_URI"); uri != "" {
		return uri
	}

	return defaultRabbitMQURI
}

func (c *Config) Client() *client.Client {
	if cfg := c.config; cfg != nil && cfg.Client != nil {
		return cfg.Client
	}

	return &client.Client{
		Name:          os.Getenv("SENSU_CLIENT_NAME"),
		Address:       os.Getenv("SENSU_ADDRESS"),
		Subscriptions: strings.Split(os.Getenv("SENSU_CLIENT_SUBSCRIPTIONS"), ","),
	}
}

func (c *Config) Checks() []*check.Check {
	if cfg := c.config; cfg != nil {
		return cfg.Checks
	}

	return []*check.Check{}
}
