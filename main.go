package main

import (
	"fmt"
	"log/slog"
	"os"

	rbacv1alpha1 "github.com/aminebt/rbac-operator/api/v1alpha1"
	"github.com/aminebt/rbac-reader/internal/server"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	//gr := &rbacv1alpha1.GridOSGroup{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	restConfig := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()

	if err := rbacv1alpha1.AddToScheme(scheme); err != nil {
		fmt.Println("failed to add rbac types to scheme", "error", err)
		os.Exit(1)
	}
	logger.Info("successfully added rbac types to scheme")

	ctrlOpts := ctrl.Options{
		Scheme:         scheme,
		LeaderElection: false,
	}

	mgr, err := ctrl.NewManager(restConfig, ctrlOpts)

	if err != nil {
		logger.Error("unable to create manager", "error", err)
		os.Exit(1)
	}

	api, err := server.NewApi(logger, mgr.GetClient(), mgr.GetScheme())
	if err != nil {
		logger.Error("failed to create api", "error", err)
		os.Exit(1)
	}

	svr, err := server.NewServer(api, logger)

	if err := svr.SetupWithManager(mgr); err != nil {
		logger.Error("failed to setup server with manager", "error", err)
		os.Exit(1)
	}

	// start the manager
	logger.Info("starting manager")
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		logger.Error("unable to start manager", "error", err)
		os.Exit(1)
	}

}
