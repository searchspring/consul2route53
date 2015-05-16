package consul2route53

import(
	"os"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the config settings for consul253
type Config struct {
	Consulhost string
	Consulport string
	Zoneid string
	Domain string
	Ttl int64
}

func(c *Config) SetConsulHost(consulhost string) {
	c.Consulhost = consulhost
}
func (c *Config) ConsulHost() string {
	return c.Consulhost
}

func(c *Config) SetConsulPort(consulport string) {
	c.Consulport = consulport
}
func (c *Config) ConsulPort() string {
	return c.Consulport
}

func (c *Config) SetZoneId(zoneid string) {
	c.Zoneid = zoneid
} 
func (c *Config) ZoneId() string {
	return c.Zoneid
}

func (c *Config) SetZone(zone string) {
	c.Domain = zone
}
func (c *Config) Zone() string {
	return c.Domain
}

func (c *Config) SetTTL(ttl int64) {
	c.Ttl = ttl
}
func (c *Config) TTL() int64 {
	return c.Ttl
}

// ReadConfig reads a config yaml file on disk and returns config or error.
func ReadConfig(filename string) (Config, error) {
	config := Config{}
	_,err := os.Stat(filename)
	if err != nil {
		return config, err
	}
	cfg,err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(cfg, &config)
	return config, err
}