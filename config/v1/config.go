package v1

type Config struct {
	XPackages []XPackageConfig `json:"xpackages"`
	Dist      string           `json:"dist"`
	Dockers   []DockerConfig   `json:"dockers"`
}

type XPackageConfig struct {
	ID       string `json:"id"`
	Dir      string `json:"dir"`
	Examples string `json:"examples"`
}

type DockerConfig struct {
	IDs            []string            `json:"ids"`
	ImageTemplates []string            `json:"image_templates"`
	Logins         []DockerConfigLogin `json:"logins"`
}

type DockerConfigLoginType string

const (
	DockerConfigLoginTypeGoogle DockerConfigLoginType = "google"
	DockerConfigLoginTypeAzure  DockerConfigLoginType = "azure"
	DockerConfigLoginTypeAWS    DockerConfigLoginType = "aws"
	DockerConfigLoginTypeBasic  DockerConfigLoginType = "basic"
)

type DockerConfigLogin struct {
	Type     DockerConfigLoginType  `json:"type"`
	Registry string                 `json:"registry"`
	Basic    DockerConfigLoginBasic `json:"basic"`
}

type DockerConfigLoginBasic struct {
	UsernameFromEnv string `json:"usernameFromEnv"`
	PasswordFromEnv string `json:"passwordFromEnv"`
}
