package envme

import (
	"context"
	"envme/lib/docker"
	"envme/lib/types"
	"envme/lib/utils"
	"fmt"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func CreateService(ctx context.Context, name, image, network string) error {
	// Create a new Docker Compose file
	config := &types.Compose{
		Services: map[string]*types.Service{
			name: {
				ContainerName: name,
				Image:         image,
				Restart:       "unless-stopped",
				Environment:   viper.GetStringSlice("env"),
				Networks:      &[]string{network},
				ExtraHosts:    &[]string{"host.docker.internal:host-gateway"},
			},
		},
		Networks: map[string]*types.Network{
			network: {
				External: true,
			},
		},
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Error marshalling config: %v\n", err)
		return err
	}

	err = utils.WriteComposeFile(name, content)
	if err != nil {
		fmt.Printf("Error writing compose file: %v\n", err)
		return err
	}

	// Run the Docker Compose file
	compose, project, err := docker.NewCompose(ctx, name)
	if err != nil {
		fmt.Printf("Error creating compose: %v\n", err)
		return err
	}

	// TODO: Expose the port with tunneling

	return compose.Up(ctx, project, api.UpOptions{})
}

func CreateDev(ctx context.Context, name, dir, template, network string) error {
	// Create a new Docker Compose file
	dir, err := utils.GetAbsPath(dir)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		return err
	}
	fmt.Printf("Creating development environment for %s in %s\n", name, dir)
	config := &types.Compose{
		Services: map[string]*types.Service{
			name: {
				ContainerName: name,
				Build: &types.Build{
					Context:    dir,
					Dockerfile: "Dockerfile",
					Target:     "development",
				},
				Restart:     "unless-stopped",
				Environment: viper.GetStringSlice("env"),
				Networks:    &[]string{network},
				ExtraHosts:  &[]string{"host.docker.internal:host-gateway"},
			},
		},
		Networks: map[string]*types.Network{
			network: {
				External: true,
			},
		},
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Error marshalling config: %v\n", err)
		return err
	}

	err = utils.WriteComposeFile(name, content)
	if err != nil {
		fmt.Printf("Error writing compose file: %v\n", err)
		return err
	}

	// Write dockerfile when template is not empty
	if template != "" && template != "(none)" {
		err = utils.WriteDockerfile(dir, template)
		if err != nil {
			fmt.Printf("Error writing Dockerfile: %v\n", err)
			return err
		}
	}

	// Run the Docker Compose file
	compose, project, err := docker.NewCompose(ctx, name)
	if err != nil {
		fmt.Printf("Error creating compose: %v\n", err)
		return err
	}

	// TODO: Expose the port with tunneling

	return compose.Up(ctx, project, api.UpOptions{})
}
