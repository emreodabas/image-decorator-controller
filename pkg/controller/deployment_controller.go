package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDeployment struct {
	Client client.Client
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileDeployment{}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

	log := log.FromContext(ctx)

	// Fetch the ReplicaSet from the cache
	deployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, request.NamespacedName, deployment)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find deployment")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch deployment: %+v", err)
	}

	// Print the Deployment
	log.Info("Reconciling deployment", "container name", deployment.Spec.Template.Spec.Containers[0].Name)

	// Set the label if it is missing
	if deployment.Labels == nil {
		deployment.Labels = map[string]string{}
	}
	if deployment.Labels["hello"] == "world" {
		return reconcile.Result{}, nil
	}

	// Update the ReplicaSet
	deployment.Labels["hello"] = "world"
	err = r.Client.Update(context.TODO(), deployment)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not write ReplicaSet: %+v", err)
	}

	return reconcile.Result{}, nil

}
