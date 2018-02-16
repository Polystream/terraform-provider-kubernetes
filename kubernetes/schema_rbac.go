package kubernetes

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func rbacRoleRefSchema(kind string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_group": {
			Type:        schema.TypeString,
			Description: "The API group of the user. Always `rbac.authorization.k8s.io`",
			Required:    true,
			Default:     "rbac.authorization.k8s.io",
		},
		"kind": {
			Type:        schema.TypeString,
			Description: "The kind of resource.",
			Default:     kind,
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the User to bind to.",
			Required:    true,
		},
	}
}

func rbacSubjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_group": {
			Type:        schema.TypeString,
			Description: "The API group of the user. Always `rbac.authorization.k8s.io`",
			Optional:    true,
			Default:     "rbac.authorization.k8s.io",
		},
		"kind": {
			Type:        schema.TypeString,
			Description: "The kind of resource.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the resource to bind to.",
			Required:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "The Namespace of the ServiceAccount",
			Optional:    true,
			Default:     "default",
		},
	}
}

func rbacRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_groups": {
			Type:        schema.TypeSet,
			Description: "APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Optional:    true,
		},
		"resources": {
			Type:        schema.TypeSet,
			Description: "A list of resources this rule applies to.  ResourceAll represents all resources",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Optional:    true,
		},
		"verbs": {
			Type:        schema.TypeSet,
			Description: "A list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.  VerbAll represents all kinds.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Required:    true,
		},
		"resource_names": {
			Type:        schema.TypeSet,
			Description: "An optional white list of names that the rule applies to.  An empty set means that everything is allowed",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Optional:    true,
		},
		"non_resource_urls": {
			Type:        schema.TypeSet,
			Description: "NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
			Optional:    true,
		},
	}
}
