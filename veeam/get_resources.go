package veeam

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/tidwall/gjson"
)

//getJobId ... fetch job id to add vm in job
func getJobID(config Config, jobName string) (string, error) {
	url := "query?type=job&filter=name==" + jobName

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return "", err
	}
	response, err := config.GetResponse(request)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", err
	}
	responseToString := string(response)
	if err != nil {
		log.Fatal(err)
	}

	jobID := gjson.Get(responseToString, "Refs.Refs.0.UID")
	if jobID.String() == "" {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", fmt.Errorf("[Error]  Error:Job Name is incorrect")
	}
	log.Println(jobID.String())

	return jobID.String(), nil
}

//getVmObject ... fetch vm object reference to add vm in job
func getVMObject(config Config, vmName string) (string, error) {
	url := "hierarchyRoots"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return "", err
	}
	response, err := config.GetResponse(request)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", err
	}
	responseToString := string(response)
	if err != nil {
		log.Fatal(err)
	}

	hierarchyRootID := gjson.Get(responseToString, "Refs.0.UID")
	log.Println(hierarchyRootID.String())

	newurl := "lookup?host=" + hierarchyRootID.String() + "&name=" + vmName + "&type=Vm"
	newrequest, err := http.NewRequest("GET", newurl, nil)
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return "", err
	}
	newresponse, err := config.GetResponse(newrequest)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", err
	}
	newresponseToString := string(newresponse)
	if err != nil {
		log.Fatal(err)
	}

	HierarchyObjectRef := gjson.Get(newresponseToString, "HierarchyItems.0.ObjectRef")
	if HierarchyObjectRef.String() == "" {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", fmt.Errorf("[Error]  Error:VM Name is incorrect")
	}
	log.Println(HierarchyObjectRef.String())
	return HierarchyObjectRef.String(), nil
}

//addVMToJob ...  add vm to the job
func addVMToJob(config Config, jobID, vmObjectRef, vmName, vmOrder, vmGpo string) ([]byte, error) {
	url := "jobs/" + jobID + "/includes"
	if vmOrder == "" {
		vmOrder = "0"
	}
	var requestbody CreateObjectInJobSpec
	requestbody.HierarchyObjRef = vmObjectRef
	requestbody.HierarchyObjName = vmName
	requestbody.Order = vmOrder
	requestbody.DisplayName = vmName
	requestbody.GuestProcessingOptions.VssSnapshotOptions.VssSnapshotMode = "RequireSuccess"
	requestbody.GuestProcessingOptions.VssSnapshotOptions.IsCopyOnly = "false"
	body, err := xml.MarshalIndent(&requestbody, "", "")
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return nil, err
	}
	response, err := config.GetResponse(request)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return nil, err
	}

	return response, nil
}

//checkVMExists ... check vm exists in job or not
func checkVMExists(config Config, jobID, vmName string) (string, error) {
	url := "jobs/" + jobID + "/includes"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return "", fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	response, err := config.GetResponse(request)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	responseToString := string(response)
	vmID := gjson.Get(responseToString, "ObjectInJobs.#(Name="+vmName+").ObjectInJobId")
	if vmID.String() == "" {
		return "", fmt.Errorf("Error while fetching a vm : %s", err)
	}
	return vmID.String(), nil
}

//deleteVMFromJob ... delete vm from the job using vm id
func deleteVMFromJob(config Config, jobID, vmID string) (string, error) {
	url := "jobs/" + jobID + "/includes/" + vmID

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("[ERROR] Error in creating http Request %s", err)
		return "", fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	response, err := config.GetResponse(request)
	if err != nil {
		log.Printf("[ERROR] Error in getting response %s", err)
		return "", fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	log.Printf("%s", response)
	return "", nil
}
