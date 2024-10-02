package v1

type Config struct {
	ProjectName string        `json:"project_name"`
	Dist        string        `json:"dist"`
	Builds      []BuildConfig `json:"builds"`
	Pushes      []PushConfig  `json:"pushes"`
}

type BuildConfig struct {
	ID              string `json:"id"`
	Dir             string `json:"dir"`
	Examples        string `json:"examples"`
	NameTemplate    string `json:"name_template"`
	RuntimeImageTar string `json:"runtime_image_tar"`
}

type PushConfig struct {
	ID             string   `json:"id"`
	Build          string   `json:"build"`
	ImageTemplates []string `json:"image_templates"`
}
