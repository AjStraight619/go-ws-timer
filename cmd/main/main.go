package main

import (
	"net/http"

	"github.com/AjStraight619/go-ws-timer/internal/routes"
)

func main() {
	r := routes.SetUpRouter()
	routes.SetUpRoutes(r)
	http.ListenAndServe(":8080", r)
}
