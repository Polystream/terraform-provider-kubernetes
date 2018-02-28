package kubernetes

import (
	api "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenDeploymentSpec(in api.DeploymentSpec) []interface{} {
	att := make(map[string]interface{})
	
	att["replicas"] = *in.Replicas
	att["selector"] = flattenLabelSelector(in.Selector)
	att["template"] = flattenPodTemplateSpec(in.Template)
	att["min_ready_seconds"] = in.MinReadySeconds
	att["progress_deadline_seconds"] = *in.ProgressDeadlineSeconds
	att["deployment_strategy"] = flattenDeploymentStrategy(in.Strategy)
	att["revision_history_limit"] = *in.RevisionHistoryLimit

	return []interface{}{att}
}

func flattenDeploymentStrategy(in api.DeploymentStrategy) []interface{} {
	att := make(map[string]interface{})

	att["type"] = string(in.Type)
	if(in.RollingUpdate != nil){
		update := make(map[string]interface{})
		update["max_unavailable"] = (*in.RollingUpdate.MaxUnavailable).String()
		update["max_surge"] = (*in.RollingUpdate.MaxSurge).String()
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
	if v, ok := m["replicas"]; ok {
		replicas := int32(v.(int))
		spec.Replicas = &replicas
	}
	if v, ok := m["selector"]; ok {
		spec.Selector = expandLabelSelector(v.([]interface{}))
	}
	if v, ok := m["template"]; ok {
		spec.Template = expandPodTemplateSpec(v.([]interface{}))
	}

	if v, ok := m["min_ready_seconds"]; ok {
		spec.MinReadySeconds = int32(v.(int))
	}
	if v, ok := m["deployment_strategy"]; ok {
		spec.Strategy = expandDeploymentStrategy(v.([]interface{}))
	}
	if v, ok := m["revision_history_limit"]; ok {
		limit := int32(v.(int))
		spec.RevisionHistoryLimit = &limit
	}
	if v, ok := m["progress_deadline_seconds"]; ok {
		seconds := int32(v.(int))
		spec.ProgressDeadlineSeconds = &seconds
	}
	return spec
}

func expandDeploymentStrategy(in []interface{}) api.DeploymentStrategy {
	if len(in) == 0 || in[0] == nil {
		return api.DeploymentStrategy{}
	}
	strategy := api.DeploymentStrategy{}
	m := in[0].(map[string]interface{})

	if v, ok := m["type"]; ok {
		strategy.Type = api.DeploymentStrategyType(v.(string))
	}

	if v, ok := m["rolling_update"]; ok {
		u := v.(map[string]interface{})
		update := api.RollingUpdateDeployment{}
		if v, ok := u["max_unavailable"]; ok {
			val := intstr.FromString(v.(string))
			update.MaxUnavailable = &val
		}
		if v, ok := u["max_surge"]; ok {
			val := intstr.FromString(v.(string))
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