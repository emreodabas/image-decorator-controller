package controller

import (
	"testing"
)

func TestReconcileDeployment(t *testing.T) {

	request := getIgnoredNSReconcileRequest()
	context := getContext()

	reconcile, err := getDeployment().Reconcile(context, request)
	if err != nil {
		t.Errorf("Error is not expected for ignored namespaces %v", err)
	}
	if reconcile.Requeue == true || reconcile.RequeueAfter > 0 {
		t.Errorf("Requeue is not expected for given ")
	}

	request = getRandomNameFromDefaultNSReconcileRequest()
	reconcile, err = getDeployment().Reconcile(context, request)
	if err == nil {
		t.Errorf("Error is expected but not found")
	}
	if reconcile.Requeue != true && reconcile.RequeueAfter != getReconciler().RequeueDuration {
		t.Errorf("Requeue is expected with given duration")
	}
}
