package cmd

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

const officialWebsiteUrl = "https://dddplayer.com"

func open(raw string) error {
	encoded := encodeURIComponent(raw)
	err := openBrowser(fmt.Sprintf("%s/#%s", officialWebsiteUrl, encoded))
	if err != nil {
		return err
	}
	return nil
}

func encodeURIComponent(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func openBrowser(url string) error {
	cmd := exec.Command("open", url)
	return cmd.Start()
}
