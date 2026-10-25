package nas

import (
	"net/http"
	"openpak/wfc-common/dwc"
)

func parseAuthRequest(r *http.Request) (map[string][]byte, error) { return dwc.ParseAuthRequest(r) }
func writeAuthResponse(w http.ResponseWriter, reply map[string]string) {
	dwc.WriteAuthResponse(w, reply)
}
func getDateTime() string { return dwc.DateTime() }
