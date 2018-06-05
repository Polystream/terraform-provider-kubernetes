package kubernetes

import (
	api "k8s.io/api/apps/v1"
)

func flattenDaemonsetSpec(in api.DaemonSetSpec) []interface{} {
	att := make(map[string]interface{})
	
	att["selector"] = flattenLabelSelector(in.Selector)

	att["template"] = flattenPodTemplateSpec(in.Template)
	att["update_strategy"] = flattenDaemonsetUpdateStrategy(in.UpdateStrategy)
	if in.RevisionHistoryLimit != nil {
		att["revision_history_limit"] = *in.RevisionHistoryLimit
	}

	att["min_ready_seconds"] = in.MinReadySeconds

	return []interface{}{att}
}

func flattenDaemonsetUpdateStrategy(in api.DaemonSetUpdateStrategy) []interface{} {
	att := make(map[string]interface{})
	att["type"] = string(in.Type)
	if(in.RollingUpdate != nil){
		update := make(map[string]interface{})
		att["rolling_update"] = []interface{}{update}
	}
	return []interface{}{att}
}

func expandDaemonsetSpec(in []interface{}) api.DaemonSetSpec {
	if len(in) == 0 || in[0] == nil {
		return api.DaemonSetSpec{}
	}
	spec := api.DaemonSetSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["selector"].([]interface{}); ok {
		spec.Selector = expandLabelSelector(v)
	}
	if v, ok := m["template"].([]interface{}); ok {
		spec.Template = expandPodTemplateSpec(v)
	}
	if v, ok := m["update_strategy"].([]interface{}); ok {
		spec.UpdateStrategy = expandDaemonSetUpdateStrategy(v)
	}
	if v, ok := m["revision_history_limit"].(int); ok {
		spec.RevisionHistoryLimit = ptrToInt32(int32(v))
	}
	if v, ok := m["min_ready_seconds"].(int); ok {
		spec.MinReadySeconds = int32(v)
	}
	return spec
}

func expandDaemonSetUpdateStrategy(in []interface{}) api.DaemonSetUpdateStrategy {
	if len(in) == 0 || in[0] == nil {
		return api.DaemonSetUpdateStrategy{}
	}
	strategy := api.DaemonSetUpdateStrategy{}
	m := in[0].(map[string]interface{})
	if v, ok := m["type"].(string); ok {
		strategy.Type = api.DaemonSetUpdateStrategyType(v)
	}
	if v, ok := m["rolling_update"]; ok {
		x := v.([]interface{})
		if len(x) > 0 && x[0] == nil {
			m = x[0].(map[string]interface{})
			strategy.RollingUpdate = &api.RollingUpdateDaemonSet{}
		}
	}

	return strategy
}