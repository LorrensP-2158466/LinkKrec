package usersession

import (
	"context"
	"encoding/gob"
)

const (
	SessionInfoKey = "sessInfo"
	GinContextKey  = "GinContextKey"
)

type UserSessionInfo struct {
	// convenience to quickly determine of this user has a completed account
	IsComplete bool `json:"ProfileCompleted"`
	IsUser     bool
	Email      string `json:"email"`
	Id         string `json:"id"`
	Cookie     string
	// TODO: More?
}

func init() {
	gob.Register(UserSessionInfo{}) // Register the type for serialization
}

func For(ctx context.Context) UserSessionInfo {
	return ctx.Value(SessionInfoKey).(UserSessionInfo)
}
