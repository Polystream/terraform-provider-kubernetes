package kubernetes

import (
	"k8s.io/api/core/v1"
	api "k8s.io/api/apps/v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenStatefulsetSpec(in api.StatefulSetSpec) []interface{} {
	att := make(map[string]interface{})
	
	if in.Replicas != nil {
		att["replicas"] = *in.Replicas
	}

	att["selector"] = flattenLabelSelector(in.Selector)

	att["template"] = flattenPodTemplateSpec(in.Template)

	if(in.VolumeClaimTemplates != nil && len(in.VolumeClaimTemplates) > 0){
		att["volume_claim_template"] = flattenVolumeClaimTemplates(in.VolumeClaimTemplates)
	}

	att["service_name"] = in.ServiceName
	att["pod_management_policy"] = string(in.PodManagementPolicy)
	att["update_strategy"] = flattenStatefulsetUpdateStrategy(in.UpdateStrategy)
	if in.RevisionHistoryLimit != nil {
		att["revision_history_limit"] = *in.RevisionHistoryLimit
	}

	return []interface{}{att}
}

func flattenPodTemplateSpec(in v1.PodTemplateSpec) []interface{} {
	att := make(map[string]interface{})
	att["metadata"] = flattenMetadata(in.ObjectMeta)
	att["spec"], _ = flattenPodSpec(in.Spec)
	return []interface{}{att}
}

func flattenVolumeClaimTemplates(in []v1.PersistentVolumeClaim) []interface{} {
	if in == nil {
		return []interface{}{}
	}
	att := make([]interface{}, len(in), len(in))
	for i, n := range in {
		m := make(map[string]interface{})
		m["metadata"] = flattenMetadata(n.ObjectMeta)
		m["spec"] = flattenPersistentVolumeClaimSpec(n.Spec)
		att[i] = m
	}
	return att
}

func flattenStatefulsetUpdateStrategy(in api.StatefulSetUpdateStrategy) []interface{} {
	att := make(map[string]interface{})
	att["type"] = string(in.Type)
	if(in.RollingUpdate != nil){
		update := make(map[string]interface{})
		if in.RollingUpdate.Partition != nil {
			update["partition"] = *in.RollingUpdate.Partition
		}
		att["rolling_update"] = []interface{}{update}
	}
	return []interface{}{att}
}

func expandStatefulsetSpec(in []interface{}) api.StatefulSetSpec {
	if len(in) == 0 || in[0] == nil {
		return api.StatefulSetSpec{}
	}
	spec := api.StatefulSetSpec{}
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
	if v, ok := m["volume_claim_template"].([]interface{}); ok {
		spec.VolumeClaimTemplates = expandVolumeClaimTemplate(v)
	}
	if v, ok := m["service_name"].(string); ok {
		spec.ServiceName = v
	}
	if v, ok := m["pod_management_policy"].(string); ok {
		spec.PodManagementPolicy = api.PodManagementPolicyType(v)
	}
	if v, ok := m["update_strategy"].([]interface{}); ok {
		spec.UpdateStrategy = expandStatefulSetUpdateStrategy(v)
	}
	if v, ok := m["revision_history_limit"].(int); ok {
		spec.RevisionHistoryLimit = ptrToInt32(int32(v))
	}
	return spec
}

func expandPodTemplateSpec(in []interface{}) v1.PodTemplateSpec {
	if len(in) == 0 || in[0] == nil {
		return v1.PodTemplateSpec{}
	}
	spec := v1.PodTemplateSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["metadata"].([]interface{}); ok {
		spec.ObjectMeta = expandMetadata(v)
	}
	if v, ok := m["spec"].([]interface{}); ok {
		spec.Spec, _ = expandPodSpec(v)
	}

	return spec
}

func expandVolumeClaimTemplate(in []interface{}) []v1.PersistentVolumeClaim {
	if len(in) == 0 || in[0] == nil {
		return []v1.PersistentVolumeClaim{}
	}
	claims := []v1.PersistentVolumeClaim{}
	for i := range in {
		claim := v1.PersistentVolumeClaim{}
		m := in[i].(map[string]interface{})
		if v, ok := m["metadata"].([]interface{}); ok {
			claim.ObjectMeta = expandMetadata(v)
		}
		if v, ok := m["spec"].([]interface{}); ok {
			claim.Spec, _ = expandPersistentVolumeClaimSpec(v)
		}
		claims = append(claims, claim)
	}
	return claims
}

func expandStatefulSetUpdateStrategy(in []interface{}) api.StatefulSetUpdateStrategy {
	if len(in) == 0 || in[0] == nil {
		return api.StatefulSetUpdateStrategy{}
	}
	strategy := api.StatefulSetUpdateStrategy{}
	m := in[0].(map[string]interface{})
	if v, ok := m["type"].(string); ok {
		strategy.Type = api.StatefulSetUpdateStrategyType(v)
	}
	if v, ok := m["rolling_update"]; ok {
		x := v.([]interface{})
		if len(x) > 0 && x[0] == nil {
			m = x[0].(map[string]interface{})
			if v, ok := m["partition"].(int); ok {
				strategy.RollingUpdate = &api.RollingUpdateStatefulSetStrategy{
					Partition: ptrToInt32(int32(v)),
				}
			}
		}
	}

	return strategy
}

func patchStatefultsetSpec(pathPrefix, prefix string, d *schema.ResourceData) (PatchOperations, error) {
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