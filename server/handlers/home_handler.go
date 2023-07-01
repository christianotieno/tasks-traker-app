package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		return
	}
}
