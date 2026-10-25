package nas

import (
	"net/http"
	"openpak/wfc-common/service"
	"os"
)

func registerRemotes() {
	key := service.Key()
	content, err := service.Proxy(os.Getenv("WFC_CONTENT_URL"), key)
	if err != nil {
		panic(err)
	}
	runtime, err := service.Proxy(os.Getenv("WFC_RUNTIME_URL"), key)
	if err != nil {
		panic(err)
	}
	sakeMux.Handle("/", content)
	raceMux.Handle("/", content)
	gamestatsMux.Handle("/", runtime)
	http.DefaultServeMux.Handle("/api/", runtime)
}
