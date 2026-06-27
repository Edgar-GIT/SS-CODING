package transporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"ss-coding/utils"
)

const ngrokAPI = "http://127.0.0.1:4040/api/tunnels"

type tunnelResponse struct {
	Tunnels []struct {
		PublicURL string `json:"public_url"`
		Proto     string `json:"proto"`
	} `json:"tunnels"`
}

var (
	mu        sync.Mutex
	ngrokCmd  *exec.Cmd
	publicURL string
)

func Installed() bool {
	_, err := exec.LookPath("ngrok")
	return err == nil
}

func Running() bool {
	mu.Lock()
	defer mu.Unlock()
	return ngrokCmd != nil && ngrokCmd.Process != nil
}

func PublicURL() string {
	mu.Lock()
	defer mu.Unlock()
	return publicURL
}

func Start(port int) (string, error) {
	mu.Lock()
	if ngrokCmd != nil && publicURL != "" {
		url := publicURL
		mu.Unlock()
		return url, nil
	}
	mu.Unlock()

	if !Installed() {
		return "", fmt.Errorf("ngrok not found — install from https://ngrok.com/download")
	}

	cmd := exec.Command("ngrok", "http", fmt.Sprint(port), "--log=stdout")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	utils.PrepareProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return "", err
	}

	url, err := waitForPublicURL(20 * time.Second)
	if err != nil {
		_ = utils.KillProcessTree(cmd)
		return "", err
	}

	mu.Lock()
	ngrokCmd = cmd
	publicURL = url
	mu.Unlock()

	go func() {
		_ = cmd.Wait()
		mu.Lock()
		if ngrokCmd == cmd {
			ngrokCmd = nil
			publicURL = ""
		}
		mu.Unlock()
	}()

	return url, nil
}

func Stop() error {
	mu.Lock()
	cmd := ngrokCmd
	mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("no ngrok tunnel running")
	}

	if err := utils.KillProcessTree(cmd); err != nil {
		return err
	}

	mu.Lock()
	if ngrokCmd == cmd {
		ngrokCmd = nil
		publicURL = ""
	}
	mu.Unlock()

	return nil
}

func waitForPublicURL(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		url, err := fetchPublicURL(client)
		if err == nil && url != "" {
			return url, nil
		}
		time.Sleep(400 * time.Millisecond)
	}

	return "", fmt.Errorf("ngrok tunnel did not become ready in time")
}

func fetchPublicURL(client *http.Client) (string, error) {
	resp, err := client.Get(ngrokAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data tunnelResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	for _, tunnel := range data.Tunnels {
		if tunnel.Proto == "https" {
			return tunnel.PublicURL, nil
		}
	}
	if len(data.Tunnels) > 0 {
		return data.Tunnels[0].PublicURL, nil
	}

	return "", fmt.Errorf("no tunnels in ngrok API response")
}
