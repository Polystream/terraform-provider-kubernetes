package kubernetes

import (
	api "k8s.io/api/batch/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenCronJobSpec(in api.CronJobSpec) []interface{} {
	att := make(map[string]interface{})
	
	att["schedule"] = in.Schedule
	if in.StartingDeadlineSeconds != nil {
		att["starting_deadline_seconds"] = *in.StartingDeadlineSeconds
	}
	att["concurrency_policy"] = string(in.ConcurrencyPolicy)
	if in.Suspend != nil {
		att["suspend"] = *in.Suspend
	}
	if in.SuccessfulJobsHistoryLimit != nil {
		att["successful_jobs_history_limit"] = *in.SuccessfulJobsHistoryLimit
	}
	if in.FailedJobsHistoryLimit != nil {
		att["failed_jobs_history_limit"] = *in.FailedJobsHistoryLimit
	}
	att["job_template"] = flattenJobTemplate(in.JobTemplate)

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

	if in.Parallelism != nil {
		att["parallelism"] = *in.Parallelism
	}
	if in.Completions != nil {
		att["completions"] = *in.Completions
	}
	if in.ActiveDeadlineSeconds != nil {
		att["active_deadline_seconds"] = *in.ActiveDeadlineSeconds
	}
	if in.BackoffLimit != nil {
		att["backoff_limit"] = *in.BackoffLimit
	}
	if in.Selector != nil {
		att["selector"] = flattenLabelSelector(in.Selector)
	}
	att["template"] = flattenPodTemplateSpec(in.Template)

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
	if v, ok := m["starting_deadline_seconds"].(int); ok {
		spec.StartingDeadlineSeconds = ptrToInt64(int64(v))
	}
	if v, ok := m["concurrency_policy"].(string); ok {
		spec.ConcurrencyPolicy = api.ConcurrencyPolicy(v)
	}
	if v, ok := m["suspend"].(bool); ok {
		spec.Suspend = ptrToBool(v)
	}
	if v, ok := m["successful_jobs_history_limit"].(int); ok {
		spec.SuccessfulJobsHistoryLimit = ptrToInt32(int32(v))
	}
	if v, ok := m["failed_jobs_history_limit"].(int); ok {
		spec.FailedJobsHistoryLimit = ptrToInt32(int32(v))
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
		spec.Parallelism = ptrToInt32(int32(v))
	}

	if v, ok := m["completions"].(int); ok {
		spec.Completions = ptrToInt32(int32(v))
	}

	if v, ok := m["active_deadline_seconds"].(int); ok {
		spec.ActiveDeadlineSeconds = ptrToInt64(int64(v))
	}

	if v, ok := m["backoff_limit"].(int); ok {
		spec.BackoffLimit = ptrToInt32(int32(v))
	}

	if v, ok := m["selector"].([]interface{}); ok {
		if len(v) > 0 {
			spec.Selector = expandLabelSelector(v)
		}
	}
	spec.ManualSelector = ptrToBool(false)

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