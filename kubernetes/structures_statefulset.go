package kubernetes

import (
	"k8s.io/api/core/v1"
	api "k8s.io/api/apps/v1"
)

func flattenStatefulsetSpec(in api.StatefulSetSpec) []interface{} {
	att := make(map[string]interface{})
	
	att["replicas"] = *in.Replicas

	att["selector"] = flattenLabelSelector(in.Selector)

	att["template"] = flattenPodTemplateSpec(in.Template)

	if(in.VolumeClaimTemplates != nil && len(in.VolumeClaimTemplates) > 0){
		att["volume_claim_template"] = flattenVolumeClaimTemplates(in.VolumeClaimTemplates)
	}

	att["service_name"] = in.ServiceName
	att["pod_management_policy"] = string(in.PodManagementPolicy)
	att["update_strategy"] = flattenStatefulsetUpdateStrategy(in.UpdateStrategy)
	att["revision_history_limit"] = *in.RevisionHistoryLimit

	return []interface{}{att}
}

func flattenPodTemplateSpec(in v1.PodTemplateSpec) []interface{} {
	att := make(map[string]interface{})
	att["metadata"] = flattenMetadata(in.ObjectMeta)
	att["spec"], _ = flattenPodSpec(in.Spec)
	return []interface{}{att}
}

func flattenVolumeClaimTemplates(in []v1.PersistentVolumeClaim) []interface{} {
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
		update["partition"] = *in.RollingUpdate.Partition
		att["rolling_update"] = update
	}
	return []interface{}{att}
}

func expandStatefulsetSpec(in []interface{}) api.StatefulSetSpec {
	if len(in) == 0 || in[0] == nil {
		return api.StatefulSetSpec{}
	}
	spec := api.StatefulSetSpec{}
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
	if v, ok := m["volume_claim_template"]; ok {
		spec.VolumeClaimTemplates = expandVolumeClaimTemplate(v.([]interface{}))
	}
	if v, ok := m["service_name"]; ok {
		spec.ServiceName = v.(string)
	}
	if v, ok := m["pod_management_policy"]; ok {
		spec.PodManagementPolicy = api.PodManagementPolicyType(v.(string))
	}
	if v, ok := m["update_strategy"]; ok {
		spec.UpdateStrategy = expandStatefulSetUpdateStrategy(v.([]interface{}))
	}
	if v, ok := m["revision_history_limit"]; ok {
		limit := int32(v.(int))
		spec.RevisionHistoryLimit = &limit
	}
	return spec
}

func expandPodTemplateSpec(in []interface{}) v1.PodTemplateSpec {
	if len(in) == 0 || in[0] == nil {
		return v1.PodTemplateSpec{}
	}
	spec := v1.PodTemplateSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["metadata"]; ok {
		spec.ObjectMeta = expandMetadata(v.([]interface{}))
	}
	if v, ok := m["spec"]; ok {
		spec.Spec, _ = expandPodSpec(v.([]interface{}))
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
		if v, ok := m["metadata"]; ok {
			claim.ObjectMeta = expandMetadata(v.([]interface{}))
		}
		if v, ok := m["spec"]; ok {
			claim.Spec, _ = expandPersistentVolumeClaimSpec(v.([]interface{}))
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
	if v, ok := m["type"]; ok {
		strategy.Type = api.StatefulSetUpdateStrategyType(v.(string))
	}
	if v, ok := m["rolling_update"]; ok {
		x := v.([]interface{})
		if len(x) > 0 && x[0] == nil {
			m = x[0].(map[string]interface{})
			if v, ok := m["partition"]; ok {
				partition := int32(v.(int))
				strategy.RollingUpdate = &api.RollingUpdateStatefulSetStrategy{
					Partition: &partition,
				}
			}
		}
	}

	return strategy
}