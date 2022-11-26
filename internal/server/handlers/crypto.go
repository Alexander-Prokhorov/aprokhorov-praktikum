package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"aprokhorov-praktikum/internal/ccrypto"
)

func Decrypt(privateKey *ccrypto.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privateKey == nil {
				next.ServeHTTP(w, r)

				return
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			r.Body.Close()

			decrypt, err := privateKey.Decrypt(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(decrypt)))

			next.ServeHTTP(w, r)
		})
	}
}
