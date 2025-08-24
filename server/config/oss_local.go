package config

type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`                   // 本地文件访问路径
	StorePath string `mapstructure:"store-path" json:"store-path" yaml:"store-path"` // 本地文件存储路径
	BaseURL   string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
}
