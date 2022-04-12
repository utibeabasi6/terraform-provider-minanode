package minanode

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Params struct {
	kubeconfig string
}

func int32Ptr(i int32) *int32 { return &i }

func CreateNode(clientset *kubernetes.Clientset, namespace string, name string, privKey string, replicas int32) (*appsv1.Deployment, error) {
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: "nginx", // "minaprotocol/mina-daemon:latest",
							//Args:  []string{"daemon", "--external-port", "8302"},
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{MountPath: "/keys", Name: "key-pair"},
								apiv1.VolumeMount{MountPath: "/root/.mina-config", Name: "mina-config"},
							},
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{Name: "MINA_PRIVKEY_PASS", Value: privKey},
								apiv1.EnvVar{Name: "LOG_LEVEL", Value: "Info"},
								apiv1.EnvVar{Name: "FILE_LOG_LEVEL", Value: "Debug"},
								apiv1.EnvVar{Name: "EXTRA_FLAGS", Value: " -block-producer-key /keys/my-wallet"},
								apiv1.EnvVar{Name: "PEER_LIST_URL", Value: "https://storage.googleapis.com/mina-seed-lists/mainnet_seeds.txt"},
							},
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8302,
								},
							},
						},
					},
					InitContainers: []apiv1.Container{
						apiv1.Container{
							Name:  name + "-key-pair",
							Image: "minaprotocol/mina-generate-keypair:1.2.0-fe51f1e",
							Args:  []string{"--privkey-path", "/keys/my-wallet"},
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{
									Name:      "key-pair",
									MountPath: "/keys",
								},
							},
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{Name: "CODA_PRIVKEY_PASS", Value: privKey},
							},
						},
						apiv1.Container{
							Name:    name + "-keys-config",
							Image:   "alpine",
							Command: []string{"sh", "-c", "cd ~", "chmod 700 /keys", "chmod 600 /keys/my-wallet"},
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{
									Name:      "key-pair",
									MountPath: "/keys",
								},
							},
							Env: []apiv1.EnvVar{
								apiv1.EnvVar{Name: "CODA_PRIVKEY_PASS", Value: privKey},
							},
						},
					},
					Volumes: []apiv1.Volume{
						apiv1.Volume{
							Name: "key-pair",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
						apiv1.Volume{
							Name: "mina-config",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	return result, err
}
