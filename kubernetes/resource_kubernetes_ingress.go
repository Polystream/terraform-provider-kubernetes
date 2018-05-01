package kubernetes

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
	api "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
)

func resourceKubernetesIngress() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesIngressCreate,
		Read:   resourceKubernetesIngressRead,
		Exists: resourceKubernetesIngressExists,
		Update: resourceKubernetesIngressUpdate,
		Delete: resourceKubernetesIngressDelete,

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("Ingress", true),
			"spec": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew: true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew: true,
							MaxItems:    1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
									"service_name": {
										Type:        schema.TypeString,
										Required:    true,
									},
									"service_port": {
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
						},
						"tls": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hosts": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"secret_name": {
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
								"host": {
									Type:        schema.TypeString,
									Optional:    true,
								},
								"http": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"paths": {
												Type:        schema.TypeList,
												Required:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"path": {
															Type:        schema.TypeString,
															Required:    true,
														},
														"backend": {
															Type:        schema.TypeList,
															Required:    true,
															MaxItems:    1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"service_name": {
																		Type:        schema.TypeString,
																		Required:    true,
																	},
																	"service_port": {
																		Type:        schema.TypeInt,
																		Required:    true,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				},
			},
		},
	}
}

func resourceKubernetesIngressCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	ingress := api.Ingress{
		ObjectMeta: metadata,
		Spec:       expandIngressSpec(d.Get("spec").([]interface{})),
	}
	log.Printf("[INFO] Creating new ingress: %#v", ingress)
	out, err := conn.ExtensionsV1beta1().Ingresses(metadata.Namespace).Create(&ingress)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new ingress: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesIngressRead(d, meta)
}

func resourceKubernetesIngressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading ingress %s", name)
	ingress, err := conn.ExtensionsV1beta1().Ingresses(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received ingress: %#v", ingress)
	err = d.Set("metadata", flattenMetadata(ingress.ObjectMeta))
	if err != nil {
		return err
	}

	flattened := flattenIngressSpec(ingress.Spec)
	log.Printf("[DEBUG] Flattened ingress spec: %#v", flattened)
	err = d.Set("spec", flattened)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesIngressUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	ingress := api.Ingress{
		ObjectMeta: metadata,
		Spec:       expandIngressSpec(d.Get("spec").([]interface{})),
	}

	data, err := json.Marshal(ingress)
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating Ingress %s: %s", d.Id(), ingress)

	out, err := conn.ExtensionsV1beta1().Ingresses(namespace).Patch(name, pkgApi.StrategicMergePatchType, data)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated Ingress: %#v", out)

	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesIngressRead(d, meta)
}

func resourceKubernetesIngressDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting ingress: %#v", name)
	err = conn.ExtensionsV1beta1().Ingresses(namespace).Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Ingress %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubernetesIngressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking ingress %s", name)
	_, err = conn.ExtensionsV1beta1().Ingresses(namespace).Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func expandIngressSpec(in []interface{}) api.IngressSpec {
	if len(in) == 0 || in[0] == nil {
		return api.IngressSpec{}
	}
	spec := api.IngressSpec{}
	m := in[0].(map[string]interface{})
	if v, ok := m["backend"].([]interface{}); ok {
		if len(v) > 0 {
			n := v[0].(map[string]interface{})
			backend := api.IngressBackend{}
			backend.ServiceName = n["service_name"].(string)
			backend.ServicePort = expandIntOrString(n["service_port"].(int))
			spec.Backend = &backend
		}
	}

	if v, ok := m["tls"].([]interface{}); ok {
		obj := make([]api.IngressTLS, len(v), len(v))
		for i, n := range v {
			tls := n.(map[string]interface{})
			obj[i] = api.IngressTLS{
				Hosts: sliceOfString(tls["hosts"].([]interface {})),
				SecretName: tls["secret_name"].(string),
			}
		}
		spec.TLS = obj
	}

	if v, ok := m["rules"].([]interface{}); ok {
		obj := make([]api.IngressRule, len(v), len(v))
		for i, n := range v {
			rule := n.(map[string]interface{})
			obj[i] = api.IngressRule{
				Host: rule["host"].(string),
			}

			value := api.IngressRuleValue{}

			if w, ok := rule["http"].([]interface{}); ok {
				if len(w) > 0 {
					http := api.HTTPIngressRuleValue{}

					x := w[0].(map[string]interface{})["paths"].([]interface{});
					paths := make([]api.HTTPIngressPath, len(x), len(x))
					for i, o := range x {
						path := o.(map[string]interface{})
						paths[i] = api.HTTPIngressPath{
							Path: path["path"].(string),
						}

						backend := (path["backend"].([]interface{}))[0].(map[string]interface{})
						paths[i].Backend = api.IngressBackend {
							ServiceName: backend["service_name"].(string),
							ServicePort: expandIntOrString(backend["service_port"].(int)),
						}
					}
					http.Paths = paths
					value.HTTP = &http
				}
			}
			obj[i].IngressRuleValue = value
		}
		spec.Rules = obj
	}
	return spec
}

func flattenIngressSpec(in api.IngressSpec) []interface{} {
	att := make(map[string]interface{})
	if in.Backend != nil {
		backend := *in.Backend
		a := make(map[string]interface{})
		a["service_name"] = backend.ServiceName
		a["service_port"] = flattenIntOrString(backend.ServicePort)
		att["backend"] = a
	}
	
	if in.TLS != nil {
		obj := make([]map[string]interface{}, len(in.TLS), len(in.TLS))
		for i, n := range in.TLS {
			a := make(map[string]interface{})
			a["hosts"] = n.Hosts
			a["secret_name"] = n.SecretName
			obj[i] = a
		}
		att["tls"] = obj
	}

	if in.Rules != nil {
		obj := make([]map[string]interface{}, len(in.Rules), len(in.Rules))
		for i, n := range in.Rules {
			a := make(map[string]interface{})
			a["host"] = n.Host
			http := *n.HTTP
			p := make([]map[string]interface{}, len(http.Paths), len(http.Paths))
			for j, o := range http.Paths {
				b := make(map[string]interface{})
				b["path"] = o.Path
				c := make(map[string]interface{})
				c["service_name"] = o.Backend.ServiceName
				c["service_port"] = flattenIntOrString(o.Backend.ServicePort)
				backend := make([]map[string]interface{}, 1, 1)
				backend[0] = c
				b["backend"] = backend
				p[j] = b
			}

			h := make([]map[string]interface{}, 1, 1)
			tmp := make(map[string]interface{})
			tmp["paths"] = p

			h[0] = tmp
			a["http"] = h
			obj[i] = a
		}
		att["rules"] = obj
	}

	return []interface{}{att}
}