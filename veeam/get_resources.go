package veeam

import (
	"bytes"
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
	requestbody := "<?xml version=\"1.0\" encoding=\"utf-8\"?><CreateObjectInJobSpec xmlns=\"http://www.veeam.com/ent/v1.0\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"><HierarchyObjRef>" + vmObjectRef + "</HierarchyObjRef><HierarchyObjName>" + vmName + "</HierarchyObjName><DisplayName>" + vmName + "</DisplayName> <Order>" + vmOrder + "</Order><GuestProcessingOptions><VssSnapshotOptions><VssSnapshotMode>RequireSuccess</VssSnapshotMode><IsCopyOnly>false</IsCopyOnly>   </VssSnapshotOptions><WindowsGuestFSIndexingOptions><FileSystemIndexingMode>ExceptSpecifiedFolders</FileSystemIndexingMode> <IncludedIndexingFolders/><ExcludedIndexingFolders><Path>%windir%</Path><Path>%ProgramFiles%</Path><Path>%ProgramFiles(x86)%</Path>  <Path>%ProgramW6432%</Path><Path>%TEMP%</Path></ExcludedIndexingFolders></WindowsGuestFSIndexingOptions><LinuxGuestFSIndexingOptions> <FileSystemIndexingMode>ExceptSpecifiedFolders</FileSystemIndexingMode><IncludedIndexingFolders/><ExcludedIndexingFolders><Path>/cdrom</Path> <Path>/dev</Path><Path>/media</Path><Path>/mnt</Path><Path>/proc</Path><Path>/tmp</Path><Path>/lost+found</Path></ExcludedIndexingFolders>    </LinuxGuestFSIndexingOptions><SqlBackupOptions><TransactionLogsProcessing>OnlyOnSuccessJob</TransactionLogsProcessing> <BackupLogsFrequencyMin>15</BackupLogsFrequencyMin><UseDbBackupRetention>true</UseDbBackupRetention><RetainDays>15</RetainDays>               </SqlBackupOptions><WindowsCredentialsId/><LinuxCredentialsId/></GuestProcessingOptions></CreateObjectInJobSpec>"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestbody)))
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
