package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	pkgApi "k8s.io/apimachinery/pkg/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/rbac/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesClusterRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesClusterRoleCreate,
		Read:   resourceKubernetesClusterRoleRead,
		Exists: resourceKubernetesClusterRoleExists,
		Update: resourceKubernetesClusterRoleUpdate,
		Delete: resourceKubernetesClusterRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": metadataSchema("clusterRole", false),
			"rule": {
				Type:        schema.TypeList,
				Description: "Subjects defines the entities to bind a ClusterRole to.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: rbacRoleSchema(),
				},
			},
		},
	}
}

func resourceKubernetesClusterRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	binding := &v1.ClusterRole{
		ObjectMeta: metadata,
		Rules: expandRBACRules(d.Get("rule").([]interface{})),
	}
	log.Printf("[INFO] Creating new ClusterRole: %#v", binding)
	binding, err := conn.Rbac().ClusterRoles().Create(binding)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new ClusterRole: %#v", binding)
	d.SetId(metadata.Name)

	return resourceKubernetesClusterRoleRead(d, meta)
}

func resourceKubernetesClusterRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	name := d.Id()
	log.Printf("[INFO] Reading ClusterRole %s", name)
	role, err := conn.Rbac().ClusterRoles().Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received ClusterRole: %#v", role)
	err = d.Set("metadata", flattenMetadata(role.ObjectMeta))
	if err != nil {
		return err
	}

	flattenedRules := flattenRBACRules(role.Rules)
	log.Printf("[DEBUG] Flattened ClusterRole ruleRef: %#v", flattenedRules)
	err = d.Set("rule", flattenedRules)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesClusterRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	name := d.Id()

	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("rule") {
		diffOps := patchRbacRule(d)
		ops = append(ops, diffOps...)
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}
	log.Printf("[INFO] Updating ClusterRole %q: %v", name, string(data))
	out, err := conn.Rbac().ClusterRoles().Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update ClusterRole: %s", err)
	}
	log.Printf("[INFO] Submitted updated ClusterRole: %#v", out)
	d.SetId(out.ObjectMeta.Name)

	return resourceKubernetesClusterRoleRead(d, meta)
}

func resourceKubernetesClusterRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	name := d.Id()
	log.Printf("[INFO] Deleting ClusterRole: %#v", name)
	err := conn.Rbac().ClusterRoles().Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] ClusterRole %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesClusterRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	name := d.Id()
	log.Printf("[INFO] Checking ClusterRole %s", name)
	_, err := conn.Rbac().ClusterRoles().Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
