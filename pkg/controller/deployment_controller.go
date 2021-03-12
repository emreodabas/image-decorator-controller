package controller

import (
	"context"
	"fmt"
	"github.com/emreodabas/image-decorator-controller/pkg/containerimage"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDeployment struct {
	Client client.Client
}

// !! Kubebuilder will read this lines and generate related resources !!
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileDeployment{}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

	log := log.FromContext(ctx)

	// Fetch the Deployment from the cache
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
	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		log.Info("Reconciling deployment", "containerimage name", container.Name)
		//TODO take prefix as variable
		imagePath, err := containerimage.CloneImage(container.Image, "kubermatico/")
		if err != nil {
			return reconcile.Result{}, nil
		}
		containers[i].Image = imagePath
	}
	//re assign updated containers to deployment
	deployment.Spec.Template.Spec.Containers = containers
	err = r.Client.Update(ctx, deployment)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not write Deployment: %+v", err)
	}

	return reconcile.Result{}, nil

}
