package warp

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func GenerateKeyPair() (privateKey string, publicKey string, err error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", err
	}

	return key.String(), key.PublicKey().String(), nil
}
