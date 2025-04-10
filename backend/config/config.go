package config

import "os"

// Get 从环境变量中获取指定键的值
func Get(key string) string {
	return os.Getenv(key)
}
