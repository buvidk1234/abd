package im

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	if WsUserID != "sendID" {
		t.Errorf("expected WsUserID to be 'sendID', got '%s'", WsUserID)
	}
	if WSSendMsg != 1001 {
		t.Errorf("expected WSSendMsg to be 1001, got %d", WSSendMsg)
	}
	if writeWait != 10000*time.Second {
		t.Errorf("expected writeWait to be 10000s, got %v", writeWait)
	}
}
