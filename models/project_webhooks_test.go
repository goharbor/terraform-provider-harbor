package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWebHookTargetsPayloadFormatJSON(t *testing.T) {
	target := WebHookTargets{
		Type:          "http",
		Address:       "https://example.com/hook",
		PayloadFormat: "CloudEvents",
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"payload_format":"CloudEvents"`) {
		t.Errorf("payload_format not serialized as expected, got %s", data)
	}

	var decoded WebHookTargets
	if err := json.Unmarshal([]byte(`{"type":"http","payload_format":"CloudEvents"}`), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.PayloadFormat != "CloudEvents" {
		t.Errorf("payload_format not read from Harbor response, got %q", decoded.PayloadFormat)
	}
}

func TestWebHookTargetsPayloadFormatOmittedWhenEmpty(t *testing.T) {
	data, err := json.Marshal(WebHookTargets{Type: "http"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "payload_format") || strings.Contains(string(data), "PayloadFormat") {
		t.Errorf("empty payload_format should be omitted, got %s", data)
	}
}
