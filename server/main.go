package main

import (
	"tiket_test/server/internal/http"

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

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Just insert after you hit login

func main() {
	http.HttpRun(viper.GetString("tiket_test/server.address") + ":" + viper.GetString("tiket_test/server.port"))
}
