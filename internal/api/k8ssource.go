package api

import (
	"context"
	"log/slog"
	"os"

	rbacv1alpha1 "github.com/aminebt/rbac-operator/api/v1alpha1"
	rbac "github.com/aminebt/rbac-operator/rbac"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type KubernetesSource struct {
	logger  *slog.Logger
	Scheme  *runtime.Scheme
	Manager manager.Manager
	Client  client.Client
}

func (ks *KubernetesSource) GetGroups(ctx context.Context) ([]rbac.Group, error) {
	groups := &rbacv1alpha1.GridOSGroupList{}

	if err := ks.Client.List(ctx, groups); err != nil {
		ks.logger.Error("failed to list groups resources", "error: ", err)
		return nil, err
	}
	ks.logger.Info("number of groups: ", len(groups.Items))
	ks.logger.Info("number of plain object groups: ", len(groups.ToPlainObject()))
	return groups.ToPlainObject(), nil

}

func NewKubernetesSource() (*KubernetesSource, manager.Manager, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	restConfig := ctrl.GetConfigOrDie()
	scheme := runtime.NewScheme()

	if err := rbacv1alpha1.AddToScheme(scheme); err != nil {
		logger.Error("failed to add rbac types to scheme", "error", err)
		return nil, nil, err
	}
	logger.Info("successfully added rbac types to scheme")

	ctrlOpts := ctrl.Options{
		Scheme:         scheme,
		LeaderElection: false,
	}

	mgr, err := ctrl.NewManager(restConfig, ctrlOpts)

	if err != nil {
		logger.Error("unable to create manager", "error", err)
		return nil, nil, err
	}

	ks := &KubernetesSource{
		logger:  logger,
		Scheme:  mgr.GetScheme(),
		Client:  mgr.GetClient(), // this is a controller-runtime client
		Manager: mgr,
	}

	return ks, mgr, nil

}

func (ks *KubernetesSource) Start() {

}
