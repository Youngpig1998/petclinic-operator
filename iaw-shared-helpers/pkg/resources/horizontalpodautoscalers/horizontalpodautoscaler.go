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
package horizontalpodautoscalers

import (
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HorizontalPodAutoscaler is a wrapper around the autoscalingv2beta2.HorizontalPodAutoscaler object that meets the
// Reconcileable interface
type HorizontalPodAutoscaler struct {
	*autoscalingv2beta2.HorizontalPodAutoscaler
}

// From returns a new Reconcileable HorizontalPodAutoscaler from a autoscalingv2beta2.HorizontalPodAutoscaler
func From(horizontalPodAutoscaler *autoscalingv2beta2.HorizontalPodAutoscaler) *HorizontalPodAutoscaler {
	return &HorizontalPodAutoscaler{HorizontalPodAutoscaler: horizontalPodAutoscaler}
}

// ShouldUpdate returns whether the resource should be updated in Kubernetes and
// the resource to update with
func (h HorizontalPodAutoscaler) ShouldUpdate(current client.Object) (bool, client.Object) {

	currentHorizontalPodAutoscaler := current.DeepCopyObject().(*autoscalingv2beta2.HorizontalPodAutoscaler)
	newHorizontalPodAutoscaler := currentHorizontalPodAutoscaler.DeepCopy()
	resources.MergeMetadata(newHorizontalPodAutoscaler, h)
	newHorizontalPodAutoscaler.Spec = h.Spec
	return !equality.Semantic.DeepEqual(newHorizontalPodAutoscaler, current), newHorizontalPodAutoscaler
}

// GetResource retrieves the resource instance
func (h HorizontalPodAutoscaler) GetResource() client.Object {
	return h.HorizontalPodAutoscaler
}

// ResourceKind retrieves the string kind of the resource
func (h HorizontalPodAutoscaler) ResourceKind() string {
	return "HorizontalPodAutoscaler"
}

// ResourceIsNil returns whether or not the resource is nil
func (h HorizontalPodAutoscaler) ResourceIsNil() bool {
	return h.HorizontalPodAutoscaler == nil
}

// NewResourceInstance returns a new instance of the same resource type
func (h HorizontalPodAutoscaler) NewResourceInstance() client.Object {
	return &autoscalingv2beta2.HorizontalPodAutoscaler{}
}
