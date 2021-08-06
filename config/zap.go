package config

type Zap struct {
	File 		  string `mapstructure:"file" json:"file"  yaml:"file"`
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	MaxSize  	  int `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"` // bytes
	MaxBackups    int `mapstructure:"maxbackups" json:"maxbackups" yaml:"maxbackups"`
	MaxAge        int `mapstructure:"maxage" json:"maxage" yaml:"maxage"`
	Compress      bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
