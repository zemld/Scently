package main

import (
	"net/http"

	"github.com/zemld/Scently/perfume-hub/api/handlers"
	"github.com/zemld/Scently/perfume-hub/api/middleware"
	"github.com/zemld/Scently/perfume-hub/internal/db/core"
)

func main() {
	core.Initiate()
	defer core.Close()

	r := http.NewServeMux()

	r.Handle("/v1/perfumes/get", middleware.Auth(http.HandlerFunc(handlers.Select)))
	r.Handle("/v1/perfumes/update", middleware.Auth(http.HandlerFunc(handlers.Update)))

	http.ListenAndServe(":8000", r)
}
