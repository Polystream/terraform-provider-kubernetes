
package kubernetes
import (
	"github.com/hashicorp/terraform/helper/schema"
)

func daemonsetSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"update_strategy": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: statefulsetUpdateSpecFields(),
			},
		},
		"min_ready_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: false,
		},
		"revision_history_limit": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}