package main

import (
    "crypto/sha1"
    "encoding/hex"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os/exec"
    "strings"
    "time"
)

type Meeting struct {
    MeetingID string `xml:"meetingID"`
    StartTime int64  `xml:"createTime"`
}

type MeetingsResponse struct {
    Meetings []Meeting `xml:"meetings>meeting"`
}

var MAX_DURATION =  4 * time.Hour
var SLEEP_TIME = 10 * time.Minute

func getBBBConfig() (string, string, error) {
    cmd := exec.Command("bbb-conf", "--secret")
    out, err := cmd.Output()
    if err != nil {
        return "", "", err
    }
    
    lines := strings.Split(string(out), "\n")
    var bbbURL, secret string
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "URL:") {
            bbbURL = strings.TrimSpace(strings.TrimPrefix(line, "URL:"))
        } else if strings.HasPrefix(line, "Secret:") {
            secret = strings.TrimSpace(strings.TrimPrefix(line, "Secret:"))
        }
    }
    
    if bbbURL == "" || secret == "" {
        return "", "", fmt.Errorf("could not parse BBB URL or secret from bbb-conf output")
    }
    
    return bbbURL, secret, nil
}

func checksum(apiCall, query, secret string) string {
    h := sha1.New()
    h.Write([]byte(apiCall + query + secret))
    return hex.EncodeToString(h.Sum(nil))
}

func apiCall(bbbURL, secret, apiCall string, params url.Values) ([]byte, error) {
    query := params.Encode()
    cs := checksum(apiCall, query, secret)
    fullUrl := fmt.Sprintf("%sapi/%s?%s&checksum=%s", bbbURL, apiCall, query, cs)

    resp, err := http.Get(fullUrl)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func cleanServer() {
    fmt.Printf("[%s] Checking local BBB server for long-running meetings...\n", time.Now().Format("15:04:05"))
    bbbURL, secret, err := getBBBConfig()
    if err != nil {
        fmt.Printf("[%s] Error getting BBB config: %v\n", time.Now().Format("15:04:05"), err)
        return
    }
    
    fmt.Printf("[%s] Using BBB URL: %s\n", time.Now().Format("15:04:05"), bbbURL)

    data, err := apiCall(bbbURL, secret, "getMeetings", url.Values{})
    if err != nil {
        fmt.Printf("[%s] Error fetching meetings: %v\n", time.Now().Format("15:04:05"), err)
        return
    }

    var meetingsResp MeetingsResponse
    if err := xml.Unmarshal(data, &meetingsResp); err != nil {
        fmt.Printf("[%s] Error parsing XML: %v\n", time.Now().Format("15:04:05"), err)
        return
    }

    now := time.Now().UnixNano() / int64(time.Millisecond)

	fmt.Println("[%s] Meeting count: %d\n", time.Now().Format("15:04:05"), len(meetingsResp.Meetings))

    for _, m := range meetingsResp.Meetings {
        if m.StartTime == 0 {
            continue
        }

        duration := time.Duration(now-m.StartTime) * time.Millisecond
        if duration >= MAX_DURATION {
            fmt.Printf("[%s] Ending meeting %s (%.2f hours)\n", time.Now().Format("15:04:05"), m.MeetingID, duration.Hours())
            params := url.Values{}
            params.Set("meetingID", m.MeetingID)
            _, err := apiCall(bbbURL, secret, "end", params)
            if err != nil {
                fmt.Printf("[%s] Error ending meeting %s: %v\n", time.Now().Format("15:04:05"), m.MeetingID, err)
            } 
        }
    }
}

func main() {
	fmt.Println("Starting BBB Local Cleaner Service")
	
	fmt.Printf("[%s] Running initial cleanup...\n", time.Now().Format("15:04:05"))
	cleanServer()
	
	ticker := time.NewTicker(SLEEP_TIME)
	defer ticker.Stop()

	for {
		<-ticker.C
		fmt.Printf("[%s] Running scheduled cleanup...\n", time.Now().Format("15:04:05"))
		cleanServer()
	}
}

