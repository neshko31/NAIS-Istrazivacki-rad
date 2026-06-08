package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

func registerWithEureka() {
    eurekaURL := os.Getenv("EUREKA_CLIENT_SERVICEURL_DEFAULTZONE")
    if eurekaURL == "" {
        eurekaURL = "http://localhost:8761/eureka"
    }

    // Hostname mora biti Docker service name, ne localhost!
    hostname := os.Getenv("EUREKA_INSTANCE_HOSTNAME")
    if hostname == "" {
        hostname = "auth-service" // fallback na Docker service name
    }

    port := 8082

    body := map[string]any{
        "instance": map[string]any{
            "instanceId": fmt.Sprintf("%s:%d", hostname, port),
            "hostName":   hostname,   // "auth_service" — Docker DNS
            "app":        "AUTH-SERVICE",
            "ipAddr":     hostname,   
			"vipAddress":  "auth-service", 
        	"secureVipAddress": "auth-service",
            "status":     "UP",
            "port":       map[string]any{"$": port, "@enabled": "true"},
            "securePort": map[string]any{"$": 443, "@enabled": "false"},
            "healthCheckUrl": fmt.Sprintf("http://%s:%d/health", hostname, port),
            "statusPageUrl":  fmt.Sprintf("http://%s:%d/health", hostname, port),
            "homePageUrl":    fmt.Sprintf("http://%s:%d/", hostname, port),
            "dataCenterInfo": map[string]any{
                "@class": "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
                "name":   "MyOwn",
            },
            "leaseInfo": map[string]any{
                "renewalIntervalInSecs": 30,
                "durationInSecs":        90,
            },
        },
    }

    jsonBody, _ := json.Marshal(body)
    url := fmt.Sprintf("%s/apps/AUTH-SERVICE", eurekaURL)

    for i := 0; i < 15; i++ {
        resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
        if err == nil && resp.StatusCode == 204 {
            log.Printf("Registered with Eureka as %s:%d", hostname, port)
            startHeartbeat(eurekaURL, hostname, port)
            return
        }
        log.Printf("Eureka registration attempt %d/15 failed, retrying in 5s...", i+1)
        time.Sleep(5 * time.Second)
    }
    log.Fatal("Could not register with Eureka after 15 attempts")
}

func startHeartbeat(eurekaURL, hostname string, port int) {
    instanceID := fmt.Sprintf("%s:%d", hostname, port)
    heartbeatURL := fmt.Sprintf("%s/apps/AUTH-SERVICE/%s", eurekaURL, instanceID)
    go func() {
        for {
            time.Sleep(30 * time.Second)
            req, err := http.NewRequest("PUT", heartbeatURL, nil)
            if err != nil {
                continue
            }
            resp, err := http.DefaultClient.Do(req)
            if err != nil || resp.StatusCode != 200 {
                log.Printf("Eureka heartbeat failed, re-registering...")
            }
        }
    }()
}