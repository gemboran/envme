package envme

import (
	"context"
	"envme/lib/docker"
	"envme/lib/types"
	"envme/lib/utils"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func CreateService(ctx context.Context, name, image, network string) error {
	// Create a new Docker Compose file
	config := &types.Compose{
		Version: "3.8",
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
		return err
	}

	err = utils.WriteComposeFile(name, content)
	if err != nil {
		return err
	}

	// Run the Docker Compose file
	compose, project, err := docker.NewCompose(ctx, name)
	if err != nil {
		return err
	}

	// TODO: Expose the port with tunneling

	return compose.Up(ctx, project, api.UpOptions{})
}
