package controller

import "testing"

func TestReconcileDaemonSet(t *testing.T) {

	request := getIgnoredNSReconcileRequest()
	context := getContext()

	reconcile, err := getDaemonSet().Reconcile(context, request)
	if err != nil {
		t.Errorf("Error is not expected for ignored namespaces %v", err)
	}
	if reconcile.Requeue == true || reconcile.RequeueAfter > 0 {
		t.Errorf("Requeue is not expected for given ")
	}

	request = getRandomNameFromDefaultNSReconcileRequest()
	reconcile, err = getDaemonSet().Reconcile(context, request)
	if err == nil {
		t.Errorf("Error is expected but not found")
	}
	if reconcile.Requeue == true && reconcile.RequeueAfter > 0 {
		t.Errorf("Requeue is expected with given duration")
	}
}
