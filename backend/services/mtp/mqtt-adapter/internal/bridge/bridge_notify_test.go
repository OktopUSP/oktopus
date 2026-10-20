package bridge

import (
	"testing"

	"github.com/OktopUSP/oktopus/backend/services/mqtt-adapter/internal/usp/usp_record"
	"google.golang.org/protobuf/proto"
)

const testSN = "os::TEST-DEVICE-0001"

// buildRecord produces a wire-format USP record with the given from_id, matching
// the record envelope obuspa publishes to the controller base topic. The real
// captured record (live device SN/BSSIDs/SSIDs) is intentionally not committed;
// this synthetic equivalent exercises the same from_id extraction path.
func buildRecord(t *testing.T, fromID string) []byte {
	t.Helper()
	record := usp_record.Record{
		Version:         "1.3",
		ToId:            "oktopus-controller",
		FromId:          fromID,
		PayloadSecurity: usp_record.Record_PLAINTEXT,
		RecordType: &usp_record.Record_NoSessionContext{
			NoSessionContext: &usp_record.NoSessionContextRecord{
				Payload: []byte("opaque-usp-message-payload"),
			},
		},
	}
	b, err := proto.Marshal(&record)
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	return b
}

// TestGetDeviceFromRecord validates that the sender endpoint ID is recoverable
// from a USP record. This is the exact path used when a notification arrives on
// the controller base topic, where the device ID is absent from the MQTT topic.
func TestGetDeviceFromRecord(t *testing.T) {
	payload := buildRecord(t, testSN)
	if got := getDeviceFromRecord(payload); got != testSN {
		t.Fatalf("getDeviceFromRecord() = %q, want %q", got, testSN)
	}
}

// TestGetDeviceFromRecord_Garbage confirms graceful failure (empty string, no
// panic) on an unparseable payload — the handler then drops the message.
func TestGetDeviceFromRecord_Garbage(t *testing.T) {
	if got := getDeviceFromRecord([]byte("not a protobuf record")); got != "" {
		t.Fatalf("getDeviceFromRecord(garbage) = %q, want empty", got)
	}
}

// TestGetDeviceFromTopic_ControllerBase documents that the base controller
// topic yields "controller" (not a device ID), which is what triggers the
// from_id fallback in mqttMessageHandler.
func TestGetDeviceFromTopic_ControllerBase(t *testing.T) {
	if got := getDeviceFromTopic("oktopus/usp/v1/controller"); got != "controller" {
		t.Fatalf("getDeviceFromTopic(base) = %q, want %q", got, "controller")
	}
	dev := "oktopus/usp/v1/controller/" + testSN
	if got := getDeviceFromTopic(dev); got != testSN {
		t.Fatalf("getDeviceFromTopic(per-device) = %q, want %q", got, testSN)
	}
}
