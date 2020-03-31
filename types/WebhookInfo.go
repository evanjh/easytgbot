package types

type WebhookInfo struct {
	Url                  string
	HasCustomCertificate bool
	PendingUpdateCount   int
	LastErrorDate        int
	LastErrorMessage     string
	MaxConnections       int
	AllowedUpdates       []string
}

func (info WebhookInfo) IsSet() bool {
	return info.Url != ""
}
