package server

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
)

func NewServer(source http.FileSystem, cacheDir string) *negroni.Negroni {

	return negroni.New(
		negronilogrus.NewMiddleware(),
		negroni.NewRecovery(),
		NewCacheMiddleware(cacheDir),
		NewFilterMiddleware(source, cacheDir),
	)
}
