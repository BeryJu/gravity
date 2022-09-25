package types

const SessionName = "gravity_session"

type RequestContextValue string

const RequestSession RequestContextValue = "session"

const (
	SessionKeyUser      = "user"
	SessionKeyOIDCState = "oidc_state"
	SessionKeyDirty     = "dirty"
)
