package veeam

import "encoding/xml"

//CreateObjectInJobSpec ...
type CreateObjectInJobSpec struct {
	XMLName                xml.Name `xml:"CreateObjectInJobSpec"`
	Text                   string   `xml:",chardata"`
	Xmlns                  string   `xml:"xmlns,attr"`
	Xsd                    string   `xml:"xsd,attr"`
	Xsi                    string   `xml:"xsi,attr"`
	HierarchyObjRef        string   `xml:"HierarchyObjRef"`
	HierarchyObjName       string   `xml:"HierarchyObjName"`
	DisplayName            string   `xml:"DisplayName"`
	Order                  string   `xml:"Order"`
	GuestProcessingOptions struct {
		Text               string `xml:",chardata"`
		VssSnapshotOptions struct {
			Text            string `xml:",chardata"`
			VssSnapshotMode string `xml:"VssSnapshotMode"`
			IsCopyOnly      string `xml:"IsCopyOnly"`
		} `xml:"VssSnapshotOptions"`
	} `xml:"GuestProcessingOptions"`
}
