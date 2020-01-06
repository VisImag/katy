package env

import (
	"log"
	"os"
	"path/filepath"
)

// User home directory
var HOMEDIR, err = os.UserHomeDir()

func SetEnv() {

	// Check for error in fetching HOMEDIR
	if err != nil {
		log.Fatalln("Error fetching $HOME")
	}

	// Set default absolute path for kubeconfig
	if _, ok := os.LookupEnv("KUBE_CONFIG_PATH"); !ok {
		path := filepath.Join(HOMEDIR, "/.kube/config")
		os.Setenv("KUBE_CONFIG_PATH", path)
	}
}
