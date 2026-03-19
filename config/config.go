package config

type MySQLConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Config struct {
	MySQL MySQLConfig
	Redis RedisConfig
	Email EmailConfig
}

var AppConfig = Config{
	MySQL: MySQLConfig{
		DSN: "root:123456@tcp(192.168.1.103)/shop?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr:     "192.168.1.103",
		Password: "123456",
		DB:       0,
	},
	Email: EmailConfig{
		Host:     "smtp.qq.com",
		Port:     587,
		Username: "235038369@qq.com",
		Password: "urthfqcguamfbhdf",
		From:     "235038369@qq.com",
	},
}
