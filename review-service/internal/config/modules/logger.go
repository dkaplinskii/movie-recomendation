package modules

type LoggerCfg struct {
	Level string `env:"LEVEL" envDefault:"info"`
}
