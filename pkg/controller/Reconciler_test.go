package controller

import (
	"context"
	"fmt"
	"github.com/emreodabas/image-decorator-controller/pkg/containerimage"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"
	"testing"
)

func TestIsIgnoredNamespaceWithNotIgnored(t *testing.T) {
	isIgnored := getDeployment().Deployment.isIgnoredNamespace("notignored" + strconv.Itoa(rand.Int()))
	if isIgnored {
		t.Errorf("Expected it is not ignored")
	}
}

func TestIsIgnoredNamespaceWithIgnored(t *testing.T) {
	ignoredNS := getRandomIgnoderNS()
	isIgnored := getDeployment().Deployment.isIgnoredNamespace(ignoredNS)
	if !isIgnored {
		t.Errorf("Expected %v is ignored", ignoredNS)
	}
}

func TestSuccessMessage(t *testing.T) {
	message, err := getDeployment().Deployment.successMessage()
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
	if message.Requeue == true || message.RequeueAfter > 0 {
		t.Errorf("Expected to message is not requeue with given duration")
	}
}

func TestRequeueMessage(t *testing.T) {
	message, err := getDeployment().Deployment.requeueMessage(fmt.Errorf("test error"))
	if err == nil {
		t.Errorf("Error is expected but not found ")
	}
	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("Error payload is not matched with triggered one %v", err.Error())
	}
	if message.Requeue == false || message.RequeueAfter == 0 {
		t.Errorf("Expected to message is requeued with given duration %v", message.RequeueAfter)
	}

	if message.RequeueAfter != getReconciler().RequeueDuration {
		t.Errorf("Expected %v duration but found %v ", getReconciler().RequeueDuration, message.RequeueAfter)
	}
}

func getDeployment() *ReconcileDeployment {
	return &ReconcileDeployment{
		Deployment: getReconciler(),
	}
}

func getDaemonSet() *ReconcileDaemonSet {
	return &ReconcileDaemonSet{
		DaemonSet: getReconciler(),
	}
}

func getReconciler() Reconciler {
	return Reconciler{
		Client:            getManager().GetClient(),
		RequeueDuration:   5000,
		IgnoredNamespaces: []string{"kube-system", "kube-public", "ignore-ns"},
		BackupRegistry:    getBackupRegistry(),
	}

}
func getManager() manager.Manager {
	mgr, _ := ctrl.NewManager(config.GetConfigOrDie(), ctrl.Options{
		Scheme:         runtime.NewScheme(),
		Port:           9443,
		LeaderElection: false,
	})

	return mgr
}

func getBackupRegistry() *containerimage.ContainerRepository {
	return &containerimage.ContainerRepository{
		RepositoryPath: "kubermatico/",
		Username:       "kubermatico",
		Password:       "$pus@L&3G?!ewTK",
		AccessToken:    "a5c4c4a4-27f8-40af-9662-a0391bee1d6d",
	}
}

func getRandomIgnoderNS() string {
	return getReconciler().IgnoredNamespaces[rand.Int()%len(getReconciler().IgnoredNamespaces)]
}

func getIgnoredNSReconcileRequest() reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: getRandomIgnoderNS(),
			Name:      "name",
		},
	}
}

func getRandomNameFromDefaultNSReconcileRequest() reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "default",
			Name:      "app" + strconv.Itoa(rand.Int()),
		},
	}
}

func getContext() context.Context {
	c, err := context.WithTimeout(context.Background(), 10000)
	if err != nil {
		return nil
	}
	return c
}
