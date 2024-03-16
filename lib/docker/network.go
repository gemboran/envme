package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func EnsureNetworkExists(ctx context.Context, client *client.Client, networkName string) error {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", networkName)

	networks, err := client.NetworkList(ctx, types.NetworkListOptions{Filters: filterArgs})
	if err != nil {
		return err
	}

	if len(networks) == 0 {
		_, err := client.NetworkCreate(ctx, networkName, types.NetworkCreate{})
		if err != nil {
			return err
		}
	}

	return nil
}
