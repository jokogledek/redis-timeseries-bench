package config

type Config struct {
	Redis    RedisConfig    `yaml:"redis"`
	Database DatabaseConfig `yaml:"database"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Pass     string `yaml:"pass"`
	DBName   string `yaml:"dbname"`
	Host     string `yaml:"host"`
}
