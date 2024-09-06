package main

import (
	"context"

	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/osquery/osquery-go/plugin/table"
)

func createCinderVolumesTable() *table.Plugin {
	columns := []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("size"),
		table.TextColumn("status"),
		table.TextColumn("created_at"),
		table.TextColumn("project_id"),
		table.TextColumn("cloud_name"),
	}

	return table.NewPlugin("cinder_volumes", columns, generateCinderVolumes)
}

func generateCinderVolumes(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	clients, err := getClientsFromCloudsYAML()
	if err != nil {
		return nil, err
	}

	var results []map[string]string

	for _, client := range clients {
		blockStorageClient, err := createBlockStorageClient(client.Client, client.ProjectID)
		if err != nil {
			return nil, err
		}

		pager := volumes.List(blockStorageClient, volumes.ListOpts{})
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			allVolumes, err := volumes.ExtractVolumes(page)
			if err != nil {
				return false, err
			}
			for _, volume := range allVolumes {
				results = append(results, map[string]string{
				  "id":          volume.ID,
				  "name":        volume.Name,
				  "status":      volume.Status,
				  "size":        fmt.Sprintf("%d", volume.Size),
				  "volume_type": volume.VolumeType,
				  "created_at":  volume.CreatedAt,
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

