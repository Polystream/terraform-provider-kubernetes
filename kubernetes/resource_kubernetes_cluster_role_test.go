package kubernetes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/api/rbac/v1"
	kubernetes "k8s.io/client-go/kubernetes"
)

func TestAccKubernetesClusterRole(t *testing.T) {
	var conf api.ClusterRole
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_cluster_role.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesClusterRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRoleConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesClusterRoleExists("kubernetes_cluster_role.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.api_group", "rbac.authorization.k8s.io"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.kind", "ClusterRole"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.name", "cluster-admin"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.api_group", "rbac.authorization.k8s.io"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.name", "notauser"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.kind", "User"),
				),
			},
			{
				Config: testAccKubernetesClusterRoleConfig_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesClusterRoleExists("kubernetes_cluster_role.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_cluster_role.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.api_group", "rbac.authorization.k8s.io"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.kind", "ClusterRole"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "role_ref.name", "cluster-admin"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.#", "3"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.api_group", "rbac.authorization.k8s.io"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.name", "notauser"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.0.kind", "User"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.1.namespace", "kube-system"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.1.name", "default"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.1.kind", "ServiceAccount"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.2.api_group", "rbac.authorization.k8s.io"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.2.name", "system:masters"),
					resource.TestCheckResourceAttr("kubernetes_cluster_role.test", "subject.2.kind", "Group"),
				),
			},
		},
	})
}

func TestAccKubernetesClusterRole_importBasic(t *testing.T) {
	resourceName := "kubernetes_cluster_role.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRoleConfig_basic(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKubernetesClusterRoleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*kubernetes.Clientset)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_cluster_role" {
			continue
		}
		name := rs.Primary.ID
		resp, err := conn.Rbac().ClusterRoles().Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("ClusterRole still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesClusterRoleExists(n string, obj *api.ClusterRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*kubernetes.Clientset)
		name := rs.Primary.ID
		resp, err := conn.Rbac().ClusterRoles().Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *resp
		return nil
	}
}

func testAccKubernetesClusterRoleConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_cluster_role" "test" {
	metadata {
		name = "%s"
	}
	rule {
		api_groups = [""]
		resources = ["pods"]
		verbs = ["get", "watch", "list"]
	}
	rule {
		api_groups = [""]
		resources = ["secrets"]
		verbs = ["get", "watch", "list"]
	}
}`, name)
}

func testAccKubernetesClusterRoleConfig_modified(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_cluster_role" "test" {
	metadata {
		name = "%s"
	}
	rule {
		api_groups = ["apps"]
		resources = ["secrets"]
		verbs = ["get", "watch", "list"]
	}
}`, name)
}
