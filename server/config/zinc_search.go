package config

type ZincSearch struct {
	Url      string `mapstructure:"url" json:"url" yaml:"url"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"Password" yaml:"Password"` // URL前缀
}
