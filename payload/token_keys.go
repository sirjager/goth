package payload

import "fmt"

func SessionKey(userID, sessionID string) string {
	return fmt.Sprintf("sess:%s:%s", userID, sessionID)
}

func UserSessionsKey(userID string) string {
	return fmt.Sprintf("sess:%s", userID)
}

func SessionRefreshKey(userID, sessionID string) string {
	return fmt.Sprintf("sess:%s:%s:%s", userID, sessionID, TypeRefresh)
}

func SessionAccessKey(userID, sessionID string) string {
	return fmt.Sprintf("sess:%s:%s:%s", userID, sessionID, TypeAccess)
}
