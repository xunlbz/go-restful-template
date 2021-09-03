package lib

import (
	"fmt"
	"regexp"
	"testing"
)

func TestGetContainerList(t *testing.T) {
	cli := NewDockerClient()
	cli.GetContainerList()
}

func TestReg(t *testing.T) {
	sourceStr := `/a/1.js`
	matched, _ := regexp.MatchString(`(/.+)*\.(html|js)`, sourceStr)
	fmt.Printf("%v\n", matched) // true
}
