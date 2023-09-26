package config

type Config struct {
	ServerPort        string
	DbUri             string
	DbName            string
	BlogCollection    string
	TemplatesLocation string
}

func NewConfig(port, dbUri, dbName, blogCollection, templatesLocation string) *Config {
	return &Config{
		ServerPort:        port,
		DbUri:             dbUri,
		DbName:            dbName,
		BlogCollection:    blogCollection,
		TemplatesLocation: templatesLocation,
	}
}
