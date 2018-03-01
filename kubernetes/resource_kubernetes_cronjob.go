package kubernetes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/api/batch/v1beta1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesCronJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesCronJobCreate,
		Read:   resourceKubernetesCronJobRead,
		Exists: resourceKubernetesCronJobExists,
		Update: resourceKubernetesCronJobUpdate,
		Delete: resourceKubernetesCronJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("CronJob", false),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec of the pod owned by the cluster",
				Required:    true,
				ForceNew: true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: cronJobSpecFields(),
				},
			},
		},
	}
}

func resourceKubernetesCronJobCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	CronJob := &v1beta1.CronJob{
		ObjectMeta: metadata,
		Spec: expandCronJobSpec(d.Get("spec").([]interface{})),
	}
	log.Printf("[INFO] Creating new CronJob: %#v", CronJob)
	CronJob, err := conn.BatchV1beta1().CronJobs(metadata.Namespace).Create(CronJob)

	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new CronJob: %#v", CronJob)
	d.SetId(buildId(CronJob.ObjectMeta))

	return resourceKubernetesCronJobRead(d, meta)
}

func resourceKubernetesCronJobRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading CronJob %s", name)
	CronJob, err := conn.BatchV1beta1().CronJobs(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	log.Printf("[INFO] Received CronJob: %#v", CronJob)
	err = d.Set("metadata", flattenMetadata(CronJob.ObjectMeta))
	if err != nil {
		return err
	}

	flattenedSpec := flattenCronJobSpec(CronJob.Spec)
	log.Printf("[DEBUG] Flattened CronJob spec: %#v", flattenedSpec)
	err = d.Set("spec", flattenedSpec)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesCronJobUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)
	
		namespace, name, err := idParts(d.Id())
		if err != nil {
			return err
		}
	
		ops := patchMetadata("metadata.0.", "/metadata/", d)
		if d.HasChange("spec") {
			specOps, err := patchCronJobSpec("/spec", "spec.0.", d)
			if err != nil {
				return err
			}
			ops = append(ops, specOps...)
		}
		data, err := ops.MarshalJSON()
		if err != nil {
			return fmt.Errorf("Failed to marshal update operations: %s", err)
		}
	
		log.Printf("[INFO] Updating CronJob %s: %s", d.Id(), ops)
	
		out, err := conn.BatchV1beta1().CronJobs(namespace).Patch(name, pkgApi.JSONPatchType, data)
		if err != nil {
			return err
		}
		log.Printf("[INFO] Submitted updated CronJob: %#v", out)
	
		d.SetId(buildId(out.ObjectMeta))
		return resourceKubernetesCronJobRead(d, meta)
}

func resourceKubernetesCronJobDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Deleting CronJob: %#v", name)
	err = conn.BatchV1beta1().CronJobs(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Printf("[INFO] CronJob %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesCronJobExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	log.Printf("[INFO] Checking CronJob %s", name)
	_, err = conn.BatchV1beta1().CronJobs(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}
