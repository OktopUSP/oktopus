package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
)

type WiFi struct {
	SSID                 string   `json:"ssid"`
	Password             string   `json:"password"`
	Security             string   `json:"security"`
	SecurityCapabilities []string `json:"securityCapabilities"`
	AutoChannelEnable    bool     `json:"autoChannelEnable"`
	Channel              int      `json:"channel"`
	ChannelBandwidth     string   `json:"channelBandwidth"`
	FrequencyBand        string   `json:"frequencyBand"`
	//PossibleChannels     		[]int    `json:"PossibleChannels"`
	SupportedChannelBandwidths []string `json:"supportedChannelBandwidths"`
}

func (a *Api) deviceWifi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	if r.Method == http.MethodGet {
		msg := utils.NewGetMsg(usp_msg.Get{
			ParamPaths: []string{
				"Device.WiFi.SSID.[Enable==true].SSID",
				//"Device.WiFi.AccessPoint.[Enable==true].SSIDReference",
				"Device.WiFi.AccessPoint.[Enable==true].Security.ModeEnabled",
				"Device.WiFi.AccessPoint.[Enable==true].Security.ModesSupported",
				//"Device.WiFi.EndPoint.[Enable==true].",
				"Device.WiFi.Radio.[Enable==true].AutoChannelEnable",
				"Device.WiFi.Radio.[Enable==true].Channel",
				"Device.WiFi.Radio.[Enable==true].CurrentOperatingChannelBandwidth",
				"Device.WiFi.Radio.[Enable==true].OperatingFrequencyBand",
				//"Device.WiFi.Radio.[Enable==true].PossibleChannels",
				"Device.WiFi.Radio.[Enable==true].SupportedOperatingChannelBandwidths",
			},
			MaxDepth: 2,
		})

		encodedMsg, err := proto.Marshal(&msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		record := utils.NewUspRecord(encodedMsg, sn)
		tr369Message, err := proto.Marshal(&record)
		if err != nil {
			log.Fatalln("Failed to encode tr369 record:", err)
		}

		//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
		a.QMutex.Lock()
		a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
		a.QMutex.Unlock()
		log.Println("Sending Msg:", msg.Header.MsgId)
		a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

		//TODO: verify in protocol and in other models, the Device.Wifi parameters. Maybe in the future, to use SSIDReference from AccessPoint
		select {
		case msg := <-a.MsgQueue[msg.Header.MsgId]:
			log.Printf("Received Msg: %s", msg.Header.MsgId)
			a.QMutex.Lock()
			delete(a.MsgQueue, msg.Header.MsgId)
			a.QMutex.Unlock()
			log.Println("requests queue:", a.MsgQueue)
			answer := msg.Body.GetResponse().GetGetResp()

			var wifi [2]WiFi

			//TODO: better algorithm, might use something faster an more reliable
			//TODO: full fill the commented wifi resources
			for _, x := range answer.ReqPathResults {
				if x.RequestedPath == "Device.WiFi.SSID.[Enable==true].SSID" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].SSID = y.ResultParams["SSID"]
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.AccessPoint.[Enable==true].Security.ModeEnabled" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].Security = y.ResultParams["Security.ModeEnabled"]
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.AccessPoint.[Enable==true].Security.ModesSupported" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].SecurityCapabilities = strings.Split(y.ResultParams["Security.ModesSupported"], ",")
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].AutoChannelEnable" {
					for i, y := range x.ResolvedPathResults {
						autoChannel, err := strconv.ParseBool(y.ResultParams["AutoChannelEnable"])
						if err != nil {
							log.Println(err)
							wifi[i].AutoChannelEnable = false
						} else {
							wifi[i].AutoChannelEnable = autoChannel
						}
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].Channel" {
					for i, y := range x.ResolvedPathResults {
						channel, err := strconv.Atoi(y.ResultParams["Channel"])
						if err != nil {
							log.Println(err)
							wifi[i].Channel = -1
						} else {
							wifi[i].Channel = channel
						}
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].CurrentOperatingChannelBandwidth" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].ChannelBandwidth = y.ResultParams["CurrentOperatingChannelBandwidth"]
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].OperatingFrequencyBand" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].FrequencyBand = y.ResultParams["OperatingFrequencyBand"]
					}
					continue
				}
				if x.RequestedPath == "Device.WiFi.Radio.[Enable==true].SupportedOperatingChannelBandwidths" {
					for i, y := range x.ResolvedPathResults {
						wifi[i].SupportedChannelBandwidths = strings.Split(y.ResultParams["SupportedOperatingChannelBandwidths"], ",")
					}
					continue
				}
			}
			json.NewEncoder(w).Encode(&wifi)
			return
		case <-time.After(time.Second * 45):
			log.Printf("Request %s Timed Out", msg.Header.MsgId)
			w.WriteHeader(http.StatusGatewayTimeout)
			a.QMutex.Lock()
			delete(a.MsgQueue, msg.Header.MsgId)
			a.QMutex.Unlock()
			log.Println("requests queue:", a.MsgQueue)
			json.NewEncoder(w).Encode("Request Timed Out")
			return
		}
	}
}
