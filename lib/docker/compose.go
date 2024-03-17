package docker

import (
	"context"
	"envme/lib/utils"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"path/filepath"
	"strings"
)

// NewCompose is a function that creates a Docker project from a given Docker Compose file.
// It takes a context and a string representation of the Docker Compose file as input.
// It returns a pointer to a types.Project which represents the Docker project.
func NewCompose(ctx context.Context, stackName string) (api.Service, *types.Project, error) {
	dir, err := utils.GetServiceDir(stackName)
	if err != nil {
		return nil, nil, err
	}

	file := filepath.Join(dir, "docker-compose.yaml")
	configDetails := types.ConfigDetails{
		WorkingDir:  dir,
		ConfigFiles: []types.ConfigFile{{Filename: file}},
		Environment: utils.ConvertEnvToMap(),
	}

	project, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(stackName, true)
	})
	if err != nil {
		return nil, nil, err
	}

	addServiceLabels(project)

	service, err := createService()
	if err != nil {
		return nil, nil, err
	}

	return service, project, nil
}

// createService is a function that creates a Docker service from a given Docker Compose service.
// It takes a pointer to a types.Service as input.
// It returns a pointer to a types.Service which represents the Docker service.
func createService() (api.Service, error) {
	var srv api.Service
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}

	dockerContext := "default"

	opts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(opts)
	if err != nil {
		return nil, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, err
}

// addServiceLabels adds the labels docker compose expects to exist on services.
// This is required for future compose operations to work, such as finding
// containers that are part of a service.
func addServiceLabels(project *types.Project) {
	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False", // default, will be overridden by `run` command
		}
		project.Services[i] = s
	}
}
