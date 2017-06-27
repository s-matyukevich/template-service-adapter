package main

import (
	"log"
	"os"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/s-matyukevich/template-service-adapter/adapter"

	"github.com/s-matyukevich/template-service-adapter/config"
)

func main() {
	stderrLogger := log.New(os.Stderr, "[template-service-adapter] ", log.LstdFlags)
	args := os.Args
	if args[1] != "--config" {
		stderrLogger.Fatal("Config argument is not provided")
	}
	config, err := config.ParseConfig(args[2])
	if err != nil {
		stderrLogger.Fatal("Error while parsing config:", err)
	}
	manifestGenerator := adapter.ManifestGenerator{Config: config}
	binder := adapter.Binder{Config: config}
	args = append([]string{args[0]}, args[3:]...)
	serviceadapter.HandleCommandLineInvocation(args, manifestGenerator, binder, nil)
}
