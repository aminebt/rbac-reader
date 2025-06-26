package main

import (
	"log/slog"
	"os"

	rbacapi "github.com/aminebt/rbac-reader/internal/api"
	"github.com/aminebt/rbac-reader/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	backend := "kubernetes"

	// initialize datasource - mgr is non nil in case of kubernetes backend
	ds, mgr, err := rbacapi.NewDataSource(backend)
	if err != nil {
		logger.Error("Unable to initialize datasource of backend type %s", backend)
		os.Exit(1)
	}

	svr, err := server.NewServer(logger)
	if err != nil {
		logger.Error("failed to create server", "error", err)
		os.Exit(1)
	}
	// create the API
	api := &rbacapi.Api{
		Logger:     logger,
		Datasource: ds,
		Backend:    backend,
		Server:     *svr,
		Manager:    mgr,
	}

	api.Start()

}
