package warp

type Generator struct {
	client *Client
}

func NewGenerator() *Generator {
	return &Generator{
		client: NewClient(),
	}
}

type GenerateRequest struct {
	EndpointAddress string
	EndpointPort    int
}

func (g *Generator) Generate(req GenerateRequest) (*Config, error) {
	privateKey, publicKey, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	device, err := g.client.RegisterDevice(publicKey)
	if err != nil {
		return nil, err
	}

	iface, err := g.client.EnableWarp(device.ID, device.Token)
	if err != nil {
		_ = g.client.DeleteDevice(device.ID, device.Token)
		return nil, err
	}

	return &Config{
		Device: Device{
			ID:    device.ID,
			Token: device.Token,
		},
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Interface: Interface{
			PeerPublicKey: iface.PeerPublicKey,
			ClientIPv4:    iface.ClientIPv4,
			ClientIPv6:    iface.ClientIPv6,
		},
	}, nil
}

func (g *Generator) DeleteDevice(deviceID, token string) error {
	return g.client.DeleteDevice(deviceID, token)
}
