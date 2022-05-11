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

package controllers

import (
	"context"
	examplev1beta1 "github.com/Youngpig1998/petClinic-operator/api/v1beta1"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/bootstrap"
	"github.com/Youngpig1998/petClinic-operator/internal/operator"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var servicesName [5]string = [5]string{"customers", "vets", "visits", "web", "gateway"}

var (
	controllerManagerName = "petclinic-operator-controller-manager"
)

// PetClinicReconciler reconciles a PetClinic object
type PetClinicReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	Config *rest.Config
}

//+kubebuilder:rbac:groups=example.njtech.edu.cn,resources=petclinics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.njtech.edu.cn,resources=petclinics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.njtech.edu.cn,resources=petclinics/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PetClinic object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PetClinicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("petclinic", req.NamespacedName)
	//log := log.FromContext(ctx)

	log.Info("1. start reconcile logic")
	// Instantialize the data structure
	instance := &examplev1beta1.PetClinic{}

	//First,query the webhook instance
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		// If there is no instance, an empty result is returned, so that the Reconcile method will not be called immediately
		if errors.IsNotFound(err) {
			log.Info("Instance not found, maybe removed")
			return reconcile.Result{}, nil
		}
		log.Error(err, "query action happens error")
		// Return error message
		return ctrl.Result{}, err
	}

	//Set the bootstrapClient's owner value as the webhook,so the resources we create then will be set reference to the webhook
	//when the webhook cr is deleted,the resources(such as deployment.configmap,issuer...) we create will be deleted too
	bootstrapClient, err := bootstrap.NewClient(r.Config, r.Scheme, controllerManagerName, instance)
	if err != nil {
		log.Error(err, "failed to initialise bootstrap client")
		return ctrl.Result{}, err
	}

	//First we create mongodb services,include geo,user,profile,recommendation,rate,reservation

	if instance.Spec.MysqlActive == true {
		statefulSet := operator.StatefulSet("mysql", instance)
		err = bootstrapClient.CreateResource("mysql", statefulSet)
		if err != nil {
			log.Error(err, "failed to create operator's mysql StatefulSet", "Name", "mysql")
			return ctrl.Result{}, err
		}
		service := operator.Service("mysql", 3306, 3306, 0)
		err = bootstrapClient.CreateResource("mysql", service)
		if err != nil {
			log.Error(err, "failed to create operator's mysql Service", "Name", "mysql")
			return ctrl.Result{}, err
		}
	}

	time.Sleep(time.Duration(10) * time.Second)

	//We create petclinic services
	for i := 0; i < 5; i++ {
		deploy := operator.Deployment(servicesName[i], instance)
		deployName := servicesName[i]
		err = bootstrapClient.CreateResource(deployName, deploy)
		if err != nil {
			log.Error(err, "failed to create operator's deployment", "Name", deployName)
			return ctrl.Result{}, err
		}

		var nodePort int32 = 0
		if servicesName[i] == "gateway" {
			nodePort = 31080
		}

		service := operator.Service(deployName, 8080, 8080, nodePort)
		err = bootstrapClient.CreateResource(deployName, service)
		if err != nil {
			log.Error(err, "failed to create operator's  Service", "Name", deployName)
			return ctrl.Result{}, err
		}

	}

	////Then we create consul service
	//deploymentForConsul := operator.DeploymentForConsul(instance)
	//err = bootstrapClient.CreateResource("consul", deploymentForConsul)
	//if err != nil {
	//	log.Error(err, "failed to create operator's consul Deployment", "Name", "consul")
	//	return ctrl.Result{}, err
	//}
	//
	////Then we create jaeger service
	//deploymentForJaeger := operator.DeploymentForJaeger(instance)
	//err = bootstrapClient.CreateResource("jaeger", deploymentForJaeger)
	//if err != nil {
	//	log.Error(err, "failed to create operator's jaeger Deployment", "Name", "jaeger")
	//	return ctrl.Result{}, err
	//}
	//
	////Then we create logic services,include search geo rate profile recommendation user
	//
	//for i := 0; i < 8; i++ {
	//
	//	var port int32 = 0
	//	if servicesName[i] == "rate" {
	//		port = 8084
	//	} else if servicesName[i] == "profile" {
	//		port = 8081
	//	} else if servicesName[i] == "reservation" {
	//		port = 8087
	//	} else if servicesName[i] == "user" {
	//		port = 8086
	//	} else if servicesName[i] == "geo" {
	//		port = 8083
	//	} else if servicesName[i] == "frontend" {
	//		port = 5000
	//	} else if servicesName[i] == "search" {
	//		port = 8082
	//	} else {
	//		port = 8085
	//	}
	//
	//	deploymentForLogic := operator.DeploymentForLogic(servicesName[i], port, instance)
	//	err = bootstrapClient.CreateResource(servicesName[i], deploymentForLogic)
	//	if err != nil {
	//		log.Error(err, "failed to create operator's logic Deployment", "Name", servicesName[i])
	//		return ctrl.Result{}, err
	//	}
	//
	//}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PetClinicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1beta1.PetClinic{}).
		Complete(r)
}
