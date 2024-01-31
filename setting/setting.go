package setting

type Settings struct {
	DataDir           string `yaml:"data_dir"`
	DeleteAfterClosed bool   `yaml:"delete_after_closed"`
}
