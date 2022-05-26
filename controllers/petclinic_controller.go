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
	"fmt"
	examplev1beta1 "github.com/Youngpig1998/petClinic-operator/api/v1beta1"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/bootstrap"
	"github.com/Youngpig1998/petClinic-operator/internal/operator"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var servicesName [5]string = [5]string{"customers", "vets", "visits", "web", "gateway"}

var isWatchPod bool = false

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
//+kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.oam.dev,resources=applications,verbs=get;list;watch;create;update;patch;delete

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

	// just a test
	mysqls := instance.Spec.Mysql
	if mysqls == nil {
		log.Info("Mysql iudshfbidsbfmjsdabnfuiwebfjhksxb ds nfdjksfnb kjdsa bfuiawbf jkdsb njksdabnf kjw ndnsjkf ")
	} else {
		log.Info(mysqls["sdsd"])
	}

	//First we create mysql statefulset

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

	//We create petclinic deployments & services
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

	//We create petclinic hpas
	for i := 0; i < 5; i++ {
		hpa := operator.HorizontalPodAutoscaler(servicesName[i], instance)
		hpaName := servicesName[i]
		err = bootstrapClient.CreateResource(hpaName, hpa)
		if err != nil {
			log.Error(err, "failed to create operator's hpa", "Name", hpaName)
			return ctrl.Result{}, err
		}

	}

	if instance.Spec.ScaleCrossCloud == true && isWatchPod == false {
		isWatchPod = true
		time.Sleep(15 * time.Second)
		go r.WatchPod(ctx, req, instance)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PetClinicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1beta1.PetClinic{}).
		Complete(r)
}

func (r *PetClinicReconciler) WatchPod(ctx context.Context, req ctrl.Request, app *examplev1beta1.PetClinic) {

	log := r.Log.WithValues("pods in cluster", req.NamespacedName)

	for {

		pods := &corev1.PodList{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "",
			},
			ListMeta: metav1.ListMeta{
				SelfLink:           "",
				ResourceVersion:    "",
				Continue:           "",
				RemainingItemCount: nil,
			},
			Items: nil,
		}
		//
		err := r.List(ctx, pods)

		if err != nil {
			log.Info("No pods")
		}

		for i := 0; i < len(pods.Items); i++ {
			if pods.Items[i].Namespace != app.Namespace {
				continue
			} else {
				if pods.Items[i].Status.Phase == "Pending" {
					fmt.Println(pods.Items[0].Name)
					fmt.Println(pods.Items[0].Status.Reason)
					fmt.Println(pods.Items[0].Status.Message)
				}
			}

		}

		time.Sleep(10 * time.Second)
	}
}
