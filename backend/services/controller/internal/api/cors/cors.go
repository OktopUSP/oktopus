package cors

import (
	"github.com/rs/cors"
	"net/http"
	"os"
	"strings"
	"fmt"
)

func GetCorsConfig() cors.Cors {
	allowedOrigins := getCorsEnvConfig()
	fmt.Println(allowedOrigins)
	return *cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application
		},
	})
}

func getCorsEnvConfig() []string {
	val, _ := os.LookupEnv("REST_API_CORS")
	if val == "" {
		return []string{"*"}
	}
	return strings.Split(val, ",")
}