package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func statefulsetSpecFields() map[string]*schema.Schema {
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
		"volume_claim_template": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: statefulsetPersistentVolumeClaimFields(),
			},
		},
		"service_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"pod_management_policy": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"update_strategy": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: statefulsetUpdateSpecFields(),
			},
		},
		"revision_history_limit": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func statefulsetTemplateFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": metadataSchema("statefulset", false),
		"spec": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSpecFields(false),
			},
		},
	}
}

func statefulsetUpdateSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"rolling_update": {
			Type:        schema.TypeList,
			Description: "",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"partition": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}
}

func statefulsetPersistentVolumeClaimFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": metadataSchema("statefulset", false),
		"spec": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"access_modes": {
						Type:        schema.TypeSet,
						Description: "A set of the desired access modes the volume should have. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#access-modes-1",
						Required:    true,
						ForceNew:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Set:         schema.HashString,
					},
					"selector": {
						Type:        schema.TypeList,
						Description: "A label query over volumes to consider for binding.",
						Optional:    true,
						ForceNew:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: labelSelectorFields(),
						},
					},
					"volume_name": {
						Type:        schema.TypeString,
						Description: "The binding reference to the PersistentVolume backing this claim.",
						Optional:    true,
						ForceNew:    true,
						Computed:    true,
					},
					"storage_class_name": {
						Type:        schema.TypeString,
						Description: "Name of the storage class requested by the claim",
						Optional:    true,
						Computed:    true,
						ForceNew:    true,
					},
					"volume_mode": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"resources": {
						Type:        schema.TypeList,
						Description: "A list of the minimum resources the volume should have. More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources",
						Required:    true,
						ForceNew:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: resourcesField(),
						},
					},
				},
			},
		},
	}
}

func labelSelectorFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match_labels": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"match_expressions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"operator": {
						Type:     schema.TypeString,
						Required: true,
					},
					"values": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:  &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
}