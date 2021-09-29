package scan

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

const BANNER_SEP = "________________________________________________________________________________\n"

func DirScan(url, wordlist, exclude string, threads int) {

	if url == "" {
		fmt.Println(url)

	}
	GetIntro(url, wordlist, exclude, threads)
	fmt.Println(BANNER_SEP)
	list := getList(wordlist)

	for i := 0; i < threads; i++ {
		min := (len(list) / threads) * i
		max := (len(list) / threads) * (i + 1)
		wg.Add(1)
		go GetOutput(list[min:max], url, exclude)
	}
	// Wait for goroutines to end
	wg.Wait()

}
func GetOutput(list []string, url, exclude string) {
	defer wg.Done()
	for i := 0; i < len(list); i++ {
		newUrl := url + "/" + list[i]
		resp, err := http.Get(newUrl)

		if err != nil {
			color.Red("The directory path in the wordlist is incorrect.", err)
		}
		split := strings.Split(exclude, ",")
		val := 0
		// Filter Excludes
		for _, s := range split {
			if strconv.Itoa(resp.StatusCode) == s {
				val = 1
				break
			}
			val = 0
		}
		// Get Output
		if val == 0 {
			if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
				color.Red("[*] %v %s : %s is not present\n", resp.StatusCode, http.StatusText(resp.StatusCode), list[i])
			} else if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				color.Green("[+] %v %s : %s", resp.StatusCode, http.StatusText(resp.StatusCode), list[i])
			} else if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
				color.Magenta("[-] %v %s: %s respond internal server error\n", resp.StatusCode, http.StatusText(resp.StatusCode), list[i])
			} else if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
				color.Cyan("[*] %v %s: %s redirection \n", resp.StatusCode, http.StatusText(resp.StatusCode), list[i])
			}
		}
	}
}

// Get Directory From Wordlist
func getList(wordlist string) []string {
	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var listlines []string
	for scanner.Scan() {
		listlines = append(listlines, scanner.Text())
	}
	file.Close()
	return listlines
}

// Get Intro
func GetIntro(url, wordlist, exclude string, threads int) {
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}
	r, err := http.Get(url)
	if err != nil {
		log.Fatal("[-] No Such Host! ")
	}
	if r.StatusCode == 200 {
		fmt.Printf(" :: Connected Target : ")
		color.Green("%s", url)
	}
	fmt.Printf(" :: Selected Wordlist : ")
	color.Yellow("%s\n", wordlist)
	fmt.Printf(" :: Threads : ")
	color.Red("%d\n", threads)
	fmt.Printf(" :: Excluded Status Codes : ")
	color.Red("%s\n\n", exclude)
}
