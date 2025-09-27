package main

import (
	"log"
	"os"

	"github.com/fernet/fernet-go"
	"gopkg.in/yaml.v2"
)

type Sumaconf struct {
	Server   string
	User     string
	Password string
	Email_to []string
}

type SUMAConfig struct {
	SUMA map[string]struct {
		User                 string   `yaml:"username"`
		Password             string   `yaml:"password"`
		Logfile              string   `yaml:"logfile"`
		Email_to             []string `yaml:"email_to"`
		Healthcheck_interval int      `yaml:"healthcheck_interval"`
		Healthcheck_email_to []string `yaml:"healthcheck_email"`
	} `yaml:"suma_api"`
}

func GetConfig(file string) *SUMAConfig {
	// Read the file
	data, err := os.ReadFile(file)
	if err != nil {
		logger.Infoln(err)
	}

	// Create a struct to hold the YAML data
	var config SUMAConfig

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Infoln(err)

	}

	key := os.Getenv("SUMAKEY")
	if len(key) == 0 {
		log.Default().Printf("SUMAKEY is not set. This might cause error for password decryption.")
	}

	return &config
}

func Decrypt(key string, cryptoText string) string {
	k := fernet.MustDecodeKeys(key)
	/* tok, err := fernet.EncryptAndSign([]byte(cryptoText), k[0])
	if err != nil {
		panic(err)
	} */
	msg := fernet.VerifyAndDecrypt([]byte(cryptoText), 0, k)
	//logger.Infoln(string(msg))

	return string(msg)
}
