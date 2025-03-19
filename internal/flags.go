package internal

import (
	"github.com/spf13/pflag"
)

var args Args

type Args struct {
	Debug bool
	V6dSocket string
	Backend string
}

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&args.V6dSocket, "v6d-socket", "/var/run/vineyard.sock", "v6d socket")
	fs.BoolVar(&args.Debug, "debug", false, "debug mode")
	fs.StringVar(&args.Backend, "backend", "mock", "backend type, mock or a object storage")
}

func GetArgs() Args {
	return args
}
