package wbstream

import (
	"fmt"
	"strings"
	"time"

	lkauth "github.com/livekit/protocol/auth"
)

// SelfHostToken generates a LiveKit access token (JWT) for connecting to a
// self-hosted LiveKit server. identity is the participant display name,
// apiKey/apiSecret are the credentials configured on the self-hosted server.
//
// The returned token is valid for 24 hours; callers that run long-lived
// sessions should regenerate it periodically (see the rejoin loops in the
// headless creator/joiner).
func SelfHostToken(roomID, identity, apiKey, apiSecret string) (string, error) {
	if roomID == "" {
		return "", fmt.Errorf("self-host token: empty room id")
	}
	if identity == "" {
		return "", fmt.Errorf("self-host token: empty identity")
	}
	if apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("self-host token: api key and secret are required")
	}
	at := lkauth.NewAccessToken(apiKey, apiSecret)
	at.SetIdentity(identity)
	at.SetValidFor(24 * time.Hour)
	grant := &lkauth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.AddGrant(grant)
	token, err := at.ToJWT()
	if err != nil {
		return "", fmt.Errorf("self-host token: %w", err)
	}
	return token, nil
}

// NormalizeServerURL ensures a self-hosted LiveKit URL carries a ws:// or
// wss:// scheme. A bare host:port (e.g. "45.13.239.115:7880") is treated as
// plain ws. Empty input is returned unchanged.
func NormalizeServerURL(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "ws://") && !strings.HasPrefix(s, "wss://") {
		return "ws://" + s
	}
	return s
}
