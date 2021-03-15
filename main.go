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
	"github.com/emreodabas/image-decorator-controller/pkg/containerimage"
	"github.com/emreodabas/image-decorator-controller/pkg/controller"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
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
	mgr, err := ctrl.NewManager(getKubeConfig(), ctrl.Options{
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
		Owns(&v1.Pod{}).
		Complete(&controller.ReconcileDeployment{
			Deployment: getReconciler(mgr),
		})
	if err != nil {
		logger.Error(err, "unable to create controller")
		os.Exit(1)
	}
	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Owns(&v1.Pod{}).
		Complete(&controller.ReconcileDaemonSet{
			DaemonSet: getReconciler(mgr),
		})
	// kubebuilder:scaffold:builder
	logger.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logger.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func setVariables() {
	if os.Getenv("ENV") == "dev" || (len(os.Args) > 1 && os.Args[1] == "dev") {
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

func getKubeConfig() *rest.Config {
	var kubeConfig *rest.Config
	configPath := viper.GetString("KUBE_CONFIG_PATH")
	if configPath != "" {
		if kubeConfigPath := configPath; kubeConfigPath != "" {
			configPath = kubeConfigPath
		}
		kubeConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			logger.Error(err, "Kube config could not build in dev")
		}

		return kubeConfig
	}

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Error(err, "Cluster config error")
	}

	return kubeConfig
}

func getReconciler(mgr ctrl.Manager) controller.Reconciler {
	return controller.Reconciler{
		Client:            mgr.GetClient(),
		RequeueDuration:   time.Duration(viper.GetInt64("REQUEUE_DURATION")),
		IgnoredNamespaces: strings.Split(viper.GetString("IGNORED_NS"), ","),
		BackupRegistry: &containerimage.ContainerRepository{
			RepositoryPath: viper.GetString("BACKUP_REGISTRY_ADDRESS"),
			Username:       viper.GetString("USERNAME"),
			Password:       viper.GetString("PASSWORD"),
			AccessToken:    viper.GetString("ACCESS_TOKEN"),
		},
	}
}
