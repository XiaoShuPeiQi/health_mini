package config

type Postgresql struct {
	DriverName string `mapstructure:"driverName" json:"driverName" yaml:"driverName"`
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Port       string `mapstructure:"port" json:"port" yaml:"port"`
	Database   string `mapstructure:"database" json:"database" yaml:"database"`
	Username   string `mapstructure:"username" json:"username" yaml:"username"`
	Password   string `mapstructure:"password" json:"password" yaml:"password"`
	DbSslmode  string `mapstructure:"dbSslmode" json:"dbSslmode" yaml:"dbSslmode"`
}
