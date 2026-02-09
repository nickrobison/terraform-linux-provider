package firewall

import (
	"fmt"
	"net/http"
)

func HandleHello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from firewall!")
	})
}
