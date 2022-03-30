package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"os"
)

func InitConfig(versionCheck bool) (cfg *Config, err error) {
	log.Infof("config load [OK]")
	cfg = &Config{}
	err = cfg.readFile()
	if err != nil {
		return
	}

	if versionCheck {
		log.Infof("skipping secret config")
		return
	}

	err = cfg.readSecret()
	if err != nil {
		return
	}

	return
}

func (cfg *Config) readSecret() (err error) {
	var (
		sec *SecretParam
		f   *os.File
	)
	log.Infof("[config][readSecret] load secret param from %s", fmt.Sprintf(`%ssecret.param.yml`, cfg.Server.SecretPath))
	f, err = os.Open(fmt.Sprintf(`%ssecret.param.yml`, cfg.Server.SecretPath))
	if err != nil {
		return
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&sec)
	if err != nil {
		return
	}

	//set db auth param
	for k, v := range cfg.Database {
		v.DBAuth = DBAuthParam{
			MasterAuthUname: sec.DbAuth[k].MasterAuthUname,
			MasterAuthPass:  sec.DbAuth[k].MasterAuthPass,
			SlaveAuthUname:  sec.DbAuth[k].SlaveAuthUname,
			SlaveAuthPass:   sec.DbAuth[k].SlaveAuthPass,
		}
	}

	cfg.Mail = sec.Mail
	log.Info("[config][readSecret] Load Secret Param success")
	return
}

func (cfg *Config) readFile() (err error) {
	var (
		f *os.File
	)

	cfg.Environment = os.Getenv(enum.APP_ENV)
	if cfg.Environment != "staging" && cfg.Environment != "production" {
		cfg.Environment = "development"
	}
	log.Infof("[cfg.Environment ADALAH] %s", cfg.Environment)

	path := []string{
		"/etc/api-config",
		"files/etc/api-config",
		"./files/etc/api-config",
		"../files/etc/api-config",
		"../../files/etc/api-config",
	}

	for _, val := range path {
		f, err = os.Open(fmt.Sprintf(`%s/%s/config.main.yml`, val, cfg.Environment))
		if err == nil {
			log.Infof("[config][init] load config file from %s", fmt.Sprintf(`%s/%s/config.main.yml`, val, cfg.Environment))
			decoder := yaml.NewDecoder(f)
			err = decoder.Decode(cfg)
			break
		}
	}

	if err != nil {
		return
	}
	cfg.setLocalAddress()
	log.Infof("[config][ReadConfig] Config load success, running on \"%s\".", cfg.Environment)
	return
}

func (cfg *Config) GetLocalAddress() string {
	return cfg.Server.LocalAddress
}

func (cfg *Config) setLocalAddress() {
	// Once Set , we do not want to Reset it again
	if cfg.Server.LocalAddress != "" {
		return
	}

	localIP, err := cfg.getOutboundIP()
	if err != nil {
		log.Errorf("[config][SetLocalIP] Error Fetching IP. Error : %+v", err)
		return
	}

	cfg.Server.LocalAddress = localIP

}

func (cfg *Config) getOutboundIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
