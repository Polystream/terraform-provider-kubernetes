package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/apps/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesStatefulSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesStatefulsetCreate,
		Read:   resourceKubernetesStatefulsetRead,
		Exists: resourceKubernetesStatefulsetExists,
		Update: resourceKubernetesStatefulsetUpdate,
		Delete: resourceKubernetesStatefulsetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("statefulset", false),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: statefulsetSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesStatefulsetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	statefulset := &v1.StatefulSet{
		ObjectMeta: metadata,
		Spec: expandStatefulsetSpec(d.Get("spec").([]interface{})),
	}
	log.Printf("[INFO] Creating new Statefulset: %#v", statefulset)
	statefulset, err := conn.AppsV1().StatefulSets(metadata.Namespace).Create(statefulset)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new Statefulset: %#v", statefulset)
	d.SetId(buildId(statefulset.ObjectMeta))

	return resourceKubernetesStatefulsetRead(d, meta)
}

func resourceKubernetesStatefulsetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Statefulset %s", name)
	statefulset, err := conn.AppsV1().StatefulSets(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received Statefulset: %#v", statefulset)
	err = d.Set("metadata", flattenMetadata(statefulset.ObjectMeta))
	if err != nil {
		return err
	}

	//TODO
	// flattenedRules := flattenRBACRules(role.Rules)
	// log.Printf("[DEBUG] Flattened Statefulset ruleRef: %#v", flattenedRules)
	// err = d.Set("rule", flattenedRules)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func resourceKubernetesStatefulsetUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Update not implemented")

	// conn := meta.(*kubernetes.Clientset)

	// name := d.Id()

	// ops := patchMetadata("metadata.0.", "/metadata/", d)
	// if d.HasChange("rule") {
	// 	diffOps := patchRbacRule(d)
	// 	ops = append(ops, diffOps...)
	// }
	// data, err := ops.MarshalJSON()
	// if err != nil {
	// 	return fmt.Errorf("Failed to marshal update operations: %s", err)
	// }
	// log.Printf("[INFO] Updating ClusterRole %q: %v", name, string(data))
	// out, err := conn.Rbac().ClusterRoles().Patch(name, pkgApi.JSONPatchType, data)
	// if err != nil {
	// 	return fmt.Errorf("Failed to update ClusterRole: %s", err)
	// }
	// log.Printf("[INFO] Submitted updated ClusterRole: %#v", out)
	// d.SetId(out.ObjectMeta.Name)

	// return resourceKubernetesClusterRoleRead(d, meta)
}

func resourceKubernetesStatefulsetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Deleting Statefulset: %#v", name)
	err = conn.AppsV1().StatefulSets(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Statefulset %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesStatefulsetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Checking Statefulset %s", name)
	_, err = conn.AppsV1().StatefulSets(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
