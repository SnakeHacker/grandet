package common

//DBConf ...
type DBConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

func (c DBConf) validate() error {
	if c.Username == "" {
		return ErrDBEmptyUsername
	}
	if c.Host == "" {
		return ErrDBEmptyHost
	}
	if c.Port == 0 {
		return ErrDBEmptyPort
	}
	if c.Database == "" {
		return ErrDBEmptyDatabase
	}

	return nil
}
