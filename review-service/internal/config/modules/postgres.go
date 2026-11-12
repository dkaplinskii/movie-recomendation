package modules

type Postgres struct {
	Host           string `env:"HOST,required"`
	Port           int    `env:"PORT" envDefault:"5432"`
	User           string `env:"USER,required"`
	Password       string `env:"PASSWORD,required"`
	DBName         string `env:"DB_NAME,required"`
	SSLMode        string `env:"SSLMODE" envDefault:"disable"`
	MigrationsPath string `env:"MIGRATIONS_PATH" envDefault:"./migrations"`
}
