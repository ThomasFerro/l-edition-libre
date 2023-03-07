package main

import (
	"github.com/ThomasFerro/l-edition-libre/api"
	"github.com/ThomasFerro/l-edition-libre/configuration"
)

func main() {
	api.Start(configuration.GetConfiguration(configuration.MONGO_DATABASE_NAME))
}
