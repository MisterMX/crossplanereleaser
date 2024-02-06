package v1

type Config struct {
	XPackages []XPackageConfig `json:"xpackages"`
	Dist      string           `json:"dist"`
}

type XPackageConfig struct {
	ID       string `json:"id"`
	Dir      string `json:"dir"`
	Examples string `json:"examples"`
}
