package main

import (
	"context"

	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/osquery/osquery-go/plugin/table"
)

func createNeutronNetworksTable() *table.Plugin {
	columns := []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		table.TextColumn("admin_state_up"),
		table.TextColumn("tenant_id"),
		table.TextColumn("project_id"),
		table.TextColumn("cloud_name"),
	}

	return table.NewPlugin("neutron_networks", columns, generateNeutronNetworks)
}

func generateNeutronNetworks(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	clients, err := getClientsFromCloudsYAML()
	if err != nil {
		return nil, err
	}

	var results []map[string]string

	for _, client := range clients {
		networkClient, err := createNetworkingClient(client.Client, client.ProjectID)
		if err != nil {
			return nil, err
		}

		pager := networks.List(networkClient, networks.ListOpts{})
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			allNetworks, err := networks.ExtractNetworks(page)
			if err != nil {
				return false, err
			}
			for _, network := range allNetworks {
				results = append(results, map[string]string{
					"id":            network.ID,
					"name":          network.Name,
					"status":        network.Status,
					"admin_state_up": fmt.Sprintf("%t", network.AdminStateUp),
					"tenant_id":     network.TenantID,
					"project_id":    client.ProjectID,
					"cloud_name":    client.CloudName,
				})
			}
			return true, nil
		})
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

