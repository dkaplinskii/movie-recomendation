package modules

type GRPC struct {
	Port int `env:"PORT" envDefault:"50051"`
}
