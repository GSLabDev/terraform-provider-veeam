package veeam

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVeeamJobVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVeeamJobVMCreate,
		Read:   resourceVeeamJobVMRead,
		Delete: resourceVeeamJobVMDelete,
		Schema: map[string]*schema.Schema{
			"job_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_order": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vm_gpo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

//resourceVeeamJobVMCreate ... resource- add vm to backup job
func resourceVeeamJobVMCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(Config)
	jobName := d.Get("job_name").(string)
	vmName := d.Get("vm_name").(string)
	vmOrder := d.Get("vm_order").(string)
	vmGpo := d.Get("vm_gpo").(string)

	//fetch job ID
	jobID, err := getJobID(config, jobName)
	if err != nil {
		log.Printf("[ERROR] Error in getting JOB ID %s", err)
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	//fetch vmObjectReference using vmname
	vmObjectRef, err := getVMObject(config, vmName)
	if err != nil {
		log.Printf("[ERROR] Error in getting VM Object Reference %s", err)
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	//add vm to backup job
	addVM, err := addVMToJob(config, jobID, vmObjectRef, vmName, vmOrder, vmGpo)
	log.Println(addVM)
	if err != nil {
		log.Printf("[Error] Error in adding vm to backup job")
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	d.SetId(jobID + "_" + vmName)
	return nil
}

// resourceVeeamJobVMRead ... resource read
func resourceVeeamJobVMRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)
	jobName := d.Get("job_name").(string)
	vmName := d.Get("vm_name").(string)

	//fetch job ID using job_name
	jobID, err := getJobID(config, jobName)
	if err != nil {
		log.Printf("[Error] Error in fetching job ID")
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	//check vm is exists in job or not
	_, error := checkVMExists(config, jobID, vmName)
	if error != nil {
		log.Printf("[Error] Error in reading Virtual Machine : %s ", error)
		d.SetId("")
	}

	return nil

}

//resourceVeeamJobVMDelete ... resource delete
func resourceVeeamJobVMDelete(d *schema.ResourceData, meta interface{}) error {
	resourceVeeamJobVMRead(d, meta)
	if d.Id() == "" {
		log.Println("[ERROR] Cannot find vm in the specified job")
		return fmt.Errorf("[ERROR] Cannot find vm in the specified job")
	}
	config := meta.(Config)
	jobName := d.Get("job_name").(string)
	vmName := d.Get("vm_name").(string)
	jobID, err := getJobID(config, jobName)
	if err != nil {
		log.Printf("[Error] Error in fetching job id")
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}

	//fetch vm id to delete vm from job
	vmID, err := checkVMExists(config, jobID, vmName)
	fmt.Println(vmID)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	//delete vm from the backup job
	response, err := deleteVMFromJob(config, jobID, vmID)
	fmt.Println(response)
	if err != nil {
		return fmt.Errorf("[Error]  Error: %s", err.Error())
	}

	return nil

}
