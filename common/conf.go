package common

import "code.infervision.com/common/go/io"

// Conf ...
type Conf struct {
	StorageDB    bool   `yaml:"storage_db"`
	StorageExcel bool   `yaml:"storage_excel"`
	DB           DBConf `yaml:"db"`
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
