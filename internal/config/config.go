package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	config_path_flag = "config-path"
	config_path_env  = "CONFIG_PATH"
)

var (
	path = flag.String(config_path_flag, "", "path to configure file")
)

type Cfg struct {
	Server struct {
		Host string `yaml:"host" env-required:"true"`
		Port string `yaml:"port" env-required:"true"`
		Env  string `yaml:"env" env-default:"local"`
	} `yaml:"server" env-required:"true"`

	AuthSecret string `yaml:"auth-secret" env-required:"true"`

	Postgres struct {
		Port     string `yaml:"port" env-required:"true"`
		Host     string `yaml:"host" env-required:"true"`
		Name     string `yaml:"name" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		User     string `yaml:"user" env-required:"true"`
		Migrate  bool   `yaml:"migrate" env-default:"false"`
		Sslmode  string `yaml:"sslmode" env-default:"disable"`
	} `yaml:"postgres" env-required:"true"`

	Notifyer struct {
		SmtpHost string `yaml:"smtp-host" env-required:"false"`
		SmtpPort string `yaml:"smtp-port" env-required:"false"`
		Email    string `yaml:"email" env-required:"false"`
		Password string `yaml:"password" env-required:"false"`
	} `yaml:"notifyer" env-required:"false"`
}

func LoadConfig() Cfg {
	cfg := &Cfg{}

	flag.Parse()
	if err := cleanenv.ReadConfig(configPath(), cfg); err != nil {
		log.Fatalf("[error] open configurate file %v", err)
	}
	return *cfg
}

func configPath() string {
	if *path == "" {
		*path = os.Getenv(config_path_env)
	}
	if _, err := os.Stat(*path); err == os.ErrNotExist {
		log.Fatalf("[error] no such configurate file %s", *path)
	}
	return *path
}
