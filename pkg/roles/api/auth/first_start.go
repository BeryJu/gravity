package auth

import (
	"encoding/base64"
	"fmt"
	"os"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
)

const firstStartHeader = `------------------------------------------------------------
Welcome to gravity! An admin user has been created for you.
Username: '%s'
Password: '%s'`

const firstStartFooter = `
Open 'http://%s:8008/' to start using Gravity!
------------------------------------------------------------
`

func (ap *AuthProvider) FirstStart(ev *roles.Event) {
	username := "admin"
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	}
	err := ap.CreateUser(ev.Context, username, password)
	if err != nil {
		ap.log.Warn("failed to create default user", zap.Error(err))
		return
	}
	text := fmt.Sprintf(firstStartHeader, username, password)

	token := os.Getenv("ADMIN_TOKEN")
	if token != "" {
		t := &types.Token{
			Key:      token,
			Username: username,
		}
		err := ap.putToken(t, ev.Context)
		if err != nil {
			ap.log.Warn("failed to create bootstrap token", zap.Error(err))
			return
		}
		text += fmt.Sprintf("\nToken: '%s'", token)
	}
	text += fmt.Sprintf(firstStartFooter, extconfig.Get().Instance.IP)
	fmt.Print(text)
}
