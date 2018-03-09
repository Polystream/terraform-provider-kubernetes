
package kubernetes
import (
	"github.com/hashicorp/terraform/helper/schema"
)

func cronJobSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"schedule": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The schedule in Cron format",
		},
		"starting_deadline_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Description:  "ptional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones",
		},
		"concurrency_policy": {
			Type:        schema.TypeString,
			Description: "Specifies how to treat concurrent executions of a Job.",
			Optional:    true,
			Default: 	 "Allow",
		},
		"suspend": {
			Type:        schema.TypeBool,
			Description: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions",
			Optional:    true,
			Default:	 false,
		},
		"successful_jobs_history_limit": {
			Type:        schema.TypeInt,
			Description: "The number of successful finished jobs to retain",
			Optional:    true,
		},
		"failed_jobs_history_limit": {
			Type:        schema.TypeInt,
			Description: "The number of failed finished jobs to retain",
			Optional:    true,
		},
		"job_template": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: jobTemplateSpecFields(),
			},
		},
	}
}

func jobTemplateSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": metadataSchema("CronJob", false),
		"spec": {
			Type:        schema.TypeList,
			Description: "",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: jobSpecFields(),
			},
		},
	}
}

func jobSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parallelism": {
			Type:        schema.TypeInt,
			Description: "Specifies the maximum desired number of pods the job should run at any given time",
			Optional:    true,
		},
		"completions": {
			Type:        schema.TypeInt,
			Description: "Specifies the desired number of successfully finished pods the job should be run with",
			Optional:    true,
		},
		"active_deadline_seconds": {
			Type:        schema.TypeInt,
			Description: "Specifies the duration in seconds relative to the startTime that the job may be active before the system tries to terminate it; value must be positive integer",
			Optional:    true,
		},
		"backoff_limit": {
			Type:        schema.TypeInt,
			Description: "Specifies the number of retries before marking this job failed",
			Optional:    true,
			Default:     6,
		},
		"selector": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: labelSelectorFields(),
			},
		},
		"template": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: statefulsetTemplateFields(),
			},
		},
	}
}