package main

import (
	"context"

	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/osquery/osquery-go/plugin/table"
)

func createNovaInstancesTable() *table.Plugin {
	columns := []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		table.TextColumn("flavor"),
		table.TextColumn("image"),
		table.TextColumn("created_at"),
		table.TextColumn("project_id"),
		table.TextColumn("cloud_name"),
	}

	return table.NewPlugin("nova_instances", columns, generateNovaInstances)
}

func generateNovaInstances(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	clients, err := getClientsFromCloudsYAML()
	if err != nil {
		return nil, err
	}

	var results []map[string]string

	for _, client := range clients {
		computeClient, err := createComputeClient(client.Client, client.ProjectID)
		if err != nil {
			return nil, err
		}

		pager := servers.List(computeClient, servers.ListOpts{})
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			allServers, err := servers.ExtractServers(page)
			if err != nil {
				return false, err
			}
			for _, server := range allServers {
				results = append(results, map[string]string{
					"id":         server.ID,
					"name":       server.Name,
					"status":     server.Status,
					"flavor":     server.Flavor["id"].(string),
					"image":      server.Image["id"].(string),
					"created_at": server.Created.String(),
					"project_id": client.ProjectID,
					"cloud_name": client.CloudName,
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

