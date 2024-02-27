package config

type Application struct {
	Name string
	Run  string
}

var ApplicationConfig = new(Application)
