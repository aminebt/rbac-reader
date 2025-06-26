package api

import (
	"context"
	"fmt"

	rbac "github.com/aminebt/rbac-operator/rbac"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Datasource interface {
	GetGroups(ctx context.Context) ([]rbac.Group, error)
}

func NewDataSource(backend string) (Datasource, manager.Manager, error) {
	switch backend {
	case "kubernetes":
		return NewKubernetesSource()
	// case "postgres":
	// 	return NewPostgresSource()
	default:
		return nil, nil, fmt.Errorf("unsupported backend: %s", backend)
	}
}
