package config

// Config describes the configuration of the service.
type Config struct {
	Application application `yaml:"application"`
}

type application struct {
	Authors []author `yaml:"authors"`
	Name    string   `yaml:"name"`
	Usage   string   `yaml:"usage"`
	Version string   `yaml:"version"`
}

type author struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// New creates and returns a configuration object for the service.
func New() Config {
	return Config{
		Application: application{
			Authors: []author{
				{
					Name:  "Angelos Giannis",
					Email: "angelos.giannis@gmail.com",
				},
			},
			Name:    "ErBuilder",
			Usage:   "CLI tool to automatically generate a .er file including a list of model files",
			Version: "0.1.0",
		},
	}
}
