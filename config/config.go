package config

type Server struct {
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	// gorm
	Mysql      Mysql      `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Postgresql Postgresql `mapstructure:"postgresql" json:"postgresql" yaml:"postgresql"`
	// oss
	Local Local `mapstructure:"local" json:"local" yaml:"local"`
}
