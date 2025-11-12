package auth

import "sync"

// ------- State Token : SID MAP ----------
var (
	stateSidMap = make(map[string]string)
	mapMu       sync.RWMutex
)

// Create + Update
func SetStateSid(stateToken, sessionID string) {
	mapMu.Lock()
	stateSidMap[stateToken] = sessionID
	mapMu.Unlock()
}

// Read
func GetStateSid(sessionID string) (string, bool) {
	mapMu.Lock()
	key, ok := stateSidMap[sessionID]
	mapMu.Unlock()
	return key, ok

}

// Delete
func DelStateSid(sessionID string) {
	mapMu.Lock()
	delete(stateSidMap, sessionID)
	mapMu.Unlock()

}

// ------- SID : WEB SESSION KEY MAP ----------

var (
	sidKeyMap  = make(map[string]string)
	mapSKMutex sync.RWMutex
)

// Create + Update
func SetSidKey(sessionID, sessionKey string) {
	mapSKMutex.Lock()
	sidKeyMap[sessionID] = sessionKey
	mapSKMutex.Unlock()
}

// Read
// sessionID (str) : SessionID
// Returns sessionKey (str), ok
func GetSidKey(sessionID string) (string, bool) {
	mapSKMutex.Lock()
	key, ok := sidKeyMap[sessionID]
	mapSKMutex.Unlock()
	return key, ok

}

// Delete
func DelSidKey(sessionID string) {
	mapSKMutex.Lock()
	delete(sidKeyMap, sessionID)
	mapSKMutex.Unlock()

}

// ------- Dual Map functions ----------

