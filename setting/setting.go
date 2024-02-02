package setting

type Settings struct {
	Language          string `yaml:"language"`
	DataDir           string `yaml:"data_dir"`
	DeleteAfterClosed bool   `yaml:"delete_after_closed"`
}
