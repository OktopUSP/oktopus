package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/leandrofars/oktopus/internal/bridge"
	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *Api) retrieveDevices(w http.ResponseWriter, r *http.Request) {
	const PAGE_SIZE_LIMIT = 50
	const PAGE_SIZE_DEFAULT = 20

	// Get specific device
	id := r.URL.Query().Get("id")
	if id != "" {
		device, err := getDeviceInfo(w, id, a.nc)
		if err != nil {
			return
		}
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Get devices with pagination
	page_n := r.URL.Query().Get("page_number")
	page_s := r.URL.Query().Get("page_size")
	var err error

	var page_number int64
	if page_n == "" {
		page_number = 0
	} else {
		page_number, err = strconv.ParseInt(page_n, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Page number must be an integer")
			return
		}
	}

	var page_size int64
	if page_s != "" {
		page_size, err = strconv.ParseInt(page_s, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Page size must be an integer")
			return
		}

		if page_size > PAGE_SIZE_LIMIT {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Page size must not exceed " + strconv.Itoa(PAGE_SIZE_LIMIT))
			return
		}

	} else {
		page_size = PAGE_SIZE_DEFAULT
	}

	total, err := getDeviceCount(w, a.nc)
	if err != nil {
		return
	}

	skip := page_number * (page_size - 1)
	if total < page_size {
		skip = 0
	}

	//TODO: fix status ordering
	statusOrder := r.URL.Query().Get("status")
	if statusOrder != "" {
		if statusOrder == "asc" {
			statusOrder = "1"
		} else if statusOrder == "desc" {
			statusOrder = "-1"
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Status order must be 'asc' or 'desc'")
			return
		}
	}

	sort := bson.M{}
	sort["status"] = 1

	//TODO: Create filters

	filter := bson.A{
		//bson.M{"$match": filter},
		bson.M{"$sort": sort}, // shows online devices first
		bson.M{"$skip": skip},
		bson.M{"$limit": page_size},
	}

	devices, err := getDevices(w, filter, a.nc)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"pages":   total / page_size,
		"page":    page_number,
		"size":    page_size,
		"devices": devices,
	})
	if err != nil {
		log.Println(err)
	}
}

type DeviceAuth struct {
	User     string `json:"id"`
	Password string `json:"password"`
}

func (a *Api) deviceAuth(w http.ResponseWriter, r *http.Request) {

	user, err := a.db.FindUser(r.Context().Value("email").(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.MarshallEncoder(err, w)
		return
	}
	if user.Level != AdminUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == http.MethodGet {

		id := r.URL.Query().Get("id")
		if id != "" {
			entry, err := a.kv.Get(r.Context(), id)
			if err != nil {
				if err == jetstream.ErrKeyNotFound {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				utils.MarshallEncoder(err, w)
				return
			}
			utils.MarshallEncoder(map[string]string{
				id: string(entry.Value()),
			}, w)
			return
		}

		entries, err := a.kv.ListKeys(r.Context(), jetstream.IgnoreDeletes())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.MarshallEncoder(err, w)
			return
		}

		listOfKeys := make(map[string]string)

		keys := entries.Keys()
		for key := range keys {
			entry, err := a.kv.Get(r.Context(), key)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				utils.MarshallEncoder(err, w)
				return
			}

			/*listOfKeys = append(listOfKeys, map[string]string{
				key: string(entry.Value()),
			})*/
			listOfKeys[key] = string(entry.Value())
		}

		utils.MarshallEncoder(listOfKeys, w)

	} else if r.Method == http.MethodDelete {

		id := r.URL.Query().Get("id")
		if id != "" {
			err := a.kv.Purge(r.Context(), id)
			if err != nil {
				if err == jetstream.ErrKeyNotFound {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				utils.MarshallEncoder(err, w)
				return
			}
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("No id provided", w)

	} else if r.Method == http.MethodPost {

		var deviceAuth DeviceAuth

		err := utils.MarshallDecoder(&deviceAuth, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.MarshallEncoder(err, w)
			return
		}

		if deviceAuth.User != "" {
			_, err := a.kv.Get(r.Context(), deviceAuth.User)

			if err != nil {

				if err == jetstream.ErrKeyNotFound {
					_, err = a.kv.PutString(r.Context(), deviceAuth.User, deviceAuth.Password)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						utils.MarshallEncoder(err, w)
					}
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				utils.MarshallEncoder(err, w)
				return
			}

			w.WriteHeader(http.StatusConflict)
			utils.MarshallEncoder("Username already exists", w)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("device must have a user", w)

	} else {
		log.Println("Unknown method used in device auth api")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (a *Api) setDeviceAlias(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("No id provided", w)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("Error to decode payload: "+err.Error(), w)
		return
	}

	payloadLen := len(payload)
	if payloadLen == 0 {
		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("No alias provided", w)
		return
	}
	if payloadLen > 50 {
		w.WriteHeader(http.StatusBadRequest)
		utils.MarshallEncoder("Alias too long", w)
		return
	}

	_, err = bridge.NatsReq[[]byte](local.NATS_ADAPTER_SUBJECT+id+".device.alias", payload, w, a.nc)
	if err != nil {
		return
	}
}
