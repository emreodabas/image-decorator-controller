package controller

import (
	"github.com/emreodabas/image-decorator-controller/pkg/containerimage"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type Reconciler struct {
	Client            client.Client
	RequeueDuration   time.Duration
	IgnoredNamespaces []string
	BackupRegistry    *containerimage.ContainerRepository
}

// for enabling leader election
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch

func (r *Reconciler) isIgnoredNamespace(name string) bool {
	for _, ns := range r.IgnoredNamespaces {
		if ns == name {
			return true
		}
	}
	return false
}

func (r *Reconciler) successMessage() (reconcile.Result, error) {
	return reconcile.Result{Requeue: false}, nil
}

func (r *Reconciler) requeueMessage(err error) (reconcile.Result, error) {
	return reconcile.Result{Requeue: true, RequeueAfter: r.RequeueDuration}, err
}

func (r *Reconciler) errorWithoutRequeue(err error) (reconcile.Result, error) {
	return reconcile.Result{Requeue: false}, err
}
