package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/nickrobison/terraform-linux-provider/provider/internal/provider"
)

var (
	version string = "dev"
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")

	opts := providerserver.ServeOpts{
		Address: "terraform.nickrobison.com/nickrobison/linux",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err)
	}

}
