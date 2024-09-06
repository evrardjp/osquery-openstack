package main

import (
	"context"

	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/osquery/osquery-go/plugin/table"
)

func createHeatStacksTable() *table.Plugin {
	columns := []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("status"),
		table.TextColumn("creation_time"),
		table.TextColumn("updated_time"),
		table.TextColumn("stack_status_reason"),
		table.TextColumn("project_id"),
		table.TextColumn("cloud_name"),
	}

	return table.NewPlugin("heat_stacks", columns, generateHeatStacks)
}

func generateHeatStacks(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	clients, err := getClientsFromCloudsYAML()
	if err != nil {
		return nil, err
	}

	var results []map[string]string

	for _, client := range clients {
		orchestrationClient, err := createOrchestrationClient(client.Client, client.ProjectID)
		if err != nil {
			return nil, err
		}

		pager := stacks.List(orchestrationClient, stacks.ListOpts{})
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			allStacks, err := stacks.ExtractStacks(page)
			if err != nil {
				return false, err
			}
			for _, stack := range allStacks {
				results = append(results, map[string]string{
					"id":                  stack.ID,
					"name":                stack.Name,
					"status":              stack.StackStatus,
					"creation_time":       stack.CreationTime.String(),
					"updated_time":        stack.UpdatedTime.String(),
					"stack_status_reason": stack.StackStatusReason,
					"project_id":          client.ProjectID,
					"cloud_name":          client.CloudName,
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

