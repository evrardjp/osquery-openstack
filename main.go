package main

import (
	"fmt"
	"os"

	"github.com/osquery/osquery-go"
)

func main() {
	// Set up the osquery extension manager client
	server, err := osquery.NewExtensionManagerServer("openstack_extension", os.Getenv("OSQUERY_SOCKET"))
	if err != nil {
		fmt.Println("Error creating osquery extension:", err)
		os.Exit(1)
	}

	// Register the tables for OpenStack services
	server.RegisterPlugin(createNovaInstancesTable())
	server.RegisterPlugin(createNeutronNetworksTable())
	server.RegisterPlugin(createCinderVolumesTable())
	server.RegisterPlugin(createHeatStacksTable())

	// Start serving the extension
	if err := server.Run(); err != nil {
		fmt.Println("Error running osquery extension:", err)
	}
}

