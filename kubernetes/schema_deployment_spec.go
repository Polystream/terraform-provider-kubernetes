
package kubernetes
import (
	"github.com/hashicorp/terraform/helper/schema"
)

func deploymentSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"replicas": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validatePositiveInteger,
			Description:  "Optional the desired number of replicas of the given Template. Value must be a positive integer.",
		},
		"selector": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: labelSelectorFields(),
			},
		},
		"template": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: statefulsetTemplateFields(),
			},
		},
		"deployment_strategy": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: deploymentStrategyFields(),
			},
		},
		"min_ready_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Default: 0,
		},
		"revision_history_limit": {
			Type:     schema.TypeInt,
			Optional: true,
			Default: 10,
		},
		"progress_deadline_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Default: 600,
		},
	}
}

func deploymentStrategyFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"rolling_update": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			Computed: true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"max_unavailable": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_surge": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}