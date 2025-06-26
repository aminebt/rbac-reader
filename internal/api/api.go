package api

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/aminebt/rbac-reader/internal/server"
	"github.com/go-chi/chi"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Backend string

// const (
// 	Kubernetes Backend = "kubernetes"
// 	Postgresql Backend = "postgresql"
// )

type Api struct {
	Logger     *slog.Logger
	Datasource Datasource
	Backend    string
	Server     server.Server
	Manager    manager.Manager
}

func NewRouter(api *Api) http.Handler {
	r := chi.NewRouter()

	r.Get("/groups", api.handlerGetGroups)

	return r
}

func (api *Api) handlerGetGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	groups, err := api.Datasource.GetGroups(ctx)
	api.Logger.Info("number of groups: ", len(groups))
	if err != nil {
		api.Logger.Error("failed to list groups", "error: ", err)
		respondWithError(w, 500, err.Error())
	}
	respondWithJSON(w, 200, groups)
}

func (api *Api) Start() {
	api.Server.SetHandler(NewRouter(api))
	if api.Backend == "kubernetes" {
		if api.Manager == nil {
			api.Logger.Error("Data source with k8s backend was not initialized correctly. Exiting")
			os.Exit(1)
		}
		if err := api.Server.SetupWithManager(api.Manager); err != nil {
			api.Logger.Error("failed to setup server with manager", "error", err)
			os.Exit(1)
		}

		// start the manager
		api.Logger.Info("starting manager")
		err := api.Manager.Start(ctrl.SetupSignalHandler())
		if err != nil {
			api.Logger.Error("unable to start manager", "error", err)
			os.Exit(1)
		}
	} else {
		err := api.Server.Start(ctrl.SetupSignalHandler())
		if err != nil {
			api.Logger.Error("unable to start server", "error", err)
			os.Exit(1)
		}
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("error marshalling payload to json", err)
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorBody struct {
		Err string `json:"error"`
	}
	respondWithJSON(w, code, errorBody{
		Err: msg,
	})
}
