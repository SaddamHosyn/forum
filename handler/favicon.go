package handler

import "net/http"

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/images/favicon.ico")
}
