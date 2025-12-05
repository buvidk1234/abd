package im

import (
	"sync"
	"time"
)

type UserPlatform struct {
	LastActive time.Time
	Clients    []*Client
}

func (u *UserPlatform) PlatformIDs() []int32 {
	platformIDs := make([]int32, 0, len(u.Clients))
	for _, c := range u.Clients {
		platformIDs = append(platformIDs, int32(c.PlatformID))
	}
	return platformIDs
}

type UserState struct {
	UserID  string
	Online  []int32
	Offline []int32
}

type UserMap interface {
	GetAll(userID string) ([]*Client, bool)
	Get(userID string, platformID int) ([]*Client, bool, bool)
	Set(userID string, v *Client)
	DeleteClients(userID string, clients []*Client) (isDeleteUser bool)
	UserState() <-chan UserState
	GetAllUserStatus(deadline time.Time, nowtime time.Time) []UserState
	RecvSubChange(userID string, platformIDs []int32) bool
}

type userMap struct {
	mu              sync.RWMutex
	userPlatformMap map[string]*UserPlatform
	ch              chan UserState
}

func newUserMap() UserMap {
	return &userMap{
		userPlatformMap: make(map[string]*UserPlatform),
		ch:              make(chan UserState, 1024),
	}
}

func (u *userMap) GetAll(userID string) ([]*Client, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	userPlatform, ok := u.userPlatformMap[userID]
	if !ok {
		return nil, false
	}
	return userPlatform.Clients, true
}

func (u *userMap) Get(userID string, platformID int) ([]*Client, bool, bool) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	userPlatform, ok := u.userPlatformMap[userID]
	if !ok {
		return nil, false, false
	}
	var clients []*Client
	for _, client := range userPlatform.Clients {
		if client.PlatformID == platformID {
			clients = append(clients, client)
		}
	}
	return clients, true, len(clients) > 0
}

func (u *userMap) Set(userID string, v *Client) {
	u.mu.Lock()
	defer u.mu.Unlock()
	userPlatform, ok := u.userPlatformMap[userID]
	if !ok {
		userPlatform = &UserPlatform{
			Clients: []*Client{v},
		}
		u.userPlatformMap[userID] = userPlatform
	} else {
		userPlatform.Clients = append(userPlatform.Clients, v)
	}
}

func (u *userMap) DeleteClients(userID string, clients []*Client) (isDeleteUser bool) {
	if len(clients) == 0 {
		return false
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	userPlatform, ok := u.userPlatformMap[userID]
	if !ok {
		return false
	}
	offline := make([]int32, 0, len(clients))
	deleteAddr := make(map[string]struct{})
	for _, c := range clients {
		deleteAddr[c.req.RemoteAddr] = struct{}{}
	}

	oldClients := userPlatform.Clients
	userPlatform.Clients = []*Client{}

	for _, c := range oldClients {
		if _, delCli := deleteAddr[c.req.RemoteAddr]; delCli {
			offline = append(offline, int32(c.PlatformID))
		} else {
			userPlatform.Clients = append(userPlatform.Clients, c)
		}
	}
	defer u.push(userID, userPlatform, offline)
	if len(userPlatform.Clients) > 0 {
		return false
	}
	delete(u.userPlatformMap, userID)
	return true
}

func (u *userMap) push(userID string, userPlatform *UserPlatform, offline []int32) bool {
	online := make([]int32, 0, len(userPlatform.Clients))
	for _, c := range userPlatform.Clients {
		online = append(online, int32(c.PlatformID))
	}
	select {
	case u.ch <- UserState{
		UserID:  userID,
		Online:  online,
		Offline: offline,
	}:
		userPlatform.LastActive = time.Now()
		return true
	default:
		return false
	}

}

func (u *userMap) UserState() <-chan UserState {
	return u.ch
}
func (u *userMap) GetAllUserStatus(deadline time.Time, nowtime time.Time) []UserState {
	u.mu.RLock()
	defer u.mu.RUnlock()

	results := make([]UserState, 0, len(u.userPlatformMap))
	for userID, userPlatform := range u.userPlatformMap {
		if deadline.Before(userPlatform.LastActive) {
			continue
		}
		online := make([]int32, 0, len(userPlatform.Clients))
		for _, c := range userPlatform.Clients {
			online = append(online, int32(c.PlatformID))
		}
		results = append(results, UserState{
			UserID: userID,
			Online: online,
		})
	}
	return results
}
func (u *userMap) RecvSubChange(userID string, platformIDs []int32) bool {
	// TODO: implement subscription change handling
	return true
}
