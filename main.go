package main

import (
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/docopt/docopt-go"
	"github.com/meatballhat/negroni-logrus"

	"github.com/paddycarey/ims/cache"
	"github.com/paddycarey/ims/server"
	"github.com/paddycarey/ims/storage"
)

var usage string = `ims.

Usage:
  ims [--storage=<src>] [--storage-credentials=<creds>] [--cache=<cch>] [--address=<address>] [--log-level=<level>] [--no-optimization]
  ims -h | --help
  ims --version

Options:
  -h --help                      Show this screen.
  --version                      Show version.
  --storage=<src>                Storage backend                        [default: ./].
  --storage-credentials=<creds>  Storage credentials JSON file          [default: storage-credentials.json].
  --cache=<cch>                  Cache backend                          [default: ./.cache].
  --address=<address>            Address that ims should bind to        [default: :5995].
  --log-level=<level>            Logging level (debug/info/warn/error)  [default: info].
  --no-optimization              Disables image optimization.`

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
	arguments, _ := docopt.Parse(usage, nil, true, "ims 0.1", false)

	// configure global logger instance
	logLevel, err := logrus.ParseLevel(arguments["--log-level"].(string))
	exitOnError(err, "Unable to initialise logger")
	logrus.SetLevel(logLevel)

	// load cache backend
	c, err := cache.LoadBackend(arguments["--cache"].(string))
	exitOnError(err, "Unable to load cache backend")

	// load storage backend
	s, err := storage.LoadBackend(arguments["--storage"].(string), arguments["--storage-credentials"].(string))
	exitOnError(err, "Unable to load storage backend")

	// run application server
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewRecovery())
	n.UseHandler(&server.Server{Cache: c, Storage: s, NoOpts: arguments["--no-optimization"].(bool)})
	n.Run(arguments["--address"].(string))
}
