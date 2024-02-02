package setting

type Settings struct {
	Locale            string `yaml:"locale"`
	DataDir           string `yaml:"data_dir"`
	DeleteAfterClosed bool   `yaml:"delete_after_closed"`
}
