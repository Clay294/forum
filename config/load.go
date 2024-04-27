package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Clay294/forum/flog"
	"github.com/caarlos0/env/v10"
)

func LoadConfig(path string) error {
	conf := new(Config)

	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		flog.Flogger().Error().Msgf("loading config from toml file failed: %s", err)
		return fmt.Errorf("loading config from toml file failed: %s", err)
	}

	err = env.Parse(conf)
	if err != nil {
		flog.Flogger().Error().Msgf("loading config from env failed: %s", err)
		return fmt.Errorf("loading config from env failed: %s", err)
	}

	globalConf = conf

	return nil
}
