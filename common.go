package main

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

type OpenStackClient struct {
	CloudName   string
	ProjectID   string
	ProjectName string
	Client      *gophercloud.ProviderClient
}

// Parse clouds.yaml and create a client for each cloud/project
func getClientsFromCloudsYAML() ([]OpenStackClient, error) {
	var clients []OpenStackClient

	opts := &clientconfig.ClientOpts{}
	clouds, err := clientconfig.GetClouds(opts)
	if err != nil {
		return nil, err
	}

	for cloudName, _ := range clouds {
		authOpts, err := clientconfig.AuthOptions(&clientconfig.ClientOpts{
			Cloud: cloudName,
		})
		if err != nil {
			return nil, err
		}

		provider, err := openstack.AuthenticatedClient(*authOpts)
		if err != nil {
			return nil, err
		}

		identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
		if err != nil {
			return nil, err
		}

		projectList := projects.List(identityClient, projects.ListOpts{})
		err = projectList.EachPage(func(page pagination.Page) (bool, error) {
			allProjects, err := projects.ExtractProjects(page)
			if err != nil {
				return false, err
			}

			for _, project := range allProjects {
				clients = append(clients, OpenStackClient{
					CloudName:   cloudName,
					ProjectID:   project.ID,
					ProjectName: project.Name,
					Client:      provider,
				})
			}

			return true, nil
		})
		if err != nil {
			return nil, err
		}
	}

	return clients, nil
}

// Create service clients for each project
func createComputeClient(provider *gophercloud.ProviderClient, projectID string) (*gophercloud.ServiceClient, error) {
	return openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		ProjectID: projectID,
	})
}

func createNetworkingClient(provider *gophercloud.ProviderClient, projectID string) (*gophercloud.ServiceClient, error) {
	return openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		ProjectID: projectID,
	})
}

func createBlockStorageClient(provider *gophercloud.ProviderClient, projectID string) (*gophercloud.ServiceClient, error) {
	return openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		ProjectID: projectID,
	})
}

func createOrchestrationClient(provider *gophercloud.ProviderClient, projectID string) (*gophercloud.ServiceClient, error) {
	return openstack.NewOrchestrationV1(provider, gophercloud.EndpointOpts{
		ProjectID: projectID,
	})
}

