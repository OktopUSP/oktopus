package api

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/cwmp"
	"github.com/leandrofars/oktopus/internal/entity"
	"github.com/leandrofars/oktopus/internal/utils"
)

type ParamData struct {
	Writable bool        `json:"writable"`
	Value    interface{} `json:"value"`
}

type WiFi struct {
	Path                 string      `json:"path"`
	Name                 ParamData   `json:"name"`
	SSID                 ParamData   `json:"ssid"`
	Password             ParamData   `json:"password"`
	Security             ParamData   `json:"security"`
	SecurityCapabilities []ParamData `json:"securityCapabilities"`
	Standard             ParamData   `json:"standard"`
	Enable               ParamData   `json:"enable"`
	Status               ParamData   `json:"status"`
}

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/mux"

// 	"github.com/leandrofars/oktopus/internal/db"
// 	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
// 	"github.com/leandrofars/oktopus/internal/utils"
// 	"google.golang.org/protobuf/proto"
// )

// type WiFi struct {
// 	SSID                 string   `json:"ssid"`
// 	Password             string   `json:"password"`
// 	Security             string   `json:"security"`
// 	SecurityCapabilities []string `json:"securityCapabilities"`
// 	AutoChannelEnable    bool     `json:"autoChannelEnable"`
// 	Channel              int      `json:"channel"`
// 	ChannelBandwidth     string   `json:"channelBandwidth"`
// 	FrequencyBand        string   `json:"frequencyBand"`
// 	//PossibleChannels     		[]int    `json:"PossibleChannels"`
// 	SupportedChannelBandwidths []string `json:"supportedChannelBandwidths"`
// }

// func (a *Api) deviceWifi(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	sn := vars["sn"]
// 	device := a.deviceExists(sn, w)

// 	if r.Method == http.MethodGet {
// 		msg := utils.NewGetMsg(usp_msg.Get{
// 			ParamPaths: []string{
// 				"Device.WiFi.SSID.[Enable==true].SSID",
// 				//"Device.WiFi.AccessPoint.[Enable==true].SSIDReference",
// 				"Device.WiFi.AccessPoint.[Enable==true].Security.ModeEnabled",
// 				"Device.WiFi.AccessPoint.[Enable==true].Security.ModesSupported",
// 				//"Device.WiFi.EndPoint.[Enable==true].",
// 				"Device.WiFi.Radio.[Enable==true].AutoChannelEnable",
// 				"Device.WiFi.Radio.[Enable==true].Channel",
// 				"Device.WiFi.Radio.[Enable==true].CurrentOperatingChannelBandwidth",
// 				"Device.WiFi.Radio.[Enable==true].OperatingFrequencyBand",
// 				//"Device.WiFi.Radio.[Enable==true].PossibleChannels",
// 				"Device.WiFi.Radio.[Enable==true].SupportedOperatingChannelBandwidths",
// 			},
// 			MaxDepth: 2,
// 		})

// 		encodedMsg, err := proto.Marshal(&msg)
// 		if err != nil {
// 			log.Println(err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		record := utils.NewUspRecord(encodedMsg, sn)
// 		tr369Message, err := proto.Marshal(&record)
// 		if err != nil {
// 			log.Fatalln("Failed to encode tr369 record:", err)
// 		}

// 		//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
// 		a.QMutex.Lock()
// 		a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
// 		a.QMutex.Unlock()
// 		log.Println("Sending Msg:", msg.Header.MsgId)

// 		if device.Mqtt == db.Online {
// 			a.Mqtt.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)
// 		} else if device.Websockets == db.Online {
// 			a.Websockets.Publish(tr369Message, "", "", false)
// 		} else if device.Stomp == db.Online {
// 			//TODO: send stomp message
// 		}

// 		//TODO: verify in protocol and in other models, the Device.Wifi parameters. Maybe in the future, to use SSIDReference from AccessPoint
// 		select {
// 		case msg := <-a.MsgQueue[msg.Header.MsgId]:
// 			log.Printf("Received Msg: %s", msg.Header.MsgId)
// 			a.QMutex.Lock()
// 			delete(a.MsgQueue, msg.Header.MsgId)
// 			a.QMutex.Unlock()
// 			log.Println("requests queue:", a.MsgQueue)
// 			answer := msg.Body.GetResponse().GetGetResp()

// 			var wifi [2]WiFi

// 			//TODO: better algorithm, might use something faster an more reliable
// 			//TODO: full fill the commented wifi resources
// 			for _, x := range answer.ReqPathResults {
// 				if x.RequestedPath == "Device.WiFi.SSID.[Enable==true].SSID" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].SSID = y.ResultParams["SSID"]
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.AccessPoint.[Enable==true].Security.ModeEnabled" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].Security = y.ResultParams["Security.ModeEnabled"]
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.AccessPoint.[Enable==true].Security.ModesSupported" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].SecurityCapabilities = strings.Split(y.ResultParams["Security.ModesSupported"], ",")
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].AutoChannelEnable" {
// 					for i, y := range x.ResolvedPathResults {
// 						autoChannel, err := strconv.ParseBool(y.ResultParams["AutoChannelEnable"])
// 						if err != nil {
// 							log.Println(err)
// 							wifi[i].AutoChannelEnable = false
// 						} else {
// 							wifi[i].AutoChannelEnable = autoChannel
// 						}
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].Channel" {
// 					for i, y := range x.ResolvedPathResults {
// 						channel, err := strconv.Atoi(y.ResultParams["Channel"])
// 						if err != nil {
// 							log.Println(err)
// 							wifi[i].Channel = -1
// 						} else {
// 							wifi[i].Channel = channel
// 						}
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].CurrentOperatingChannelBandwidth" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].ChannelBandwidth = y.ResultParams["CurrentOperatingChannelBandwidth"]
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].OperatingFrequencyBand" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].FrequencyBand = y.ResultParams["OperatingFrequencyBand"]
// 					}
// 					continue
// 				}
// 				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].SupportedOperatingChannelBandwidths" {
// 					for i, y := range x.ResolvedPathResults {
// 						wifi[i].SupportedChannelBandwidths = strings.Split(y.ResultParams["SupportedOperatingChannelBandwidths"], ",")
// 					}
// 					continue
// 				}
// 			}
// 			json.NewEncoder(w).Encode(&wifi)
// 			return
// 		case <-time.After(time.Second * 45):
// 			log.Printf("Request %s Timed Out", msg.Header.MsgId)
// 			w.WriteHeader(http.StatusGatewayTimeout)
// 			a.QMutex.Lock()
// 			delete(a.MsgQueue, msg.Header.MsgId)
// 			a.QMutex.Unlock()
// 			log.Println("requests queue:", a.MsgQueue)
// 			json.NewEncoder(w).Encode("Request Timed Out")
// 			return
// 		}
// 	}
// }

func (a *Api) deviceWifi(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	sn := vars["sn"]

	device, err := getDeviceInfo(w, sn, a.nc)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {

		if device.Cwmp == entity.Online {

			if a.enterpise.Enable {
				a.getEnterpriseResource("wifi", "get", device, sn, w, []byte{}, "cwmp", "098")
				return
			}

			var (
				NUMBER_OF_WIFI_PARAMS_TO_GET = 5
			)

			var wlans []WiFi

			payload := cwmp.GetParameterNames("InternetGatewayDevice.LANDevice.", 1)

			_, response, err := cwmpInteraction[cwmp.GetParameterNamesResponse](sn, []byte(payload), w, a.nc)
			if err != nil {
				return
			}

			parameters_to_get_values := []string{}
			wlanConfigurationInstances := 0

			// x = InternetGatewayDevice.LanDevice.*.
			for _, x := range response.ParameterList {

				payload = cwmp.GetParameterNames(x.Name+"WLANConfiguration.", 1)

				_, response, err := cwmpInteraction[cwmp.GetParameterNamesResponse](sn, []byte(payload), w, a.nc)
				if err != nil {
					return
				}

				// y = InternetGatewayDevice.LanDevice.*.WLANConfiguration.*
				for _, y := range response.ParameterList {
					wlans = append(wlans, WiFi{})

					payload = cwmp.GetParameterNames(y.Name, 1)

					_, response, err := cwmpInteraction[cwmp.GetParameterNamesResponse](sn, []byte(payload), w, a.nc)
					if err != nil {
						return
					}

					// z = InternetGatewayDevice.LanDevice.*.WLANConfiguration.*.<Parameter>
					for _, z := range response.ParameterList {
						path := strings.Split(z.Name, ".")
						parameter := path[len(path)-1]

						switch parameter {
						case "Enable":
							wlans[wlanConfigurationInstances].Enable.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						case "Name":
							wlans[wlanConfigurationInstances].Name.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						case "Status":
							wlans[wlanConfigurationInstances].Status.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						case "SSID":
							wlans[wlanConfigurationInstances].SSID.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						case "Standard":
							wlans[wlanConfigurationInstances].Standard.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						case "KeyPassphrase":
							wlans[wlanConfigurationInstances].Password.Writable = cwmp.ParamTypeIsWritable(z.Writable)
						}

					}

					parameters_to_get_values = append(
						parameters_to_get_values,
						y.Name+"Enable",
						y.Name+"Name",
						y.Name+"Status",
						y.Name+"SSID",
						y.Name+"Standard",
						y.Name+"PreSharedKey.1.KeyPassphrase",
					)

					wlans[wlanConfigurationInstances].Path = y.Name
					wlanConfigurationInstances = wlanConfigurationInstances + 1
				}
			}

			payload = cwmp.GetParameterMultiValues(parameters_to_get_values)

			_, parameterValuesResp, err := cwmpInteraction[cwmp.GetParameterValuesResponse](sn, []byte(payload), w, a.nc)
			if err != nil {
				return
			}

			i := 0
			wlanIndex := 0

			for _, a := range parameterValuesResp.ParameterList {
				path := strings.Split(a.Name, ".")
				parameter := path[len(path)-1]

				switch parameter {
				case "Enable":
					wlans[wlanIndex].Enable.Value = a.Value
				case "Name":
					wlans[wlanIndex].Name.Value = a.Value
				case "Status":
					wlans[wlanIndex].Status.Value = a.Value
				case "SSID":
					wlans[wlanIndex].SSID.Value = a.Value
				case "Standard":
					wlans[wlanIndex].Standard.Value = a.Value
				case "KeyPassphrase":
					wlans[wlanIndex].Password.Value = a.Value
				}

				i = i + 1
				if i == (NUMBER_OF_WIFI_PARAMS_TO_GET + 1) {
					wlanIndex = wlanIndex + 1
					i = 0
				}
			}

			utils.MarshallEncoder(wlans, w)

			return
		}

		if device.Mqtt == entity.Online || device.Stomp == entity.Online || device.Websockets == entity.Online {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write(utils.Marshall("This feature is only working with CWMP devices"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall("Device is Offline"))
	}

	if r.Method == http.MethodPut {

		if device.Cwmp == entity.Online {

			if a.enterpise.Enable {
				payload, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(utils.Marshall(err.Error()))
					return
				}
				a.getEnterpriseResource("wifi", "set", device, sn, w, payload, "cwmp", "098")
				return
			}

			var body []WiFi

			err := utils.MarshallDecoder(&body, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.Marshall("Couldn't decode received payload, err: " + err.Error()))
				return
			}

			fmtBody := map[string]string{}

			for _, x := range body {

				if x.Name.Value != nil {
					fmtBody[x.Path+"Name"] = x.Name.Value.(string)
				}

				if x.SSID.Value != nil {
					fmtBody[x.Path+"SSID"] = x.SSID.Value.(string)
				}

				if x.Enable.Value != nil {
					fmtBody[x.Path+"Enable"] = x.Enable.Value.(string)
				}
			}

			payload := cwmp.SetParameterMultiValues(fmtBody)

			_, setParameterValuesResp, err := cwmpInteraction[cwmp.SetParameterValuesResponse](sn, []byte(payload), w, a.nc)
			if err != nil {
				return
			}

			if setParameterValuesResp.Status == cwmp.ALL_OK {
				log.Printf("All parameters sent to the cpe %s were applied", device.SN)
				w.Write(utils.Marshall(cwmp.ALL_OK))
				return
			}

			if setParameterValuesResp.Status == cwmp.SOME_PENDING {
				log.Printf("All parameters sent to the cpe %s were committed, but not all of them were applied, maybe you need to wait sometime or have a reboot", device.SN)
				w.Write(utils.Marshall(cwmp.SOME_PENDING))
				return
			}

		}

		if device.Mqtt == entity.Online || device.Stomp == entity.Online || device.Websockets == entity.Online {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write(utils.Marshall("This feature is only working with CWMP devices"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall("Device is Offline"))
	}

}
