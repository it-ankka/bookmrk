package config

type config struct {
	JsonFile string
	Address  string
}

var Settings = config{
	JsonFile: "bookmarks.json",
	Address:  ":8080",
}
