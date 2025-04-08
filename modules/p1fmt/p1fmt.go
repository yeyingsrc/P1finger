package p1fmt

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

func PrintfShoot(prefix, targetUrl, webTitle string, fingerTag []string, StatusCode int) {
	titleWidth := 50
	fingerTagStr := strings.Join(fingerTag, ", ")
	paddedWebTitle := runewidth.FillRight(webTitle, titleWidth)

	fmt.Printf("%v %-60v [%-3v] Title: %s \033[31m%v\033[0m\n", prefix, targetUrl, StatusCode, paddedWebTitle, fingerTagStr)
}

func PrintMiss(prefix, targetUrl, webTitle string, fingerTag []string, StatusCode int) {
	titleWidth := 50
	fingerTagStr := strings.Join(fingerTag, ", ")
	paddedWebTitle := runewidth.FillRight(webTitle, titleWidth)

	fmt.Printf("%v %-60v [%-3v] Title: %s \033[31m%v\033[0m\n", prefix, targetUrl, StatusCode, paddedWebTitle, fingerTagStr)
}

func PrintReqFail(prefix, targetUrl string, fingerTag []string) {
	fingerTagStr := strings.Join(fingerTag, ", ")
	fmt.Printf("%v %-45v \u001B[31m%v\u001B[0m\n", prefix, targetUrl, fingerTagStr)
}
