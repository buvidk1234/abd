package im

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func newTestClient(platformID int, addr string) *Client {
	return &Client{
		PlatformID: platformID,
		req: &http.Request{
			RemoteAddr: addr,
		},
	}
}

func TestUserMap_SetGetDelete(t *testing.T) {
	um := newUserMap()

	c1 := newTestClient(1, "addr1")
	c2 := newTestClient(2, "addr2")

	um.Set("u1", c1)
	um.Set("u1", c2)

	// GetAll
	all, ok := um.GetAll("u1")
	if !ok {
		t.Fatalf("expected user present")
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(all))
	}

	// Get by platform
	byPlat, ok, exists := um.Get("u1", 1)
	if !ok || !exists || len(byPlat) != 1 {
		t.Fatalf("expected to find platform 1 client")
	}

	// Delete one client and expect user still exists
	ch := um.UserState()
	deleted := um.DeleteClients("u1", []*Client{c1})
	if deleted {
		t.Fatalf("expected user not deleted after removing one of two clients")
	}

	// Expect a state change with offline platform 1
	select {
	case st := <-ch:
		if !reflect.DeepEqual(st.Offline, []int32{1}) {
			t.Fatalf("unexpected offline platforms: %#v", st.Offline)
		}
		if !reflect.DeepEqual(st.Online, []int32{2}) {
			t.Fatalf("unexpected online platforms: %#v", st.Online)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for user state update")
	}

	// Now delete remaining client and expect user removed
	deleted = um.DeleteClients("u1", []*Client{c2})
	if !deleted {
		t.Fatalf("expected user deleted after removing last client")
	}

	// Expect another state with offline platform 2
	select {
	case st := <-ch:
		if !reflect.DeepEqual(st.Offline, []int32{2}) {
			t.Fatalf("unexpected offline after final delete: %#v", st.Offline)
		}
		if len(st.Online) != 0 {
			t.Fatalf("expected no online platforms, got %#v", st.Online)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for second user state update")
	}
}

func TestUserMap_GetAllUserStatus_And_RecvSubChange(t *testing.T) {
	um := newUserMap()

	// Add a user (LastActive will be zero time)
	c := newTestClient(3, "addr3")
	um.Set("u2", c)

	// deadline = now -> since LastActive is zero, user should be included
	res := um.GetAllUserStatus(time.Now(), time.Now())
	if len(res) != 1 {
		t.Fatalf("expected 1 user in GetAllUserStatus, got %d", len(res))
	}

	if res[0].UserID != "u2" {
		t.Fatalf("unexpected userid: %s", res[0].UserID)
	}

	// RecvSubChange currently returns true (TODO)
	if !um.RecvSubChange("u2", []int32{3}) {
		t.Fatalf("expected RecvSubChange to return true")
	}
}
