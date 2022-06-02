package dependency_inversion


type Datasource interface {
	GetById(id string) string
}

type App struct {
	ds Datasource
}

func (a *App) getDataSource() *Datasource {
	return &a.ds
}

func NewApp(ds Datasource) *App {
	return &App{ds: ds}
}

var _ Datasource = (*Redis)(nil)

type Redis struct {
	c *Config
}

func (r *Redis) GetById(id string) string {
	return id
}

func NewRedis(c *Config) Datasource {
	return &Redis{c}
}

func NewConfig() *Config {
	return &Config{ip: "123.2344.323", port: "8080"}
}

type Config struct {
	ip string
	port string
}




