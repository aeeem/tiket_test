package main

import (
	internal_http "tiket_test/provider/internal/http"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`./config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

//	@title			Tech Test Api Docs
//	@version		1.0
//	@description	This is a sample tiket_test/server Petstore tiket_test/server.
//	@termsOfService	http://swagger.io/terms/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Just insert after you hit login
func main() {

	internal_http.HttpRun(viper.GetString("server.address")+":"+viper.GetString("server.port"), 2)

}
