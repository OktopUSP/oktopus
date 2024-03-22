package api

// func (a *Api) retrieveDevices(w http.ResponseWriter, r *http.Request) {
// 	const PAGE_SIZE_LIMIT = 50
// 	const PAGE_SIZE_DEFAULT = 20

// 	// Get specific device
// 	id := r.URL.Query().Get("id")
// 	if id != "" {
// 		device, err := a.Db.RetrieveDevice(id)
// 		if err != nil {
// 			if err == mongo.ErrNoDocuments {
// 				json.NewEncoder(w).Encode("Device id: " + id + " not found")
// 				return
// 			}
// 			json.NewEncoder(w).Encode(err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		err = json.NewEncoder(w).Encode(device)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		return
// 	}

// 	// Get devices with pagination
// 	page_n := r.URL.Query().Get("page_number")
// 	page_s := r.URL.Query().Get("page_size")
// 	var err error

// 	var page_number int64
// 	if page_n == "" {
// 		page_number = 0
// 	} else {
// 		page_number, err = strconv.ParseInt(page_n, 10, 64)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode("Page number must be an integer")
// 			return
// 		}
// 	}

// 	var page_size int64
// 	if page_s != "" {
// 		page_size, err = strconv.ParseInt(page_s, 10, 64)

// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode("Page size must be an integer")
// 			return
// 		}

// 		if page_size > PAGE_SIZE_LIMIT {
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode("Page size must not exceed " + strconv.Itoa(PAGE_SIZE_LIMIT))
// 			return
// 		}

// 	} else {
// 		page_size = PAGE_SIZE_DEFAULT
// 	}

// 	total, err := a.Db.RetrieveDevicesCount(bson.M{})
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Unable to get devices count from database")
// 		return
// 	}

// 	skip := page_number * (page_size - 1)
// 	if total < page_size {
// 		skip = 0
// 	}

// 	//TODO: Create filters
// 	//TODO: Create sorting

// 	sort := bson.M{}
// 	sort["status"] = 1

// 	filter := bson.A{
// 		//bson.M{"$match": filter},
// 		bson.M{"$sort": sort}, // shows online devices first
// 		bson.M{"$skip": skip},
// 		bson.M{"$limit": page_size},
// 	}

// 	devices, err := a.Db.RetrieveDevices(filter)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Unable to aggregate database devices info")
// 		return
// 	}

// 	err = json.NewEncoder(w).Encode(map[string]interface{}{
// 		"pages":   total / page_size,
// 		"page":    page_number,
// 		"size":    page_size,
// 		"devices": devices,
// 	})
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
