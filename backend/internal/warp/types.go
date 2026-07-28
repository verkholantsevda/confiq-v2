package warp

type Device struct {
	ID    string
	Token string
}

type Interface struct {
	PeerPublicKey string

	ClientIPv4 string
	ClientIPv6 string
}

type Config struct {
	Device

	PrivateKey string
	PublicKey  string

	Interface
}
