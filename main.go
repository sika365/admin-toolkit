package main

import (
	"flag"
	"path"

	simutils "github.com/alifakhimi/simple-utils-go"

	"github.com/sika365/admin-tools/service"
)

func main() {
	var (
		configPath string
	)

	flag.StringVar(&configPath, "c", path.Join(simutils.CurrentDirectory(), "config.json"),
		// usage
		"config path with json extension")

	flag.Parse()

	service.New(configPath)
}
