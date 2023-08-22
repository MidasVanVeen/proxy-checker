package main

import (
	"bufio"
	"os"
	"strconv"
	"time"

	"github.com/MidasVanVeen/proxy-checker/pkg/proxychecker"
	"github.com/MidasVanVeen/proxy-checker/pkg/proxyutils"
)

var (
	uas = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (Windows NT 6.1; rv:52.0) Gecko/20100101 Firefox/52.0",
	}
	refs = []string{
		"https://www.google.com/",
		"https://www.bing.com/",
		"https://duckduckgo.com/",
	}
)

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return proxies
}

func main() {
	checker, err := proxychecker.NewChecker(proxyutils.SOCKS5, 3*time.Second, 3)
	if err != nil {
		panic(err)
	}
	proxies := readFile("proxies.txt")
	cleanproxies := checker.CleanList(proxies, &uas, &refs)
	println("Checked " + strconv.Itoa(len(proxies)) + " proxies. " + strconv.Itoa(len(cleanproxies)) + " are working.")
}
