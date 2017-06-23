package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/s-matyukevich/template-service-adapter/adapter"

	"github.com/s-matyukevich/template-service-adapter/config"
)

func main() {
	stderrLogger := log.New(os.Stderr, "[template-service-adapter] ", log.LstdFlags)
	ex, err := os.Executable()
	if err != nil {
		stderrLogger.Fatal(err)
	}
	config, err := config.ParseConfig(filepath.Join(path.Dir(ex), "../config/config.yml"))
	if err != nil {
		stderrLogger.Fatal("config", err)
	}
	manifestGenerator := adapter.ManifestGenerator{Config: config}
	binder := adapter.Binder{Config: config}
	serviceadapter.HandleCommandLineInvocation(os.Args, manifestGenerator, binder, nil)
}
