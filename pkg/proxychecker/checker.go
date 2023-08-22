package proxychecker

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/MidasVanVeen/proxy-checker/pkg/proxyutils"
)

type Checker struct {
	protocol proxyutils.Protocol
	timeout  time.Duration
	retries  int
	checkUrl *url.URL
}

func NewChecker(protocol proxyutils.Protocol, timeout time.Duration, retries int) (*Checker, error) {
	url, err := url.Parse("http://httpbin.org/ip")
	if err != nil {
		return nil, err
	}
	return &Checker{
		protocol: protocol,
		timeout:  timeout,
		retries:  retries,
		checkUrl: url,
	}, nil
}

func (c *Checker) Check(proxy string, ua string, ref string) bool {
	if !proxyutils.ValidateProxy(proxy) || proxy == "" {
		return false
	}
	retries := 0
	transport, err := proxyutils.ConfigureTransport(c.protocol, proxy)
	if err != nil {
		println("error: " + err.Error())
		return false
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   c.timeout,
	}
	var reader io.ReadCloser
	var jsonResp httpbinResponse

	for retries < c.retries {
		req, err := http.NewRequest("GET", c.checkUrl.String(), reader)
		if err != nil {
			retries++
			continue
		}
		req.Header.Add("user-agent", ua)
		req.Header.Add("referer", ref)
		req.Header.Add("accept", "*/*")
		req.Header.Add("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			retries++
			continue
		}
		if resp.StatusCode != 200 {
			retries++
			continue
		}
		err = json.NewDecoder(resp.Body).Decode(&jsonResp)
		if err != nil {
			retries++
			continue
		}

		if jsonResp.Origin == strings.Split(proxy, ":")[0] {
			return true
		}
		retries++
	}
	return false
}

func (c *Checker) CleanList(proxies []string, uas *[]string, refs *[]string, validCallback func(proxy string)) []string {
	var wg sync.WaitGroup
	var cleanProxies []string
	cleanChannel := make(chan string)
	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			if c.Check(proxy, (*uas)[rand.Intn(len(*uas))], (*refs)[rand.Intn(len(*refs))]) {
				cleanChannel <- proxy
				validCallback(proxy)
			}
		}(proxy)
	}
	go func() {
		wg.Wait()
		close(cleanChannel)
	}()
	for p := range cleanChannel {
		cleanProxies = append(cleanProxies, p)
	}
	return cleanProxies
}
