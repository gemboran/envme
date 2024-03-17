package types

type Compose struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
	Networks map[string]*Network `yaml:"networks,omitempty"`
	Volumes  map[string]*Volume  `yaml:"volumes,omitempty"`
}

type Service struct {
	ContainerName string    `yaml:"container_name,omitempty"`
	Image         string    `yaml:"image"`
	Restart       string    `yaml:"restart,omitempty"`
	Volumes       []string  `yaml:"volumes,omitempty"`
	Environment   []string  `yaml:"environment,omitempty"`
	Command       string    `yaml:"command,omitempty"`
	Networks      *[]string `yaml:"networks,omitempty"`
	ExtraHosts    *[]string `yaml:"extra_hosts,omitempty"`
}

type Network struct {
	Driver   string `yaml:"driver,omitempty"`
	External bool   `yaml:"external,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type Volume string
