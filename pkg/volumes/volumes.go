package volumes

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"mysql-operator/controllers"

	//"mysql-operator/controllers"
	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	//"sigs.k8s.io/controller-runtime/pkg/client"
	//ctrl "sigs.k8s.io/controller-runtime"

)
type Volumes struct {
	client.Client
	Scheme *runtime.Scheme
}


func (v *Volumes) CreateVolumes(ns string, sc string, name string, size string,ctx context.Context) error {
	//
	foundConfigMap := &apiv1.ConfigMap{}
	err := v.Get(ctx, types.NamespacedName{Name: "kube-root-ca.crt", Namespace: "default"}, foundConfigMap)
	if err != nil {
		// If a configMap name is provided, then it must exist
		// You will likely want to create an Event for the user to understand why their reconcile is failing.
		return err
	}
	fmt.Println(foundConfigMap)
	return nil
}


