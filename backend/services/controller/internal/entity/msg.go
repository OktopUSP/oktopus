package entity

import "time"

type DataType interface {
	[]map[string]interface{} | *string | Device | int64 | []Device | []VendorsCount | []ProductClassCount | []StatusCount | time.Duration | []byte
}

type MsgAnswer[T DataType] struct {
	Code int
	Msg  T
}
