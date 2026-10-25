package nas

import (
	"net/http"
	"openpak/wfc-common/service"
	"os"
)

func registerRemotes() {
	key := service.Key()
	for _, v := range []struct {
		mux *http.ServeMux
		env string
	}{{sakeMux, "WFC_SAKE_URL"}, {raceMux, "WFC_RACE_URL"}, {gamestatsMux, "WFC_GAMESTATS_URL"}, {dlsMux, "WFC_DLS_URL"}} {
		proxy, e := service.Proxy(os.Getenv(v.env), key)
		if e != nil {
			panic(e)
		}
		v.mux.Handle("/", proxy)
	}
	api, e := service.Proxy(os.Getenv("WFC_API_URL"), key)
	if e != nil {
		panic(e)
	}
	http.DefaultServeMux.Handle("/api/", api)
}
