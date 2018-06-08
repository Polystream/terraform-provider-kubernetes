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

func resourceKubernetesRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesRoleCreate,
		Read:   resourceKubernetesRoleRead,
		Exists: resourceKubernetesRoleExists,
		Update: resourceKubernetesRoleUpdate,
		Delete: resourceKubernetesRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("role", false),
			"rule": {
				Type:        schema.TypeList,
				Description: "Subjects defines the entities to bind a Role to.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: rbacRoleSchema(),
				},
			},
		},
	}
}

func resourceKubernetesRoleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	binding := &v1.Role{
		ObjectMeta: metadata,
		Rules: expandRBACRules(d.Get("rule").([]interface{})),
	}
	log.Printf("[INFO] Creating new Role: %#v", binding)
	binding, err := conn.Rbac().Roles(metadata.Namespace).Create(binding)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new Role: %#v", binding)
	d.SetId(buildId(binding.ObjectMeta))

	return resourceKubernetesRoleRead(d, meta)
}

func resourceKubernetesRoleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Role %s", name)
	role, err := conn.Rbac().Roles(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received Role: %#v", role)
	err = d.Set("metadata", flattenMetadata(role.ObjectMeta))
	if err != nil {
		return err
	}

	flattenedRules := flattenRBACRules(role.Rules)
	log.Printf("[DEBUG] Flattened Role ruleRef: %#v", flattenedRules)
	err = d.Set("rule", flattenedRules)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	ops := patchMetadata("metadata.0.", "/metadata/", d)
	if d.HasChange("rule") {
		diffOps := patchRbacRule(d)
		ops = append(ops, diffOps...)
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}
	log.Printf("[INFO] Updating Role %q: %v", name, string(data))
	out, err := conn.Rbac().Roles(namespace).Patch(name, pkgApi.JSONPatchType, data)
	if err != nil {
		return fmt.Errorf("Failed to update Role: %s", err)
	}
	log.Printf("[INFO] Submitted updated Role: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesRoleRead(d, meta)
}

func resourceKubernetesRoleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Role: %#v", name)
	err = conn.Rbac().Roles(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] Role %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking Role %s", name)
	_, err = conn.Rbac().Roles(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
