package scan

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func SubScan(url, wordlist, exclude string, threads int) {

	GetIntro(url, wordlist, exclude, threads)
	fmt.Println(BANNER_SEP)
	list := getList(wordlist)
	if strings.Contains(url, "http") {
		split := strings.Split(url, "//")
		url = split[1]
	}

	for i := 0; i < threads; i++ {
		min := (len(list) / threads) * i
		max := (len(list) / threads) * (i + 1)
		wg.Add(1)
		go Req(url, list[min:max])
	}
	// Wait for goroutines to end
	wg.Wait()

}
func Req(url string, list []string) {

	for _, s := range list {

		wHttpUrl := "http://" + s + "." + url
		wHttpsUrl := "https://" + s + "." + url

		var netClient = &http.Client{
			Timeout: time.Second * 10,
		}
		response, err := netClient.Get(wHttpUrl)
		response2, err := netClient.Get(wHttpsUrl)

		if err == nil {
			fmt.Println("[+] - URL: ", wHttpUrl, " ", response.StatusCode)
			fmt.Println("[+] - URL: ", wHttpsUrl, " ", response2.StatusCode)
		}

	}

}
