package stan

type Config struct {
	Url         string `yaml:"url"`
	ClusterName string `yaml:"cluster_name"`
	ClientId    string `yaml:"client_id"`
}
