//+build wireinject

package app

import (
	"database/sql"
	"github.com/google/wire"
)

var wireSet = wire.NewSet(
	NewApp,
	NewRouter,
	v1.WireSet,
	service.WireSet,
	store.WireSet,
	wire.Bind(new(*service.UrlShortenerService), new(*store.URLShortenerStore)),
)

var wireSetWithDb = wire.NewSet(
	wireSet,
	db.WireSet,
)

func Initialize() (*App, error) {
	panic(wire.Build(wireSetWithDb))
}

func InitializeWithDb(*sql.DB) (*App, error) {
	panic(wire.Build(wireSet))
}
