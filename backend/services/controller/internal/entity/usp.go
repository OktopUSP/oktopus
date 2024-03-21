package entity

import "github.com/leandrofars/oktopus/internal/usp/usp_msg"

type UspType interface {
	usp_msg.GetSupportedDM
}