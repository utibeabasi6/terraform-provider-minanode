package minanode

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func nodeCreate(d *schema.ResourceData, m interface{}) error {
	kubeconfig := m.(Params).kubeconfig
	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	replicas := d.Get("replicas").(int)
	if replicas == 0 {
		replicas = 1
	}
	privKey := d.Get("privkey").(string)

	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return errors.New("unable to authenticate with the kubeconfig file")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("error while creating client")
	}
	deployment, err := CreateNode(clientset, namespace, name, privKey, int32(replicas))
	if err == nil {
		d.Set("deployment", deployment.Name)
		d.SetId(name)
	}
	return err
}

func nodeUpdate(d *schema.ResourceData, m interface{}) error {
	return nodeRead(d, m)
}

func nodeRead(d *schema.ResourceData, m interface{}) error {
	kubeconfig := m.(Params).kubeconfig
	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	id := d.Id()

	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	replicas := d.Get("replicas").(int)
	if replicas == 0 {
		replicas = 1
	}
	privKey := d.Get("privkey").(string)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return errors.New("unable to authenticate with the kubeconfig file")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("error while creating client")
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	result, err := deploymentsClient.Get(context.TODO(), id, metav1.GetOptions{})
	if err == nil {
		d.Set("deployment", result.Name)
		d.Set("name", name)
		d.Set("privkey", privKey)
		d.Set("replicas", replicas)
		d.Set("namespace", namespace)
	}
	return err
}

func nodeDelete(d *schema.ResourceData, m interface{}) error {
	kubeconfig := m.(Params).kubeconfig
	name := d.Id()
	namespace := d.Get("namespace").(string)
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return errors.New("unable to authenticate with the kubeconfig file")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("unable to create kubernetes client")
	}
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err = deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err == nil {
		d.SetId("")
	}
	return err
}

func resourceNode() *schema.Resource {
	return &schema.Resource{
		Read:   nodeRead,
		Create: nodeCreate,
		Update: nodeUpdate,
		Delete: nodeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the deployment",
			},
			"privkey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private key for the key pair",
			},
			"replicas": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of nodes to deploy",
				Default:     1,
			},
			"namespace": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The namespace to deploy resources into",
				Default:     "default",
			},
			"deployment": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
