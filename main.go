package main

import (
	"runtime"

	"github.com/codegangsta/negroni"
	"github.com/docopt/docopt-go"
	"github.com/meatballhat/negroni-logrus"

	"github.com/paddycarey/ims/cache"
	"github.com/paddycarey/ims/server"
	"github.com/paddycarey/ims/storage"
)

var usage string = `ims.

Usage:
  ims [--storage=<src>] [--storage-credentials=<creds>] [--cache=<cch>] [--address=<address>]
  ims -h | --help
  ims --version

Options:
  -h --help                      Show this screen.
  --version                      Show version.
  --storage=<src>                Storage backend                  [default: ./].
  --storage-credentials=<creds>  Storage credentials JSON file    [default: storage-credentials.json].
  --cache=<cch>                  Cache backend                    [default: ./.cache].
  --address=<address>            Address that ims should bind to  [default: :5995].`

func main() {

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse command line args, exiting if required
	arguments, _ := docopt.Parse(usage, nil, true, "ims 0.1", false)

	// load cache backend
	c, err := cache.LoadBackend(arguments["--cache"].(string))
	if err != nil {
		panic(err)
	}

	// load storage backend
	s, err := storage.LoadBackend(arguments["--storage"].(string), arguments["--storage-credentials"].(string))
	if err != nil {
		panic(err)
	}

	// run application server
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewRecovery())
	n.UseHandler(&server.Server{Cache: c, Storage: s})
	n.Run(arguments["--address"].(string))
}
