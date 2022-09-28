package config

// Config describes the configuration of the service.
type Config struct {
	Application application `yaml:"application"`
	Settings    settings
}

type application struct {
	Authors []author
	Name    string
	Usage   string
	Version string
}

type author struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type settings struct {
	AllowedColumnNameCaseValues []string
	AllowedTableNameCaseValues  []string
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
			Version: "v0.7.0",
		},
		Settings: settings{
			AllowedColumnNameCaseValues: []string{"snake_case", "camelCase", "screaming_snake_case", "kebab_case"},
			AllowedTableNameCaseValues:  []string{"snake_case", "camelCase", "screaming_snake_case", "kebab_case"},
		},
	}
}
