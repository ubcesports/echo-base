package auth

type Application struct {
	AppName   string
	KeyId     string
	HashedKey []byte
}

type APIKey struct {
	KeyId   string
	APIKey  string
	AppName string
}
