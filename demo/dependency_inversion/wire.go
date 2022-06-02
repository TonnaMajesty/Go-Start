// +build wireinject

package dependency_inversion

import "github.com/google/wire"

func initApp() *App {
	panic(wire.Build(NewRedis, NewConfig, NewApp))
}
