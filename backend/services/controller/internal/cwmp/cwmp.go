package cwmp

import (
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
A successful response to SetParameterValues method returns an integer enumeration defined as follows:

0 = All Parameter changes have been validated and applied.
1 = All Parameter changes have been validated and committed, but some or all are not yet applied (for example, if a reboot is required before the new values are applied).
*/
type SetParameterValuesResponseStatus int

const (
	ALL_OK       = 0
	SOME_PENDING = 1
)

/* -------------------------------------------------------------------------- */

const WRITABLE = "1"

func ParamTypeIsWritable(param string) bool {
	return param == WRITABLE
}

type SoapEnvelope struct {
	XMLName xml.Name
	Header  SoapHeader
	Body    SoapBody
}

type SoapHeader struct {
	Id string `xml:"ID"`
}
type SoapBody struct {
	CWMPMessage CWMPMessage `xml:",any"`
}

type CWMPMessage struct {
	XMLName xml.Name
}

type EventStruct struct {
	EventCode  string
	CommandKey string
}

type ParameterValueStruct struct {
	Name  string
	Value string
}

type ParameterInfoStruct struct {
	Name     string
	Writable string
}

type ParameterAttributeStruct struct {
	Name         string
	Notification int
	AccessList   []string
}

type SetParameterValues_ struct {
	ParameterList []ParameterValueStruct `xml:"Body>SetParameterValues>ParameterList>ParameterValueStruct"`
	ParameterKey  string                 `xml:"Body>SetParameterValues>ParameterKey>string"`
}

type GetParameterValues_ struct {
	ParameterNames []string `xml:"Body>GetParameterValues>ParameterNames>string"`
}

type GetParameterNames_ struct {
	ParameterPath []string `xml:"Body>GetParameterNames>ParameterPath"`
	NextLevel     string   `xml:"Body>GetParameterNames>NextLevel"`
}

type GetParameterAttributes_ struct {
	ParameterNames []string `xml:"Body>GetParameterAttributes>ParameterNames>string"`
}

type GetParameterValuesResponse struct {
	ParameterList []ParameterValueStruct `xml:"Body>GetParameterValuesResponse>ParameterList>ParameterValueStruct"`
}

type GetParameterNamesResponse struct {
	ParameterList []ParameterInfoStruct `xml:"Body>GetParameterNamesResponse>ParameterList>ParameterInfoStruct"`
}

type GetParameterAttributesResponse struct {
	ParameterList []ParameterAttributeStruct `xml:"Body>GetParameterAttributesResponse>ParameterList>ParameterAttributeStruct"`
}

type CWMPInform struct {
	DeviceId      DeviceID               `xml:"Body>Inform>DeviceId"`
	Events        []EventStruct          `xml:"Body>Inform>Event>EventStruct"`
	ParameterList []ParameterValueStruct `xml:"Body>Inform>ParameterList>ParameterValueStruct"`
}

func (s *SoapEnvelope) KindOf() string {
	return s.Body.CWMPMessage.XMLName.Local
}

func (i *CWMPInform) GetEvents() string {
	res := ""
	for idx := range i.Events {
		res += i.Events[idx].EventCode
	}

	return res
}

func (i *CWMPInform) GetConnectionRequest() string {
	for idx := range i.ParameterList {
		// valid condition for both tr98 and tr181
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.ManagementServer.ConnectionRequestURL") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetSoftwareVersion() string {
	for idx := range i.ParameterList {
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.DeviceInfo.SoftwareVersion") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetHardwareVersion() string {
	for idx := range i.ParameterList {
		if strings.HasSuffix(i.ParameterList[idx].Name, "Device.DeviceInfo.HardwareVersion") {
			return i.ParameterList[idx].Value
		}
	}

	return ""
}

func (i *CWMPInform) GetDataModelType() string {
	if strings.HasPrefix(i.ParameterList[0].Name, "InternetGatewayDevice") {
		return "TR098"
	} else if strings.HasPrefix(i.ParameterList[0].Name, "Device") {
		return "TR181"
	}

	return ""
}

type DeviceID struct {
	Manufacturer string
	OUI          string
	SerialNumber string
	ProductClass string
}

func InformResponse(mustUnderstand string) string {
	mustUnderstandHeader := ""
	if mustUnderstand != "" {
		mustUnderstandHeader = `<cwmp:ID soap:mustUnderstand="1">` + mustUnderstand + `</cwmp:ID>`
	}

	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header>` + mustUnderstandHeader + `</soap:Header>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:InformResponse>
      <MaxEnvelopes>1</MaxEnvelopes>
    </cwmp:InformResponse>
  </soap:Body>
</soap:Envelope>`
}

func GetParameterValues(leaf string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterValues>
      <ParameterNames>
      	<string>` + leaf + `</string>
      </ParameterNames>
    </cwmp:GetParameterValues>
  </soap:Body>
</soap:Envelope>`
}

func GetParameterMultiValues(leaves []string) string {
	msg := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterValues>
      <ParameterNames>`

	for idx := range leaves {
		msg += `<string>` + leaves[idx] + `</string>`

	}
	msg += `</ParameterNames>
    </cwmp:GetParameterValues>
  </soap:Body>
</soap:Envelope>`
	return msg
}

type SetParameterValuesResponse struct {
	Status int `xml:"Body>SetParameterValuesResponse>Status"`
}

func SetParameterValues(leaf string, value string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:SetParameterValues>
      <ParameterList soapenc:arrayType="cwmp:ParameterValueStruct[1]">
		  <ParameterValueStruct>
			  <Name>` + leaf + `</Name>
			  <Value>` + value + `</Value>
		  </ParameterValueStruct>
      </ParameterList>
      <ParameterKey>LC1309` + randToken() + `</ParameterKey>
    </cwmp:SetParameterValues>
  </soap:Body>
</soap:Envelope>`
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func SetParameterMultiValues(data map[string]string) string {
	msg := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:SetParameterValues>
      <ParameterList soapenc:arrayType="cwmp:ParameterValueStruct[` + fmt.Sprint(len(data)) + `]">`

	for key, value := range data {
		msg += `<ParameterValueStruct>
			  <Name>` + key + `</Name>
			  <Value>` + value + `</Value>
		  </ParameterValueStruct>`
	}

	msg += `</ParameterList>
      <ParameterKey>LC1309` + randToken() + `</ParameterKey>
    </cwmp:SetParameterValues>
  </soap:Body>
</soap:Envelope>`

	return msg
}

func GetParameterNames(leaf string, nextlevel int) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:GetParameterNames>
      <ParameterPath>` + leaf + `</ParameterPath>
      <NextLevel>` + strconv.Itoa(nextlevel) + `</NextLevel>
    </cwmp:GetParameterNames>
  </soap:Body>
</soap:Envelope>`
}

func FactoryReset() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:FactoryReset/>
  </soap:Body>
</soap:Envelope>`
}

func Download(filetype, url, username, password, filesize string) string {
	// 3 Vendor Configuration File
	// 1 Firmware Upgrade Image

	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:Download>
      <CommandKey>MSDWK</CommandKey>
      <FileType>` + filetype + `</FileType>
      <URL>` + url + `</URL>
      <Username>` + username + `</Username>
      <Password>` + password + `</Password>
      <FileSize>` + filesize + `</FileSize>
      <TargetFileName></TargetFileName>
      <DelaySeconds>0</DelaySeconds>
      <SuccessURL></SuccessURL>
      <FailureURL></FailureURL>
    </cwmp:Download>
  </soap:Body>
</soap:Envelope>`
}

func CancelTransfer() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:CancelTransfer>
      <CommandKey></CommandKey>
    <cwmp:CancelTransfer/>
  </soap:Body>
</soap:Envelope>`
}

type TimeWindowStruct struct {
	WindowStart string
	WindowEnd   string
	WindowMode  string
	UserMessage string
	MaxRetries  string
}

func (window *TimeWindowStruct) String() string {
	return `<TimeWindowStruct>
<WindowStart>` + window.WindowStart + `</WindowStart>
<WindowEnd>` + window.WindowEnd + `</WindowEnd>
<WindowMode>` + window.WindowMode + `</WindowMode>
<UserMessage>` + window.UserMessage + `</UserMessage>
<MaxRetries>` + window.MaxRetries + `</MaxRetries>
</TimeWindowStruct>`
}

func ScheduleDownload(filetype, url, username, password, filesize string, windowslist []fmt.Stringer) string {
	ret := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Header/>
  <soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <cwmp:ScheduleDownload>
      <CommandKey>MSDWK</CommandKey>
      <FileType>` + filetype + `</FileType>
      <URL>` + url + `</URL>
      <Username>` + username + `</Username>
      <Password>` + password + `</Password>
      <FileSize>` + filesize + `</FileSize>
      <TargetFileName></TargetFileName>
      <TimeWindowList>`

	for _, op := range windowslist {
		ret += op.String()
	}

	ret += `</TimeWindowList>
    </cwmp:ScheduleDownload>
  </soap:Body>
</soap:Envelope>`

	return ret
}

type InstallOpStruct struct {
	Url                  string
	Uuid                 string
	Username             string
	Password             string
	ExecutionEnvironment string
}

func (op *InstallOpStruct) String() string {
	return `<InstallOpStruct>
	<URL>` + op.Url + `</URL>
	<UUID>` + op.Uuid + `</UUID>
	<Username>` + op.Username + `</Username>
	<Password>` + op.Password + `</Password>
	<ExecutionEnvRef>` + op.ExecutionEnvironment + `</ExecutionEnvRef>
</InstallOpStruct>`
}

type UpdateOpStruct struct {
	Uuid     string
	Version  string
	Url      string
	Username string
	Password string
}

func (op *UpdateOpStruct) String() string {
	return `<UpdateOpStruct>
<UUID>` + op.Uuid + `</UUID>
<Version>` + op.Version + `</Version>
<URL>` + op.Url + `</URL>
<Username>` + op.Username + `</Username>
<Password>` + op.Password + `</Password>
</UpdateOpStruct>`
}

type UninstallOpStruct struct {
	Uuid                 string
	Version              string
	ExecutionEnvironment string
}

func (op *UninstallOpStruct) String() string {
	return `<UninstallOpStruct>
<UUID>` + op.Uuid + `</UUID>
<Version>` + op.Version + `</Version>
<ExecutionEnvRef>` + op.ExecutionEnvironment + `</ExecutionEnvRef>
</UninstallOpStruct>`
}

func ChangeDuState(ops []fmt.Stringer) string {
	ret := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:schemaLocation="urn:dslforum-org:cwmp-1-0 ..\schemas\wt121.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soap:Header/>
<soap:Body soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<cmwp:ChangeDUState>
<Operations>`

	for _, op := range ops {
		ret += op.String()
	}

	ret += `</Operations>
<CommandKey></CommandKey>
</cmwp:ChangeDUState>
</soap:Body>
</soap:Envelope>`

	return ret
}

// CPE side

func Inform(serial string) string {
	return `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0"><soap:Header><cwmp:ID soap:mustUnderstand="1">5058</cwmp:ID></soap:Header>
	<soap:Body><cwmp:Inform><DeviceId><Manufacturer>ADB Broadband</Manufacturer>
<OUI>0013C8</OUI>
<ProductClass>VV5522</ProductClass>
<SerialNumber>` + serial + `</SerialNumber>
</DeviceId>
<Event soap-enc:arrayType="cwmp:EventStruct[1]">
<EventStruct><EventCode>6 CONNECTION REQUEST</EventCode>
<CommandKey></CommandKey>
</EventStruct>
</Event>
<MaxEnvelopes>1</MaxEnvelopes>
<CurrentTime>` + time.Now().Format(time.RFC3339) + `</CurrentTime>
<RetryCount>0</RetryCount>
<ParameterList soap-enc:arrayType="cwmp:ParameterValueStruct[8]">
<ParameterValueStruct><Name>InternetGatewayDevice.ManagementServer.ConnectionRequestURL</Name>
<Value xsi:type="xsd:string">http://104.199.175.27:7547/ConnectionRequest-` + serial + `</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.ManagementServer.ParameterKey</Name>
<Value xsi:type="xsd:string"></Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceSummary</Name>
<Value xsi:type="xsd:string">InternetGatewayDevice:1.2[](Baseline:1,EthernetLAN:1,WiFiLAN:1,ADSLWAN:1,EthernetWAN:1,QoS:1,QoSDynamicFlow:1,Bridging:1,Time:1,IPPing:1,TraceRoute:1,DeviceAssociation:1,UDPConnReq:1),VoiceService:1.0[1](TAEndpoint:1,SIPEndpoint:1)</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.HardwareVersion</Name>
<Value xsi:type="xsd:string">` + serial + `</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.ProvisioningCode</Name>
<Value xsi:type="xsd:string">ABCD</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SoftwareVersion</Name>
<Value xsi:type="xsd:string">4.0.8.17785</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.DeviceInfo.SpecVersion</Name>
<Value xsi:type="xsd:string">1.0</Value>
</ParameterValueStruct>
<ParameterValueStruct><Name>InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.ExternalIPAddress</Name>
<Value xsi:type="xsd:string">12.0.0.10</Value>
</ParameterValueStruct>
</ParameterList>
</cwmp:Inform>
</soap:Body></soap:Envelope>`
}

/*
func BuildGetParameterValuesResponse(serial string, leaves GetParameterValues_) string {
	ret := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	<soap:Header><cwmp:ID soap:mustUnderstand="1">3</cwmp:ID></soap:Header>
	<soap:Body><cwmp:GetParameterValuesResponse>`

	db, _ := sqlite3.Open("/tmp/cpe.db")

	n_leaves := 0
	var temp string
	for _, leaf := range leaves.ParameterNames {
		sql := "select key, value, tipo from params where key like '" + leaf + "%'"
		for s, err := db.Query(sql); err == nil; err = s.Next() {
			n_leaves++
			var key string
			var value string
			var tipo string
			s.Scan(&key, &value, &tipo)
			temp += `<ParameterValueStruct>
			<Name>` + key + `</Name>
			<Value xsi:type="` + tipo + `">` + value + `</Value>
			</ParameterValueStruct>`
		}
	}

	ret += `<ParameterList soap-enc:arrayType="cwmp:ParameterValueStruct[` + strconv.Itoa(n_leaves) + `]">`
	ret += temp
	ret += `</ParameterList></cwmp:GetParameterValuesResponse></soap:Body></soap:Envelope>`

	return ret
}

func BuildGetParameterNamesResponse(serial string, leaves GetParameterNames_) string {
	ret := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:cwmp="urn:dslforum-org:cwmp-1-0">
	<soap:Header><cwmp:ID soap:mustUnderstand="1">69</cwmp:ID></soap:Header>
	<soap:Body><cwmp:GetParameterNamesResponse>`
	db, _ := sqlite3.Open("/tmp/cpe.db")

	obj := make(map[string]bool)
	var temp string
	for _, leaf := range leaves.ParameterPath {
		fmt.Println(leaf)
		sql := "select key, value, tipo from params where key like '" + leaf + "%'"
		for s, err := db.Query(sql); err == nil; err = s.Next() {
			var key string
			var value string
			var tipo string
			s.Scan(&key, &value, &tipo)
			var sp = strings.Split(strings.Split(key, leaf)[1], ".")
			nextlevel, _ := strconv.Atoi(leaves.NextLevel)
			if nextlevel == 0 {
				root := leaf
				obj[root] = true
				for idx := range sp {
					if idx == len(sp)-1 {
						root = root + sp[idx]
					} else {
						root = root + sp[idx] + "."
					}
					obj[root] = true
				}
			} else {
				if !obj[sp[0]] {
					if len(sp) > 1 {
						obj[leaf+sp[0]+"."] = true
					} else {
						obj[leaf+sp[0]] = true
					}

				}
			}

		}
	}

	for o := range obj {
		temp += `<ParameterInfoStruct>
				<Name>` + o + `</Name>
				<Writable>true</Writable>
				</ParameterInfoStruct>`
	}

	fmt.Println(len(obj))
	ret += `<ParameterList soap-enc:arrayType="cwmp:ParameterInfoStruct[]">`
	ret += temp
	ret += `</ParameterList></cwmp:GetParameterNamesResponse></soap:Body></soap:Envelope>`

	return ret
}
*/
