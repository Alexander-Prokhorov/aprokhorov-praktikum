package ccrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

const size = 512

type PrivateKey struct {
	key *rsa.PrivateKey
}

func NewPrivateKeyFromFile(filepath string) (*PrivateKey, error) {
	rawKey, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(rawKey)

	key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)

	return &PrivateKey{key: key}, err
}

func (pk *PrivateKey) GetKey() *rsa.PrivateKey {
	return pk.key
}

func (pk *PrivateKey) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk.key),
	})

	if _, err = file.Write(pubBytes); err != nil {
		return err
	}

	return file.Close()
}

func (pk *PrivateKey) Decrypt(msg []byte) ([]byte, error) {
	msgLen := len(msg)
	step := 512
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptPKCS1v15(rand.Reader, pk.key, msg[start:finish])
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

type PublicKey struct {
	key *rsa.PublicKey
}

func NewPublicKeyFromFile(filepath string) (*PublicKey, error) {
	rawKey, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(rawKey)

	key, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)

	return &PublicKey{key: key}, err
}

func (pk *PublicKey) GetKey() *rsa.PublicKey {
	return pk.key
}

func (pk *PublicKey) WriteToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pk.key),
	})

	if _, err = file.Write(pubBytes); err != nil {
		return err
	}

	return file.Close()
}

func (pk *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	msgLen := len(msg)
	step := size - 11
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pk.key, msg[start:finish])
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

type KeyPair struct {
	priv *PrivateKey
	pub  *PublicKey
}

func NewKeyPair() (*KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, size*8)
	return &KeyPair{
		priv: &PrivateKey{key: key},
		pub:  &PublicKey{key: &key.PublicKey},
	}, err
}

func (k *KeyPair) GetPriv() *rsa.PrivateKey {
	return k.priv.key
}

func (k *KeyPair) GetPub() *rsa.PublicKey {
	return k.pub.key
}

func (k *KeyPair) WriteKeyToFile() error {
	if err := k.pub.WriteToFile("id_rsa.pub"); err != nil {
		return err
	}
	return k.priv.WriteToFile("id_rsa")
}
