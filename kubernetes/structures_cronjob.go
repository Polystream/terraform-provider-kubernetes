package kubernetes

import (
	api "k8s.io/api/batch/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenCronJobSpec(in api.CronJobSpec) []interface{} {
	att := make(map[string]interface{})
	
	att["schedule"] = in.Schedule
	att["starting_deadline_seconds"] = *in.StartingDeadlineSeconds
	att["concurrency_policy"] = string(in.ConcurrencyPolicy)
	att["suspend"] = *in.Suspend
	att["successful_jobs_history_limit"] = *in.SuccessfulJobsHistoryLimit
	att["failed_jobs_history_limit"] = *in.FailedJobsHistoryLimit
	att["template"] = flattenJobTemplate(in.JobTemplate)

	return []interface{}{att}
}

func flattenJobTemplate(in api.JobTemplateSpec) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = flattenMetadata(in.ObjectMeta)
	att["spec"] = flattenJobSpec(in.Spec)

	return []interface{}{att}
}

func flattenJobSpec(in batchv1.JobSpec) []interface{} {
	att := make(map[string]interface{})

	att["parallelism"] = *in.Parallelism
	att["completions"] = *in.Completions
	att["active_deadline_seconds"] = *in.ActiveDeadlineSeconds
	att["backoff_limit"] = *in.BackoffLimit
	att["selector"] = flattenLabelSelector(in.Selector)
	att["job_template"] = flattenPodTemplateSpec(in.Template)

	return []interface{}{att}
}

func expandCronJobSpec(in []interface{}) api.CronJobSpec {
	if len(in) == 0 || in[0] == nil {
		return api.CronJobSpec{}
	}
	spec := api.CronJobSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["schedule"].(string); ok {
		spec.Schedule = v
	}
	if v, ok := m["starting_deadline_seconds"]; ok {
		time := int64(v.(int))
		spec.StartingDeadlineSeconds = &time
	}
	if v, ok := m["concurrency_policy"]; ok {
		spec.ConcurrencyPolicy = api.ConcurrencyPolicy(v.(string))
	}
	if v, ok := m["suspend"].(bool); ok {
		spec.Suspend = &v
	}
	if v, ok := m["successful_jobs_history_limit"].(int); ok {
		limit := int32(v)
		spec.SuccessfulJobsHistoryLimit = &limit
	}
	if v, ok := m["failed_jobs_history_limit"].(int); ok {
		limit := int32(v)
		spec.FailedJobsHistoryLimit = &limit
	}
	if v, ok := m["job_template"].([]interface{}); ok {
		spec.JobTemplate = expandJobTemplateSpec(v)
	}
	return spec
}

func expandJobTemplateSpec(in []interface{}) api.JobTemplateSpec {
	if len(in) == 0 || in[0] == nil {
		return api.JobTemplateSpec{}
	}
	spec := api.JobTemplateSpec{}
	m := in[0].(map[string]interface{})

	if v, ok := m["metadata"].([]interface{}); ok {
		spec.ObjectMeta = expandMetadata(v)
	}

	if v, ok := m["spec"].([]interface{}); ok {
		spec.Spec = expandJobSpec(v)
	}

	return spec
}

func expandJobSpec(in []interface{}) batchv1.JobSpec {
	if len(in) == 0 || in[0] == nil {
		return batchv1.JobSpec{}
	}
	spec := batchv1.JobSpec{}
	m := in[0].(map[string]interface{})

	if v, ok := m["parallelism"].(int); ok {
		p := int32(v)
		spec.Parallelism = &p
	}

	if v, ok := m["completions"].(int); ok {
		p := int32(v)
		spec.Completions = &p
	}

	if v, ok := m["active_deadline_seconds"].(int); ok {
		p := int64(v)
		spec.ActiveDeadlineSeconds = &p
	}

	if v, ok := m["backoff_limit"].(int); ok {
		p := int32(v)
		spec.BackoffLimit = &p
	}

	if v, ok := m["selector"].([]interface{}); ok {
		if len(v) > 0 {
			spec.Selector = expandLabelSelector(v)
		}
	}
	manualSelector := false
	spec.ManualSelector = &manualSelector

	if v, ok := m["template"].([]interface{}); ok {
		spec.Template = expandPodTemplateSpec(v)
	}

	return spec
}

func patchCronJobSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)

	//TODO

	return ops, nil
}