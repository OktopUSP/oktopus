package usp_handler

import (
	"fmt"
	"testing"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_msg"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_record"
	"google.golang.org/protobuf/proto"
)

const testSN = "os::TEST-DEVICE-0001"

// buildOperationCompleteRecord produces a wire-format USP record carrying an
// OperationComplete notification with `n` synthetic neighbor results. It mirrors
// the exact structure obuspa emits for an async NeighboringWiFiDiagnostic — the
// real captured record (with live device SN/BSSIDs/SSIDs) is intentionally NOT
// committed; this synthetic equivalent validates the same parse path.
func buildOperationCompleteRecord(t *testing.T, fromID string, n int) []byte {
	t.Helper()

	outputArgs := map[string]string{"Status": "Success"}
	for i := 1; i <= n; i++ {
		outputArgs[fmt.Sprintf("Result.%d.SSID", i)] = fmt.Sprintf("test-ssid-%d", i)
		outputArgs[fmt.Sprintf("Result.%d.BSSID", i)] = fmt.Sprintf("00:00:5e:00:53:%02x", i%256)
		outputArgs[fmt.Sprintf("Result.%d.Channel", i)] = "36"
		outputArgs[fmt.Sprintf("Result.%d.SignalStrength", i)] = "-70"
	}

	msg := usp_msg.Msg{
		Body: &usp_msg.Body{
			MsgBody: &usp_msg.Body_Request{
				Request: &usp_msg.Request{
					ReqType: &usp_msg.Request_Notify{
						Notify: &usp_msg.Notify{
							SubscriptionId: "test-sub",
							Notification: &usp_msg.Notify_OperComplete{
								OperComplete: &usp_msg.Notify_OperationComplete{
									ObjPath:     "Device.LocalAgent.Request.1.",
									CommandName: "NeighboringWiFiDiagnostic()",
									CommandKey:  "test-key",
									OperationResp: &usp_msg.Notify_OperationComplete_ReqOutputArgs{
										ReqOutputArgs: &usp_msg.Notify_OperationComplete_OutputArgs{
											OutputArgs: outputArgs,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	msgBytes, err := proto.Marshal(&msg)
	if err != nil {
		t.Fatalf("marshal msg: %v", err)
	}

	record := usp_record.Record{
		Version:         "1.3",
		ToId:            "oktopus-controller",
		FromId:          fromID,
		PayloadSecurity: usp_record.Record_PLAINTEXT,
		RecordType: &usp_record.Record_NoSessionContext{
			NoSessionContext: &usp_record.NoSessionContextRecord{
				Payload: msgBytes,
			},
		},
	}

	recBytes, err := proto.Marshal(&record)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	return recBytes
}

// TestParseOperationComplete validates the consumer-side parse against a
// synthetic OperationComplete record with 45 neighbor results. Confirms the
// command name and that all result fields are recovered into output args.
func TestParseOperationComplete(t *testing.T) {
	const neighbors = 45
	data := buildOperationCompleteRecord(t, testSN, neighbors)

	result, ok := parseOperationComplete(testSN, data)
	if !ok {
		t.Fatal("parseOperationComplete returned ok=false for a valid OperationComplete record")
	}
	if result.SN != testSN {
		t.Errorf("SN = %q, want %q", result.SN, testSN)
	}
	if result.CommandName != "NeighboringWiFiDiagnostic()" {
		t.Errorf("CommandName = %q, want NeighboringWiFiDiagnostic()", result.CommandName)
	}
	if result.Failed {
		t.Errorf("Failed = true, want false (ErrMsg=%q)", result.ErrMsg)
	}
	// 45 neighbors * 4 fields + Status = 181 output args.
	if len(result.OutputArgs) < 100 {
		t.Errorf("OutputArgs count = %d, want >= 100", len(result.OutputArgs))
	}
	if _, ok := result.OutputArgs["Result.1.SSID"]; !ok {
		t.Errorf("expected Result.1.SSID in output args; got %d keys", len(result.OutputArgs))
	}
	t.Logf("parsed %d output args for %s", len(result.OutputArgs), result.CommandName)
}

// TestParseOperationComplete_Garbage confirms graceful rejection.
func TestParseOperationComplete_Garbage(t *testing.T) {
	if _, ok := parseOperationComplete("sn", []byte("garbage")); ok {
		t.Fatal("expected ok=false for garbage payload")
	}
}
