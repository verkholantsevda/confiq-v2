package warp

import (
	"strconv"
	"strings"
)

type TemplateData struct {
	PrivateKey string
	PublicKey  string

	PeerPublicKey string

	ClientIPv4 string
	ClientIPv6 string

	Endpoint string
	Port     int
}

func Render(templateText string, data TemplateData) (string, error) {
	replacer := strings.NewReplacer(
		"{{ private_key }}", data.PrivateKey,
		"{{ public_key }}", data.PublicKey,
		"{{ peer_public_key }}", data.PeerPublicKey,
		"{{ client_ipv4 }}", data.ClientIPv4,
		"{{ client_ipv6 }}", data.ClientIPv6,
		"{{ endpoint }}", data.Endpoint,
		"{{ port }}", strconv.Itoa(data.Port),
	)

	return replacer.Replace(templateText), nil
}
