package main

import (
	"log"

	"aprokhorov-praktikum/internal/ccrypto"
	"aprokhorov-praktikum/internal/logger"
)

const logLevel = 0

func main() {
	// Init Logger
	logger, err := logger.NewLogger("geKeys.log", logLevel)
	if err != nil {
		log.Fatal("cannot initialize zap.logger")
	}

	logger.Info("Generating key pair ...")

	keys, err := ccrypto.NewKeyPair()
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Debug("Save priv and pub keys to files...")
	if err = keys.WriteKeyToFile(); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Debug("Keys Saved!")
}
