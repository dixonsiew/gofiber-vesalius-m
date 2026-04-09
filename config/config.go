package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type WSVesaliusConfig struct {
    ServerBaseUrl string
    NricWithDash  string
}

func Config(key string) string {
    // load .env file
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Print("Error loading .env file")
    }
    // Return the value of the variable
    return os.Getenv(key)
}

func GetPatientDocumentCode() string {
    return Config("patient.document.code")
}

func GetIpayTestEnv() string {
    return Config("payment.ipay.testenv")
}

func GetWSVesaliusConfig() WSVesaliusConfig {
    return WSVesaliusConfig{
        ServerBaseUrl: Config("ws.vesalius.server-baseurl"),
        NricWithDash:  Config("ws.vesalius.nric.withDash"),
    }
}
