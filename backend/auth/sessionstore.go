package auth

import "sync"

// ------- State Token : SID MAP ----------
var (
	stateSidMap = make(map[string]string)
	mapMu       sync.RWMutex
)

// Create + Update
func setStateSid(stateToken, sessionID string) {
	mapMu.Lock()
	stateSidMap[stateToken] = sessionID
	mapMu.Unlock()
}

// Read
func getStateSid(sessionID string) (string, bool) {
	mapMu.Lock()
	key, ok := stateSidMap[sessionID]
	mapMu.Unlock()
	return key, ok

}

// Delete
func delStateSid(sessionID string) {
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
func setSidKey(sessionID, sessionKey string) {
	mapMu.Lock()
	sidKeyMap[sessionID] = sessionKey
	mapMu.Unlock()
}

// Read
// sessionID (str) : SessionID
// Returns sessionKey (str), ok
func getSidKey(sessionID string) (string, bool) {
	mapMu.Lock()
	key, ok := sidKeyMap[sessionID]
	mapMu.Unlock()
	return key, ok

}

// Delete
func delSidKey(sessionID string) {
	mapMu.Lock()
	delete(sidKeyMap, sessionID)
	mapMu.Unlock()

}

// ------- Dual Map functions ----------

func consumeAndSetSID() {}
