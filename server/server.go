package server

import (
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"

	"github.com/paddycarey/ims/storage"
)

func NewServer(source storage.FileSystem, cacheDir string) *negroni.Negroni {

	return negroni.New(
		negronilogrus.NewMiddleware(),
		negroni.NewRecovery(),
		NewCacheMiddleware(cacheDir),
		NewFilterMiddleware(source, cacheDir),
	)
}
