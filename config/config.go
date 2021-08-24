package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper

	APIURL   string `mapstructure:"api_url"`
	APIToken string `mapstructure:"api_token"`
}

var Defaults = &Config{
	APIURL:   "",
	APIToken: "",
}

func New() (*Config, error) {
	c := &Config{
		viper: newViper(),
	}

	c.SetDefault("api_url", Defaults.APIURL)
	if err := c.BindEnv("url"); err != nil {
		return nil, err
	}

	c.SetDefault("api_token", Defaults.APIToken)
	if err := c.BindEnv("api_token"); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Load() error {
	err := c.viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return c.viper.Unmarshal(c)
		}
		return err
	}

	return c.viper.Unmarshal(c)
}

func (c *Config) AllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}

func (c *Config) ConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}

func (c *Config) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}

func (c *Config) SetConfigFile(file string) {
	c.viper.SetConfigFile(file)
}

func (c *Config) BindEnv(name string) error {
	return c.viper.BindEnv(name)
}

func (c *Config) BindFlagValue(key string, flag viper.FlagValue) error {
	return c.viper.BindFlagValue(key, flag)
}

func (c *Config) BindFlagValues(flags viper.FlagValueSet) error {
	return c.viper.BindFlagValues(flags)
}

func (c *Config) BindPFlag(key string, flag *pflag.Flag) error {
	return c.viper.BindPFlag(key, flag)
}

func (c *Config) BindPFlags(flags *pflag.FlagSet) error {
	return c.viper.BindPFlags(flags)
}

func (c *Config) WriteConfig() error {
	//nolint:errorlint
	switch err := c.viper.WriteConfig().(type) {
	case nil:
		return nil
	case viper.ConfigFileNotFoundError:
		homedir, e := os.UserHomeDir()
		if e != nil {
			return e
		}
		fp := filepath.Join(homedir, ".katapult")
		if err := os.MkdirAll(fp, 0o777); err != nil {
			return err
		}
		fp = filepath.Join(fp, "katapult.yaml")
		if err := ioutil.WriteFile(fp, []byte{}, 0o600); err != nil {
			return err
		}
		return c.WriteConfig()
	default:
		return err
	}
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("katapult")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/katapult")
	v.AddConfigPath("$HOME/.katapult")

	v.SetEnvPrefix("katapult")

	return v
}
