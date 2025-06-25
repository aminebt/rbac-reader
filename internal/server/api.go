package server

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	rbacv1alpha1 "github.com/aminebt/rbac-operator/api/v1alpha1"
	"github.com/go-chi/chi"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Api struct {
	logger *slog.Logger
	Scheme *runtime.Scheme
	Client client.Client
}

func NewApi(logger *slog.Logger, client client.Client, scheme *runtime.Scheme) (*Api, error) {
	api := &Api{
		logger: logger,
		Scheme: scheme,
		Client: client, // this is a controller-runtime client
	}

	return api, nil
}

func NewRouter(api *Api) http.Handler {
	r := chi.NewRouter()

	r.Get("/groups", api.handlerGetGroups)

	return r
}

func (api *Api) handlerGetGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	groups := &rbacv1alpha1.GridOSGroupList{}
	if err := api.Client.List(ctx, groups); err != nil {
		api.logger.Error("failed to list groups", "error: ", err)
		respondWithError(w, 500, err.Error())
	}

	respondWithJSON(w, 200, groups)
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
