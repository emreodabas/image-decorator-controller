/*
Copyright 2021 emreo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/emreodabas/image-decorator-controller/pkg/controller"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"strings"
	"time"
)

// +kubebuilder:scaffold:imports

var (
	logger = log.Log.WithName("entrypoint")
	scheme = runtime.NewScheme()
)

func init() {
	log.SetLogger(zap.New())
	clientgoscheme.AddToScheme(scheme)
	setVariables()
	// +kubebuilder:scaffold:scheme

}

func main() {
	// Setup a Manager
	logger.Info("setting up manager")
	mgr, err := ctrl.NewManager(config.GetConfigOrDie(), ctrl.Options{
		Scheme:         scheme,
		Port:           9443,
		LeaderElection: false,
	})
	if err != nil {
		logger.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile Deployment
	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Owns(&corev1.Pod{}).
		Complete(&controller.ReconcileDeployment{
			Client:            mgr.GetClient(),
			RequeueDuration:   time.Duration(viper.GetInt64("REQUEUE_DURATION")),
			IgnoredNamespaces: strings.Split(viper.GetString("IGNORED_NS"), ","),
			BackupRegistry:    viper.GetString("BACKUP_REGISTRY_ADDRESS"),
		})
	if err != nil {
		logger.Error(err, "unable to create controller")
		os.Exit(1)
	}
	// kubebuilder:scaffold:builder
	logger.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func setVariables() {

	if os.Getenv("ENV") == "dev" || os.Args[1] == "dev" {
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			logger.Error(err, "Error reading config file, %s")
		}
	} else {
		viper.AutomaticEnv()
	}
}
