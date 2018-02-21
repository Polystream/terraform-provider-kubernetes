package kubernetes

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/apps/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesDeploymentCreate,
		Read:   resourceKubernetesDeploymentRead,
		Exists: resourceKubernetesDeploymentExists,
		Update: resourceKubernetesDeploymentUpdate,
		Delete: resourceKubernetesDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("Deployment", false),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				ForceNew: true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: deploymentSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	Deployment := &v1.Deployment{
		ObjectMeta: metadata,
		Spec: expandDeploymentSpec(d.Get("spec").([]interface{})),
	}
	log.Printf("[INFO] Creating new Deployment: %#v", Deployment)
	Deployment, err := conn.AppsV1().Deployments(metadata.Namespace).Create(Deployment)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new Deployment: %#v", Deployment)
	d.SetId(buildId(Deployment.ObjectMeta))
	name := Deployment.ObjectMeta.Name

	pending := make([]string, *Deployment.Spec.Replicas)
	for i := range pending {
		pending[i] = fmt.Sprintf("%v", i)
	}

	stateConf := &resource.StateChangeConf{
		Target:  []string{fmt.Sprintf("%v", *Deployment.Spec.Replicas)},
		Pending: pending,
		Timeout: 20 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			out, err := conn.AppsV1().Deployments(metadata.Namespace).Get(name, meta_v1.GetOptions{})
			if err != nil {
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "", err
			}

			statusPhase := fmt.Sprintf("%v", out.Status.ReadyReplicas)
			log.Printf("[DEBUG] Deployment %s ready replicas: %#v", out.Name, statusPhase)
			return out, statusPhase, nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		lastWarnings, wErr := getLastWarningsForObject(conn, Deployment.ObjectMeta, "Deployment", 3)
		if wErr != nil {
			return wErr
		}
		return fmt.Errorf("%s%s", err, stringifyEvents(lastWarnings))
	}

	return resourceKubernetesDeploymentRead(d, meta)
}

func resourceKubernetesDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Deployment %s", name)
	Deployment, err := conn.AppsV1().Deployments(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received Deployment: %#v", Deployment)
	err = d.Set("metadata", flattenMetadata(Deployment.ObjectMeta))
	if err != nil {
		return err
	}

	flattenedSpec := flattenDeploymentSpec(Deployment.Spec)
	log.Printf("[DEBUG] Flattened Deployment spec: %#v", flattenedSpec)
	err = d.Set("spec", flattenedSpec)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	
		namespace, name, err := idParts(d.Id())
		if err != nil {
			return err
		}
	
		ops := patchMetadata("metadata.0.", "/metadata/", d)
		data, err := ops.MarshalJSON()
		if err != nil {
			return fmt.Errorf("Failed to marshal update operations: %s", err)
		}
	
		log.Printf("[INFO] Updating Deployment %s: %s", d.Id(), ops)
	
		out, err := conn.AppsV1().Deployments(namespace).Patch(name, pkgApi.JSONPatchType, data)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Submitted updated Deployment: %#v", out)
	
		d.SetId(buildId(out.ObjectMeta))
		return resourceKubernetesDeploymentRead(d, meta)
}

func resourceKubernetesDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Deleting Deployment: %#v", name)
	err = conn.AppsV1().Deployments(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Deployment %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesDeploymentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Checking Deployment %s", name)
	_, err = conn.AppsV1().Deployments(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
