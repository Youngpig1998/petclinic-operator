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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PetClinicSpec defines the desired state of PetClinic
type PetClinicSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//
	Replicas int32 `json:"replicas"`

	MysqlActive bool `json:"mysqlActive"`

	ScaleCrossCloud bool `json:"scaleCrossCloud"`

	Mysql map[string]string `json:"mysql,omitempty" protobuf:"bytes,8,rep,name=mysql"`

	Hpa map[string]string `json:"hpa,omitempty" protobuf:"bytes,8,rep,name=hpa"`
}

// PetClinicStatus defines the observed state of PetClinic
type PetClinicStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PetClinic is the Schema for the petclinics API
type PetClinic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PetClinicSpec   `json:"spec,omitempty"`
	Status PetClinicStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PetClinicList contains a list of PetClinic
type PetClinicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PetClinic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PetClinic{}, &PetClinicList{})
}
