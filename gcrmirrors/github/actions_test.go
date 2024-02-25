package github

import (
	"fmt"
	"testing"
)

func TestScanWorkflows(t *testing.T) {
	sourceDir := "/Users/xiexianbin/workspace/code/github.com/kbcx/gcr.io"
	actions, err := ScanWorkflows(sourceDir)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return
	}
	for _, action := range actions {
		t.Logf("%v", action)
	}

	mirrors := ParseMirrorAction(actions, sourceDir)
	for _, mirror := range mirrors {
		fmt.Println("mirror:", mirror)
	}

	imageMaps := ParseSourceImages(mirrors, sourceDir)
	fmt.Println("imageMaps:")
	for _, imageMap := range imageMaps {
		fmt.Println("imageMap:", imageMap)
	}
}
