package kubernetes

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/apps/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesDaemonSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesDaemonsetCreate,
		Read:   resourceKubernetesDaemonsetRead,
		Exists: resourceKubernetesDaemonsetExists,
		Update: resourceKubernetesDaemonsetUpdate,
		Delete: resourceKubernetesDaemonsetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("Daemonset", false),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				ForceNew: true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: daemonsetSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesDaemonsetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	daemonset := &v1.DaemonSet{
		ObjectMeta: metadata,
		Spec: expandDaemonsetSpec(d.Get("spec").([]interface{})),
	}
	log.Printf("[INFO] Creating new Daemonset: %#v", daemonset)
	daemonset, err := conn.AppsV1().DaemonSets(metadata.Namespace).Create(daemonset)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new Daemonset: %#v", daemonset)
	d.SetId(buildId(daemonset.ObjectMeta))

	return resourceKubernetesDaemonsetRead(d, meta)
}

func resourceKubernetesDaemonsetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Daemonset %s", name)
	daemonset, err := conn.AppsV1().DaemonSets(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received Daemonset: %#v", daemonset)
	err = d.Set("metadata", flattenMetadata(daemonset.ObjectMeta))
	if err != nil {
		return err
	}

	flattenedSpec := flattenDaemonsetSpec(daemonset.Spec)
	log.Printf("[DEBUG] Flattened Daemonset spec: %#v", flattenedSpec)
	err = d.Set("spec", flattenedSpec)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesDaemonsetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	
		namespace, name, err := idParts(d.Id())
		if err != nil {
			return err
		}
	
		metadata := expandMetadata(d.Get("metadata").([]interface{}))
		daemonset := &v1.DaemonSet{
			ObjectMeta: metadata,
			Spec: expandDaemonsetSpec(d.Get("spec").([]interface{})),
		}

		data, err := json.Marshal(daemonset)
		if err != nil {
			return fmt.Errorf("Failed to marshal update operations: %s", err)
		}
	
		log.Printf("[INFO] Updating Daemonset %s: %s", d.Id(), daemonset)
	
		out, err := conn.AppsV1().DaemonSets(namespace).Patch(name, pkgApi.StrategicMergePatchType, data)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Submitted updated Daemonset: %#v", out)
	
		d.SetId(buildId(out.ObjectMeta))
		return resourceKubernetesDaemonsetRead(d, meta)
}

func resourceKubernetesDaemonsetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Deleting Daemonset: %#v", name)
	err = conn.AppsV1().DaemonSets(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Daemonset %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesDaemonsetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Checking Daemonset %s", name)
	_, err = conn.AppsV1().DaemonSets(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
