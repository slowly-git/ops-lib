package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// How many seconds the program should wait before trying to connect to the dashboard again
const RetryTimeout = 5

type grafanaConfig struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Access    string `json:"access"`
	IsDefault bool   `json:"isDefault"`
	URL       string `json:"url"`
	Password  string `json:"password"`
	User      string `json:"user"`
	Database  string `json:"database"`
}

func main() {

	envParams := map[string]string{
		"grafana_user":               "admin",
		"grafana_passwd":             "admin",
		"grafana_port":               "3000",
		"influxdb_host":              "monitoring-influxdb",
		"influxdb_port":              "8086",
		"influxdb_database":          "k8s",
		"influxdb_user":              "root",
		"influxdb_password":          "root",
		"influxdb_service_url":       "",
		"dashboard_location":         "/dashboards",
		"gf_auth_anonymous_enabled":  "true",
		"gf_security_admin_user":     "",
		"gf_security_admin_password": "",
		"gf_server_http_port":        "",
		"gf_server_protocol":         "http",
		"backend_access_mode":        "proxy",
	}

	//如果在环境变量中检测到有envParams中的key,则将envParams中key相关的value替换为环境变量中的值
	for k := range envParams {
		if v := os.Getenv(strings.ToUpper(k)); v != "" {
			envParams[k] = v
		}
	}

	if envParams["influxdb_service_url"] == "" {
		envParams["influxdb_service_url"] = fmt.Sprintf("http://%s:%s", envParams["influxdb_host"], envParams["influxdb_port"])
	}

	cfg := grafanaConfig{
		Name:      "influxdb-datasource",
		Type:      "influxdb",
		Access:    envParams["backend_access_mode"],
		IsDefault: true,
		URL:       envParams["influxdb_service_url"],
		User:      envParams["influxdb_user"],
		Password:  envParams["influxdb_password"],
		Database:  envParams["influxdb_database"],
	}
	// Override setup env vars with Grafana configuration env vars if present
	adminUser := envParams["grafana_user"]
	if user, ok := envParams["gf_security_admin_user"]; ok && len(user) != 0 {
		adminUser = user
	}
	adminPassword := envParams["grafana_passwd"]
	if password, ok := envParams["gf_security_admin_password"]; ok && len(password) != 0 {
		adminPassword = password
	}
	httpPort := envParams["grafana_port"]
	if port, ok := envParams["gf_server_http_port"]; ok && len(port) != 0 {
		httpPort = port
	}

	grafanaURL := fmt.Sprintf("%s://%s:%s@localhost:%s", envParams["gf_server_protocol"], adminUser, adminPassword, httpPort)

	for {
		res, err := http.Get(grafanaURL + "/api/org")
		if err != nil {
			fmt.Printf("Can't access the Grafana dashboard. Error: %v. Retrying after %d seconds...\n", err, RetryTimeout)
			time.Sleep(RetryTimeout * time.Second)
			continue
		}

		_, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("Can't access the Grafana dashboard. Error: %v. Retrying after %d seconds...\n", err, RetryTimeout)
			time.Sleep(RetryTimeout * time.Second)
			continue
		}

		fmt.Println("Connected to the Grafana dashboard.")
		break
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(cfg)

	for {
		_, err := http.Post(grafanaURL+"/api/datasources", "application/json; charset=utf-8", b)
		if err != nil {
			fmt.Printf("Failed to configure the Grafana dashboard. Error: %v. Retrying after %d seconds...\n", err, RetryTimeout)
			time.Sleep(RetryTimeout * time.Second)
			continue
		}

		fmt.Println("The datasource for the Grafana dashboard is now set.")
		break
	}

	dashboardDir := envParams["dashboard_location"]
	files, err := ioutil.ReadDir(dashboardDir)
	if err != nil {
		fmt.Printf("Failed to read the directory the json files should be in. Exiting... Error: %v\n", err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dashboardDir, file.Name())
		jsonbytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Failed to read the json file: %s. Proceeding with the next one. Error: %v\n", filePath, err)
			continue
		}

		_, err = http.Post(grafanaURL+"/api/dashboards/db", "application/json; charset=utf-8", bytes.NewReader(jsonbytes))
		if err != nil {
			fmt.Printf("Failed to post the json file: %s. Proceeding with the next one. Error: %v\n", filePath, err)
			continue
		}
	}
}
