package v1

type Config struct {
	ProjectName string           `json:"project_name"`
	XPackages   []XPackageConfig `json:"xpackages"`
	Dist        string           `json:"dist"`
	Dockers     []DockerConfig   `json:"dockers"`
}

type XPackageConfig struct {
	ID           string `json:"id"`
	Dir          string `json:"dir"`
	Examples     string `json:"examples"`
	NameTemplate string `json:"name_template"`
}

type DockerConfig struct {
	IDs            []string `json:"ids"`
	ImageTemplates []string `json:"image_templates"`
}
