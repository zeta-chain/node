package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	port = ":8000"

	basePath             = "/opt/my-local-ton/myLocalTon"
	liteClientConfigPath = basePath + "/genesis/db/my-ton-local.config.json"
	settingsPath         = basePath + "/settings.json"

	faucetJSONKey = "faucetWalletSettings"
)

func main() {
	http.HandleFunc("/faucet.json", errorWrapper(faucetHandler))
	http.HandleFunc("/lite-client.json", errorWrapper(liteClientHandler))
	http.HandleFunc("/status", errorWrapper(statusHandler))

	//nolint:gosec
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func errorWrapper(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			errResponse(w, http.StatusInternalServerError, err)
		}
	}
}

// Handler for the /faucet.json route
func faucetHandler(w http.ResponseWriter, _ *http.Request) error {
	faucet, err := extractFaucetFromSettings(settingsPath)
	if err != nil {
		return err
	}

	jsonResponse(w, http.StatusOK, faucet)
	return nil
}

// liteClientHandler returns lite json client config
// and alters localhost to docker IP if needed.
func liteClientHandler(w http.ResponseWriter, _ *http.Request) error {
	data, err := os.ReadFile(liteClientConfigPath)
	if err != nil {
		return fmt.Errorf("could not read lite client config: %w", err)
	}

	dockerIP := os.Getenv("DOCKER_IP")

	if dockerIP != "" {
		altered, err := alterConfigIP(data, dockerIP)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, err)
			return nil
		}

		data = altered
	}

	jsonResponse(w, http.StatusOK, json.RawMessage(data))
	return nil
}

// Handler for the /status route
func statusHandler(w http.ResponseWriter, _ *http.Request) error {
	if _, err := os.Stat(liteClientConfigPath); err != nil {
		return fmt.Errorf("lite client config %q not found: %w", liteClientConfigPath, err)
	}

	faucet, err := extractFaucetFromSettings(settingsPath)
	if err != nil {
		return err
	}

	type faucetShape struct {
		Created bool `json:"created"`
	}

	var fs faucetShape
	if err = json.Unmarshal(faucet, &fs); err != nil {
		return fmt.Errorf("failed to parse faucet settings: %w", err)
	}

	if !fs.Created {
		return errors.New("faucet is not created yet")
	}

	jsonResponse(w, http.StatusOK, map[string]string{"status": "OK"})
	return nil
}

func extractFaucetFromSettings(filePath string) (json.RawMessage, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read faucet settings: %w", err)
	}

	var keyValue map[string]json.RawMessage
	if err := json.Unmarshal(data, &keyValue); err != nil {
		return nil, fmt.Errorf("failed to parse faucet settings: %w", err)
	}

	faucet, ok := keyValue[faucetJSONKey]
	if !ok {
		return nil, errors.New("faucet settings not found in JSON")
	}

	return faucet, nil
}

func errResponse(w http.ResponseWriter, status int, err error) {
	jsonResponse(w, status, map[string]string{"error": err.Error()})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		b = []byte("Failed to marshal JSON")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//nolint:errcheck
	w.Write(b)
}

// TON's lite client config contains the IP of the node.
// And it's localhost, we need to change it to the docker IP.
func alterConfigIP(config []byte, ipString string) ([]byte, error) {
	const localhost = uint32(2130706433)

	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP: %q", ipString)
	}

	return bytes.ReplaceAll(
		config,
		uint32ToBytes(localhost),
		uint32ToBytes(ip2int(ip)),
	), nil
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func uint32ToBytes(n uint32) []byte {
	return []byte(fmt.Sprintf("%d", n))
}
