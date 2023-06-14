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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MysqlSpec defines the desired state of Mysql
type MysqlSpec struct {

	// 定义复制配置,默认是true，Operator将生成主从MySQL deployment
	// 定义字段类型校验
	// Run "make" and "make install" to regenerate code after modifying this file
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format:= bool
	Replication bool `json:"replication,true"`
	// 定义MySQL规格套餐,默认small,只支持small,medium,large三个属性
	// 可以在mysql-operator/pkg/constants中进行修改对应Instance的配置
	// +kubebuilder:validation:Enum:= small;medium;large
	Instance string `json:"instance,small"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format:= string
	StorageClass string `json:"storageclass"`
	// +kubebuilder:validation:Required
	Databases []map[string]string `json:"databases"`
	// +kubebuilder:validation:Required
	Version string `json:"version,mysql:5.7.35"`
	// +optional
	Phase string `json:"phase,omitempty"`
}

// MysqlStatus defines the observed state of Mysql
// 由于公司使用k8s版本较老,启用和更新status状态存在问题,所以Operator全部使用spec.phase作为状态列
type MysqlStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replication",type= string,JSONPath=`.spec.replication`
// +kubebuilder:printcolumn:name="Instance",type=string,JSONPath=`.spec.instance`
// +kubebuilder:printcolumn:name="StorageClass",type=string,JSONPath=`.spec.storageclass`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.spec.phase`
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

type Mysql struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MysqlSpec   `json:"spec,omitempty"`
	Status MysqlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MysqlList contains a list of Mysql
type MysqlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mysql `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mysql{}, &MysqlList{})
}
