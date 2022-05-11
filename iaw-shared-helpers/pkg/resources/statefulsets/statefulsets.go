// ------------------------------------------------------ {COPYRIGHT-TOP} ---
// IBM Confidential
// OCO Source Materials
// 5900-AEO
//
// Copyright IBM Corp. 2021
//
// The source code for this program is not published or otherwise
// divested of its trade secrets, irrespective of what has been
// deposited with the U.S. Copyright Office.
// ------------------------------------------------------ {COPYRIGHT-END} ---

package statefulsets

import (
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Deployment is a wrapper around the appsv1.Deployment object that meets the
// Reconcileable interface
type StatefulSet struct {
	*appsv1.StatefulSet
}

// From returns a new Reconcileable Deployment from a appsv1.Deployment
func From(statefulset *appsv1.StatefulSet) *StatefulSet {
	return &StatefulSet{StatefulSet: statefulset}
}

// ShouldUpdate returns whether the resource should be updated in Kubernetes and
// the resource to update with
func (s StatefulSet) ShouldUpdate(current client.Object) (bool, client.Object) {
	newStatefulSet := current.DeepCopyObject().(*appsv1.StatefulSet)
	resources.MergeMetadata(newStatefulSet, s)
	resources.MergeMetadata(&newStatefulSet.Spec.Template, &s.Spec.Template)
	mergedTemplate := newStatefulSet.Spec.Template
	newStatefulSet.Spec = s.Spec
	newStatefulSet.Spec.Template.ObjectMeta = mergedTemplate.ObjectMeta
	return !equality.Semantic.DeepEqual(newStatefulSet, current), newStatefulSet
}

// GetResource retrieves the resource instance
func (s StatefulSet) GetResource() client.Object {
	return s.StatefulSet
}

// ResourceKind retrieves the string kind of the resource
func (s StatefulSet) ResourceKind() string {
	return "StatefulSet"
}

// ResourceIsNil returns whether or not the resource is nil
func (s StatefulSet) ResourceIsNil() bool {
	return s.StatefulSet == nil
}

// NewResourceInstance returns a new instance of the same resource type
func (s StatefulSet) NewResourceInstance() client.Object {
	return &appsv1.StatefulSet{}
}
