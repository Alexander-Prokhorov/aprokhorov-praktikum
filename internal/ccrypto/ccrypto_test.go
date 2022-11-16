package ccrypto_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/ccrypto"
)

func Test_GenKeysForTest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "First Test",
			wantErr: false,
		},
		{
			name:    "Second Test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// test
			keys, err := ccrypto.NewKeyPair()
			assert.NoError(t, err)
			crypt, err := rsa.EncryptPKCS1v15(rand.Reader, keys.GetPub(), []byte(tt.name))
			assert.NoError(t, err)
			decrypt, err := rsa.DecryptPKCS1v15(rand.Reader, keys.GetPriv(), crypt)
			assert.NoError(t, err)

			assert.Equal(t, string(decrypt), tt.name)
		})
	}
}
