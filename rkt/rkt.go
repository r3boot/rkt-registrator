package rkt

import (
	"github.com/r3boot/rkt-registrator/utils"
	"os"
)

var Log utils.Log
var Rkt_dir string
var Cni_dir string

func Setup(l utils.Log, rkt_dir string, cni_dir string) (err error) {
	Log = l

	// Sanity check
	if _, err = os.Stat(rkt_dir); err != nil {
		Log.Fatal(rkt_dir + " does not exist")
	}

	if _, err = os.Stat(cni_dir); err != nil {
		Log.Fatal(cni_dir + " does not exist")
	}

	if _, err = os.Stat(rkt_dir + "/pods"); err != nil {
		Log.Fatal(rkt_dir + "/pods does not exist, forgot to run setup-data-dir.sh?")
	}

	Rkt_dir = rkt_dir
	Cni_dir = cni_dir

	return
}
