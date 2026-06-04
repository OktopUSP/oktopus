package usp_handler

import (
	"encoding/json"
	"log"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_msg"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_record"
	"google.golang.org/protobuf/proto"
)

// OperationResult is the structured payload published on device.v1.notify when
// an async Operate finishes. It carries the command output args back to any
// downstream consumer (controller API, frontend socket push, persistence).
type OperationResult struct {
	SN          string            `json:"sn"`
	CommandName string            `json:"command_name"`
	CommandKey  string            `json:"command_key"`
	ObjPath     string            `json:"obj_path"`
	OutputArgs  map[string]string `json:"output_args,omitempty"`
	Failed      bool              `json:"failed"`
	ErrMsg      string            `json:"err_msg,omitempty"`
}

// HandleDeviceNotify processes unsolicited USP notifications that arrive on the
// controller base topic (routed by the mqtt-adapter to the .notify subject).
// It currently handles OperationComplete — the result of an async Operate —
// and republishes it as a structured event for downstream consumers.
func (h *Handler) HandleDeviceNotify(device, subject string, data []byte, mtp string, ack func()) {
	defer ack()

	result, ok := parseOperationComplete(device, data)
	if !ok {
		return
	}

	log.Printf("Notify: OperationComplete from %s for %s (key=%s, %d output args, failed=%t)",
		device, result.CommandName, result.CommandKey, len(result.OutputArgs), result.Failed)

	payload, err := json.Marshal(result)
	if err != nil {
		log.Printf("Notify: failed to marshal OperationResult: %v", err)
		return
	}

	if err := h.nc.Publish("device.v1.notify", payload); err != nil {
		log.Printf("Notify: failed to publish operation result: %v", err)
	}
}

// parseOperationComplete decodes a USP record into an OperationResult. It
// returns ok=false (and logs) for unparseable records or notifications that are
// not OperationComplete. Pure function — no side effects — so it is unit
// testable against captured wire data.
func parseOperationComplete(device string, data []byte) (OperationResult, bool) {
	var record usp_record.Record
	if err := proto.Unmarshal(data, &record); err != nil {
		log.Printf("Notify: failed to unmarshal USP record: %v", err)
		return OperationResult{}, false
	}

	var msg usp_msg.Msg
	if err := proto.Unmarshal(record.GetNoSessionContext().GetPayload(), &msg); err != nil {
		log.Printf("Notify: failed to unmarshal USP message: %v", err)
		return OperationResult{}, false
	}

	notify := msg.GetBody().GetRequest().GetNotify()
	if notify == nil {
		log.Printf("Notify: message on .notify subject is not a Notify")
		return OperationResult{}, false
	}

	operComplete := notify.GetOperComplete()
	if operComplete == nil {
		// Other notification types (ValueChange, ObjectCreation, Event, ...)
		// are not yet consumed downstream — acknowledge and drop for now.
		log.Printf("Notify: ignoring non-OperationComplete notification from %s", device)
		return OperationResult{}, false
	}

	result := OperationResult{
		SN:          device,
		CommandName: operComplete.GetCommandName(),
		CommandKey:  operComplete.GetCommandKey(),
		ObjPath:     operComplete.GetObjPath(),
	}

	if outArgs := operComplete.GetReqOutputArgs(); outArgs != nil {
		result.OutputArgs = outArgs.GetOutputArgs()
	} else if failure := operComplete.GetCmdFailure(); failure != nil {
		result.Failed = true
		result.ErrMsg = failure.GetErrMsg()
	}

	return result, true
}
