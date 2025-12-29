package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zemld/Scently/perfume-hub/api/handlers"
	"github.com/zemld/Scently/perfume-hub/api/middleware"
	"github.com/zemld/Scently/perfume-hub/internal/db/core"
)

func main() {
	time.Sleep(5 * time.Second)
	core.Initiate()
	defer core.Close()
	r := http.NewServeMux()

	r.Handle("/v1/perfumes/get", middleware.Auth(http.HandlerFunc(handlers.Select)))
	r.Handle("/v1/perfumes/update", middleware.Auth(http.HandlerFunc(handlers.Update)))

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PERFUME_HUB_PORT")), r)
}
