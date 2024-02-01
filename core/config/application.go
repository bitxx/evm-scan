package config

type Application struct {
	Name    string
	NodeUrl string
}

var ApplicationConfig = new(Application)
