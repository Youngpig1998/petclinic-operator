package operator

import (
	examplev1beta1 "github.com/Youngpig1998/petClinic-operator/api/v1beta1"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources/applications"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources/deployments"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources/horizontalpodautoscalers"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources/services"
	"github.com/Youngpig1998/petClinic-operator/iaw-shared-helpers/pkg/resources/statefulsets"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	oamv1beta1 "github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"strconv"
)

func Service(serviceName string, port int32, targetPort int32, nodePort int32) resources.Reconcileable {

	serviceSpec := corev1.ServiceSpec{
		Ports: []corev1.ServicePort{{
			Name: "http",
			Port: port,
			TargetPort: intstr.IntOrString{
				IntVal: targetPort,
				StrVal: strconv.Itoa(int(targetPort)),
			},
			//NodePort: nodePort,
		},
		},
		Selector: map[string]string{
			"app": serviceName,
		},
		Type: "ClusterIP",
	}

	objectMeta := metav1.ObjectMeta{
		Name: serviceName,
		Labels: map[string]string{
			"svc": serviceName,
		},
	}

	if nodePort != 0 {
		serviceSpec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: port,
				TargetPort: intstr.IntOrString{
					IntVal: targetPort,
					StrVal: strconv.Itoa(int(targetPort)),
				},
				NodePort: nodePort,
			},
			},
			Selector: map[string]string{
				"app": serviceName,
			},
			Type: "NodePort",
		}
	}

	if serviceName == "mysql" {
		objectMeta = metav1.ObjectMeta{
			Name: serviceName,
		}
		serviceSpec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "tcp",
				Port: port,
				TargetPort: intstr.IntOrString{
					IntVal: targetPort,
					StrVal: strconv.Itoa(int(targetPort)),
				},
			},
			},
			Selector: map[string]string{
				"app": serviceName,
			},
			Type: "ClusterIP",
			//ClusterIP: "None",
		}

	}

	service := &corev1.Service{
		ObjectMeta: objectMeta,
		Spec:       serviceSpec,
	}

	return services.From(service)
}

func StatefulSet(servicesName string, app *examplev1beta1.PetClinic) resources.Reconcileable {

	statefulSetName := servicesName

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: statefulSetName,
			Labels: map[string]string{
				"app": statefulSetName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: pointer.Int32Ptr(app.Spec.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": statefulSetName,
				},
			},
			ServiceName: statefulSetName,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": statefulSetName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  statefulSetName,
						Image: "mysql:5.7",
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "petclinic",
							},
							{
								Name:  "MYSQL_DATABASE",
								Value: "petclinic",
							},
						},
						ImagePullPolicy: "IfNotPresent",
						VolumeMounts: []corev1.VolumeMount{{
							Name:      servicesName,
							MountPath: "/var/lib/mysql",
						},
						},
					},
					},
					ServiceAccountName: "nfs-provisioner",
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: servicesName,
					Annotations: map[string]string{
						"volume.beta.kubernetes.io/storage-class": "managed-nfs-storage",
					},
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse("10Gi"),
						},
					},
				},
			},
			},
		},
	}

	return statefulsets.From(statefulSet)
}

func Deployment(serviceName string, app *examplev1beta1.PetClinic) resources.Reconcileable {

	return deployments.From(createDeployment(serviceName, app))
}

func HorizontalPodAutoscaler(horizontalPodAutoscalerName string, app *examplev1beta1.PetClinic) resources.Reconcileable {

	horizontalPodAutoscaler := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      horizontalPodAutoscalerName,
			Namespace: app.Namespace,
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       horizontalPodAutoscalerName,
				APIVersion: "apps/v1",
			},
			MinReplicas: pointer.Int32Ptr(app.Spec.Replicas),
			MaxReplicas: 5,
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Type: "Resource",
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2beta2.MetricTarget{
							Type:               "Utilization",
							AverageUtilization: pointer.Int32Ptr(40),
						},
					},
				},
			},
		},
	}

	return horizontalpodautoscalers.From(horizontalPodAutoscaler)
}

func getKeys(m map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

func Application(serviceName string, app *examplev1beta1.PetClinic) resources.Reconcileable {

	deploy := createDeployment(serviceName, app)

	application := &oamv1beta1.Application{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "petclinic-customer",
			Namespace: app.Namespace,
		},
		Spec: oamv1beta1.ApplicationSpec{
			Components: []common.ApplicationComponent{
				{
					Name: "petclinic-customer",
					Type: "k8s-objects",
					Properties: &runtime.RawExtension{
						Object: deploy,
					},
				},
			},

			Policies: []oamv1beta1.AppPolicy{
				{
					Name: "topology-default",
					Type: "topology",
					Properties: &runtime.RawExtension{
						Raw: []byte(`{"namespace":"default","clusters":"['k8s-master']"}`),
					},
				},
			},
			Workflow: &oamv1beta1.Workflow{
				Steps: []oamv1beta1.WorkflowStep{
					{
						Name: "deploy2default",
						Type: "deploy",
						Properties: &runtime.RawExtension{
							Raw: []byte(`{"policies":"['topology-default']"}`),
						},
					},
				},
			},
		},
	}

	return applications.From(application)
}

func createDeployment(serviceName string, app *examplev1beta1.PetClinic) *appsv1.Deployment {

	deployName := serviceName
	imageName := "youngpig/spring-petclinic-" + serviceName + "-service:1.0.0.RELEASE"

	env := []corev1.EnvVar{
		{
			Name:  "JAVA_OPTS",
			Value: "-XX:MinRAMPercentage=50.0 -XX:MaxRAMPercentage=80.0 -XX:+HeapDumpOnOutOfMemoryError",
		},
		{
			Name:  "SERVER_PORT",
			Value: "8080",
		},
	}

	if app.Spec.MysqlActive == true {
		env = []corev1.EnvVar{
			{
				Name:  "JAVA_OPTS",
				Value: "-XX:MinRAMPercentage=50.0 -XX:MaxRAMPercentage=80.0 -XX:+HeapDumpOnOutOfMemoryError",
			},
			{
				Name:  "SERVER_PORT",
				Value: "8080",
			},
			{
				Name:  "SPRING_PROFILES_ACTIVE",
				Value: "mysql",
			},
			{
				Name:  "DATASOURCE_URL",
				Value: "jdbc:mysql://mysql/petclinic",
			},
			{
				Name:  "DATASOURCE_USERNAME",
				Value: "root",
			},
			{
				Name:  "DATASOURCE_PASSWORD",
				Value: "petclinic",
			},
			{
				Name:  "DATASOURCE_INIT_MODE",
				Value: "always",
			},
		}
	}

	if serviceName == "web" {
		env = []corev1.EnvVar{
			{
				Name:  "JAVA_OPTS",
				Value: "-XX:MinRAMPercentage=50.0 -XX:MaxRAMPercentage=80.0 -XX:+HeapDumpOnOutOfMemoryError",
			},
			{
				Name:  "SERVER_PORT",
				Value: "8080",
			},
			{
				Name:  "VISITS_SERVICE_ENDPOINT",
				Value: "http://visits:8080",
			},
			{
				Name:  "CUSTOMERS_SERVICE_ENDPOINT",
				Value: "http://customers:8080",
			},
		}
		imageName = "youngpig/spring-petclinic-" + serviceName + "-app:1.0.0.RELEASE"
	} else if serviceName == "gateway" {
		env = []corev1.EnvVar{
			{
				Name:  "JAVA_OPTS",
				Value: "-XX:MinRAMPercentage=50.0 -XX:MaxRAMPercentage=80.0 -XX:+HeapDumpOnOutOfMemoryError",
			},
			{
				Name:  "SERVER_PORT",
				Value: "8080",
			},
			{
				Name:  "WEB_APP_ENDPOINT",
				Value: "http://web:8080",
			},
			{
				Name:  "VETS_SERVICE_ENDPOINT",
				Value: "http://vets:8080",
			},
			{
				Name:  "VISITS_SERVICE_ENDPOINT",
				Value: "http://visits:8080",
			},
			{
				Name:  "CUSTOMERS_SERVICE_ENDPOINT",
				Value: "http://customers:8080",
			},
		}
		imageName = "youngpig/spring-petclinic-cloud-" + serviceName + ":1.0.0.RELEASE"
	}

	// Instantialize the data structure
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: app.Namespace,
			Name:      deployName,
			//Labels: map[string]string{
			//	"io.kompose.service": deployName,
			//},
		},
		Spec: appsv1.DeploymentSpec{
			// The replica is computed
			Replicas: pointer.Int32Ptr(app.Spec.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deployName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deployName,
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{{
						Image:           "busybox:1.30",
						ImagePullPolicy: "IfNotPresent",
						Name:            "test-mysql",
						Command:         []string{"sh", "-c", "until ping mysql -c 1 ; do echo waiting for mysql...;sleep 2;done;"},
					}},
					RestartPolicy: corev1.RestartPolicy("Always"),
					Containers: []corev1.Container{{
						Image:           imageName,
						ImagePullPolicy: "IfNotPresent",
						Name:            deployName,
						Env:             env,
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/actuator/health",
									Port: intstr.IntOrString{
										Type:   intstr.Int,
										IntVal: 8080,
										StrVal: "8080",
									},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       10,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu":    resource.MustParse("2"),
								"memory": resource.MustParse("1024Mi"),
							},
							Requests: corev1.ResourceList{
								"cpu":    resource.MustParse("1"),
								"memory": resource.MustParse("512Mi"),
							},
						},
					}},
				},
			},
		},
	}

	return deployment
}
