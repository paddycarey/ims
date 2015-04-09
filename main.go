package main

import (
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"

	"github.com/paddycarey/ims/cache"
	"github.com/paddycarey/ims/config"
	"github.com/paddycarey/ims/server"
	"github.com/paddycarey/ims/storage"
)

// exitOnError checks that an error is not nil. If the passed value is an
// error, it is logged and the program exits with an error code of 1
func exitOnError(err error, prefix string) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal(prefix)
	}
}

func main() {

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse command line args, exiting if required
	args := config.ParseCLIConfig()

	// configure global logger instance
	logLevel, err := logrus.ParseLevel(args.LogLevel)
	exitOnError(err, "Unable to initialise logger")
	logrus.SetLevel(logLevel)

	// load cache backend
	c, err := cache.LoadBackend(args.Cache)
	exitOnError(err, "Unable to load cache backend")

	// load storage backend
	s, err := storage.LoadBackend(args.Storage)
	exitOnError(err, "Unable to load storage backend")

	// run application server
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewRecovery())
	n.UseHandler(&server.Server{Cache: c, Storage: s, NoOpts: args.NoOptimization})
	n.Run(args.Address)
}
