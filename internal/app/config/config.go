package config

var Cfg *Config

type Config struct {
	Database struct {
		DSN         string `mapstructure:"dsn"`
		MaxIdleConn int    `mapstructure:"max_idle_conn"`
		MaxOpenConn int    `mapstructure:"max_open_conn"`
	} `mapstructure:"database"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	// 其他配置...
}

// 加载配置
func LoadConfig() (*Config, error) {
	var cfg Config
	// 加载配置逻辑...
	Cfg = &cfg
	return &cfg, nil
}

/*
// Get 从环境变量中获取指定键的值
func Get(key string) string {
	return os.Getenv(key)
}*/
