package kubernetes

import (
	api "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenDeploymentSpec(in api.DeploymentSpec) []interface{} {
	att := make(map[string]interface{})
	if in.Replicas != nil {
		att["replicas"] = *in.Replicas
	}
	att["selector"] = flattenLabelSelector(in.Selector)
	att["template"] = flattenPodTemplateSpec(in.Template)
	att["min_ready_seconds"] = in.MinReadySeconds
	if in.ProgressDeadlineSeconds != nil {
		att["progress_deadline_seconds"] = *in.ProgressDeadlineSeconds
	}
	att["deployment_strategy"] = flattenDeploymentStrategy(in.Strategy)
	if in.RevisionHistoryLimit != nil {
		att["revision_history_limit"] = *in.RevisionHistoryLimit
	}

	return []interface{}{att}
}

func flattenDeploymentStrategy(in api.DeploymentStrategy) []interface{} {
	att := make(map[string]interface{})

	att["type"] = string(in.Type)
	if(in.RollingUpdate != nil){
		update := make(map[string]interface{})
		if in.RollingUpdate.MaxUnavailable != nil {
			update["max_unavailable"] = (*in.RollingUpdate.MaxUnavailable).String()
		}
		if in.RollingUpdate.MaxSurge != nil {
			update["max_surge"] = (*in.RollingUpdate.MaxSurge).String()
		}
		att["rolling_update"] = []interface{}{update}
	}

	return []interface{}{att}
}

func expandDeploymentSpec(in []interface{}) api.DeploymentSpec {
	if len(in) == 0 || in[0] == nil {
		return api.DeploymentSpec{}
	}
	spec := api.DeploymentSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["replicas"].(int); ok {
		spec.Replicas = ptrToInt32(int32(v))
	}
	if v, ok := m["selector"].([]interface{}); ok {
		spec.Selector = expandLabelSelector(v)
	}
	if v, ok := m["template"].([]interface{}); ok {
		spec.Template = expandPodTemplateSpec(v)
	}

	if v, ok := m["min_ready_seconds"].(int); ok {
		spec.MinReadySeconds = int32(v)
	}
	if v, ok := m["deployment_strategy"].([]interface{}); ok {
		spec.Strategy = expandDeploymentStrategy(v)
	}
	if v, ok := m["revision_history_limit"].(int); ok {
		spec.RevisionHistoryLimit = ptrToInt32(int32(v))
	}
	if v, ok := m["progress_deadline_seconds"].(int); ok {
		spec.ProgressDeadlineSeconds = ptrToInt32(int32(v))
	}
	return spec
}

func expandDeploymentStrategy(in []interface{}) api.DeploymentStrategy {
	if len(in) == 0 || in[0] == nil {
		return api.DeploymentStrategy{}
	}
	strategy := api.DeploymentStrategy{}
	m := in[0].(map[string]interface{})

	if v, ok := m["type"].(string); ok {
		strategy.Type = api.DeploymentStrategyType(v)
	}

	if v, ok := m["rolling_update"].(map[string]interface{}); ok {
		update := api.RollingUpdateDeployment{}
		if v, ok := v["max_unavailable"].(string); ok {
			val := intstr.FromString(v)
			update.MaxUnavailable = &val
		}
		if v, ok := v["max_surge"].(string); ok {
			val := intstr.FromString(v)
			update.MaxSurge = &val
		}
	}

	return strategy
}

func patchDeploymentSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
	ops := make([]PatchOperation, 0)

	if d.HasChange(prefix + "replicas") {
		v := d.Get(prefix + "replicas").(int)
		ops = append(ops, &ReplaceOperation{
			Path:  pathPrefix + "/replicas",
			Value: v,
		})
	}

	if d.HasChange(prefix + "template.0.metadata") {
		metadataOps := patchMetadata(pathPrefix + "/template/metadata/", prefix + "template.0.metadata.0.", d)
		
		ops = append(ops, metadataOps...)
	}

	if d.HasChange(prefix + "template.0.spec") {
		podOps, err := patchPodSpec(pathPrefix + "/template/spec", prefix + "template.0.spec.0.", d)
		if err != nil {
			return nil, err
		}
		
		ops = append(ops, podOps...)
	}

	return ops, nil
}