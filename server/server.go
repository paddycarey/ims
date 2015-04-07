package server

import (
	"github.com/codegangsta/negroni"
	"github.com/pilu/xrequestid"

	"github.com/paddycarey/ims/storage"
)

func NewServer(source storage.FileSystem, cacheDir string) *negroni.Negroni {

	return negroni.New(
		xrequestid.New(16),
		NewLoggerMiddleware(),
		negroni.NewRecovery(),
		NewCacheMiddleware(cacheDir),
		NewFilterMiddleware(source, cacheDir),
	)
}
