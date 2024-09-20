package handlers

import "net/http"

func Tools(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/tools.html")
}
func IpLookup(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/ip_lookup.html")
}
