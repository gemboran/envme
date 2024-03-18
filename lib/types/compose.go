package types

type Compose struct {
	Version  string              `yaml:"version,omitempty"`
	Services map[string]*Service `yaml:"services"`
	Networks map[string]*Network `yaml:"networks,omitempty"`
	Volumes  map[string]*Volume  `yaml:"volumes,omitempty"`
}

type Service struct {
	ContainerName string    `yaml:"container_name,omitempty"`
	Build         *Build    `yaml:"build,omitempty"`
	Image         string    `yaml:"image,omitempty"`
	Restart       string    `yaml:"restart,omitempty"`
	Volumes       []string  `yaml:"volumes,omitempty"`
	Environment   []string  `yaml:"environment,omitempty"`
	Command       string    `yaml:"command,omitempty"`
	Networks      *[]string `yaml:"networks,omitempty"`
	ExtraHosts    *[]string `yaml:"extra_hosts,omitempty"`
}

type Build struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
	Target     string `yaml:"target,omitempty"`
}

type Network struct {
	Driver   string `yaml:"driver,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type Volume string
