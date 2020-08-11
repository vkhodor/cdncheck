package main

import (
	"fmt"
	"github.com/vkhodor/cdncheck/pkg/cloudconfig"
)

var version = "0.0.1"

func main() {
	r53client := cloudconfig.NewRoute53()
	fmt.Println(r53client)
}
