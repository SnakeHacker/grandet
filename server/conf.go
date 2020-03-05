package server

import (
	"time"

	"code.infervision.com/common/go/io"
)

// Conf ...
type Conf struct {
	StorageDB    bool   `yaml:"storage_db"`
	StorageExcel bool   `yaml:"storage_excel"`
	DB           DBConf `yaml:"db"`
	Web          struct {
		Port int `yaml:"port"`
		// Timeout config
		WriteTimeoutInSec time.Duration `yaml:"write_timeout_sec"`
		ReadTimeoutInSec  time.Duration `yaml:"read_timeout_sec"`
		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout_sec"`
		IdleTimeoutInSec  time.Duration `yaml:"idle_timeout_sec"`
		ShutdownWaitInSec time.Duration `yaml:"shutdown_wait_sec"`
	} `yaml:"web"`
}

// LoadConf load config from yaml
func LoadConf(yamlPath string) (conf Conf, err error) {
	if err = io.LoadYAML(yamlPath, &conf); err != nil {
		return
	}
	if err = conf.validate(); err != nil {
		return
	}
	return
}

func (c Conf) validate() error {
	return nil
}
