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
package applications

import (
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources"
	oamv1beta1 "github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Application is a wrapper around the oamv1beta1.Application object that meets the
// Reconcileable interface
type Application struct {
	*oamv1beta1.Application
}

// From returns a new Reconcileable Service from a oamv1beta1.Application
func From(application *oamv1beta1.Application) *Application {
	return &Application{Application: application}
}

// ShouldUpdate returns whether the resource should be updated in Kubernetes and
// the resource to update with
func (a Application) ShouldUpdate(current client.Object) (bool, client.Object) {

	currentApplication := current.DeepCopyObject().(*oamv1beta1.Application)
	newApplication := currentApplication.DeepCopy()
	resources.MergeMetadata(newApplication, a)
	newApplication.Spec = a.Spec
	return !equality.Semantic.DeepEqual(newApplication, current), newApplication
}

// GetResource retrieves the resource instance
func (a Application) GetResource() client.Object {
	return a.Application
}

// ResourceKind retrieves the string kind of the resource
func (a Application) ResourceKind() string {
	return "Application"
}

// ResourceIsNil returns whether or not the resource is nil
func (a Application) ResourceIsNil() bool {
	return a.Application == nil
}

// NewResourceInstance returns a new instance of the sme resource type
func (a Application) NewResourceInstance() client.Object {
	return &oamv1beta1.Application{}
}
