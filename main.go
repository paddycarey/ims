package main

import (
	"runtime"

	"github.com/docopt/docopt-go"

	"github.com/paddycarey/ims/server"
	"github.com/paddycarey/ims/storage"
)

var usage string = `ims.

Usage:
  ims [--source=<src>] [--cache=<cch>] [--address=<address>]
  ims -h | --help
  ims --version

Options:
  -h --help            Show this screen.
  --version            Show version.
  --source=<src>       Source directory                 [default: ./].
  --cache=<cch>        Cache directory                  [default: ./.cache].
  --address=<address>  Address that ims should bind to  [default: :5995].`

func main() {

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// parse command line args, exiting if required
	arguments, _ := docopt.Parse(usage, nil, true, "ims 0.1", false)

	// load storage backend
	src, err := storage.LoadBackend(arguments["--source"].(string))
	if err != nil {
		panic(err)
	}

	// run application server
	srv := server.NewServer(src, arguments["--cache"].(string))
	srv.Run(arguments["--address"].(string))
}
