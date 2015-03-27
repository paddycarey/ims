package main

import (
	"net/http"
	"runtime"

	"github.com/codegangsta/negroni"
	"github.com/docopt/docopt-go"
	"github.com/meatballhat/negroni-logrus"

	"github.com/paddycarey/ims/middleware"
)

var usage string = `ims.

Usage:
  ims [--source=<src>] [--cache=<cache>]
  ims -h | --help
  ims --version

Options:
  -h --help        Show this screen.
  --version        Show version.
  --source=<src>   Source directory  [default: ./].
  --cache=<cache>  Cache directory   [default: ./.cache/].`

func main() {

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse command line args, exiting if required
	arguments, _ := docopt.Parse(usage, nil, true, "ims 0.1", false)

	// set up filesystems for source and cache directories
	cacheDir := arguments["--cache"].(string)
	cacheFS := http.Dir(cacheDir)
	srcFS := http.Dir(arguments["--source"].(string))

	n := negroni.New(
		negronilogrus.NewMiddleware(),
		negroni.NewRecovery(),
		middleware.NewCache(cacheFS),
		middleware.NewProcessor(srcFS, cacheDir),
		middleware.NewNotFound(),
	)
	n.Run(":5995")

}
