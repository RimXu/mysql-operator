/*
Copyright 2022.

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
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	databasev1 "mysql-operator/api/v1"
	"mysql-operator/controllers"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	//+kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
	//setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(databasev1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	// 初始化log类型
	// logrus.SetFormatter(&logrus.TextFormatter{})
	// 日志格式为保留彩色，完整时间戳，时间格式为微秒
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:03:04.0000",
	})
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	//flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	//flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	//flag.BoolVar(&enableLeaderElection, "leader-elect", false,
	//	"Enable leader election for controller manager. "+
	//		"Enabling this will ensure there is only one active controller manager.")
	//opts := zap.Options{
	//	Development: true,
	//}
	//opts.BindFlags(flag.CommandLine)
	//flag.Parse()

	logrus.Infof("metrics-bind-address:%s", "8080")
	logrus.Infof("health-probe-address:%s", "8081")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "5e3c7de0.operator.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		logrus.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.MysqlReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		logrus.Error(err, "unable to create controller", "controller", "Mysql")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		logrus.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		logrus.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	logrus.Info("start manager successfully")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logrus.Error(err, "problem running manager")
		os.Exit(1)
	}
}
