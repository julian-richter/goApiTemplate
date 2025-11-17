package app

type Env string

const (
	EnvProd Env = "prod"
	EnvDev  Env = "dev"
)

func (e Env) Valid() bool {
	return e == EnvProd || e == EnvDev
}

type Config struct {
	ApplicationName string
	Port            string
	Env             Env
}
