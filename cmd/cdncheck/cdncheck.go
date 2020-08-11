package main

import (
	"fmt"
	"github.com/vkhodor/cdncheck/pkg/cloudconfig"
	"os"
)

var version = "0.0.1"

func main() {
	r53client := cloudconfig.NewCloudRoute53("Z2WXU28CDS7KHT", "content.algorithmic.bid.")
	fmt.Println(r53client.Status())

	if r53client.Status() == "normal" {
		r53client.Fallback()
		os.Exit(0)
	}
	r53client.Normal()
	os.Exit(0)
}
