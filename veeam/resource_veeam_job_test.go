package veeam

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testBasicPreCheckVirtualMachine(t *testing.T) {
	if v := os.Getenv("VEEAM_JOB_NAME"); v == "" {
		t.Fatal("VEEAM_BACKUP_JOB must be set for acceptance tests")
	}
	if v := os.Getenv("VEEAM_VM_NAME"); v == "" {
		t.Fatal("VEEAM_VIRTUAL_MACHINE must be set for acceptance tests")
	}

}

func TestAccAddVMToJob_Basic(t *testing.T) {
	resourceName := "veeam_job_vm.AddVM"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testBasicPreCheckVirtualMachine(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy(resourceName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVirtualMachineConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMachineExists(resourceName),
				),
			},
		},
	})
}
func testAccCheckVirtualMachineDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		config := testAccProvider.Meta().(Config)
		jobName := rs.Primary.Attributes["job_name"]
		vmName := rs.Primary.Attributes["vm_name"]
		jobID, err := getJobID(config, jobName)
		vmID, err := checkVMExists(config, jobID, vmName)
		response, err := deleteVMFromJob(config, jobID, vmID)
		log.Printf("%s", response)
		if err != nil {
			return fmt.Errorf("Virtual Machine still exists: %v", err)
		}
		return nil
	}
}

func testAccCheckVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No request ID is set")
		}

		config := testAccProvider.Meta().(Config)
		//splits the jobID and VM name from the ID
		jobVM := strings.SplitN(rs.Primary.ID, "_", 2)
		vmID, err := checkVMExists(config, jobVM[0], jobVM[1])
		if err != nil {
			log.Printf("[ERROR]error %s", err.Error())
		}
		if vmID == "" {
			return fmt.Errorf("vm does not exists")
		}
		return nil
	}
}

func testAccCheckVirtualMachineConfigBasic() string {
	return fmt.Sprintf(`

	provider "veeam" {
		server_ip  = "%s"
		port = "%s"
		username = "%s"
		password="%s"
	  }
resource "veeam_job_vm" "AddVM" {
  job_name = "%s"
  vm_name = "%s"
}`, os.Getenv("VEEAM_SERVER_IP"),
		os.Getenv("VEEAM_SERVER_PORT"),
		os.Getenv("VEEAM_SERVER_USERNAME"),
		os.Getenv("VEEAM_SERVER_PASSWORD"),
		os.Getenv("VEEAM_JOB_NAME"),
		os.Getenv("VEEAM_VM_NAME"))
}
