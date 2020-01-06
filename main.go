package main

import (
	"fmt"
	"os"

	"github.com/VisImag/katy/env"
)

func main() {
	env.SetEnv()
	fmt.Println(os.Getenv("KUBE_CONFIG_PATH"))
}
