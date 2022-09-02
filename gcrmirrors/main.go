package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gcrmirrors/github"
)

var publicDir string
var sourceDir string

const (
	MirrorRepos  = "mirror-repos.json"
	MirrorImages = "mirror-images.json"
)

func init() {
	flag.StringVar(&publicDir, "publicDir", "./public", "json public dir")
	flag.StringVar(&sourceDir, "sourceDir", "", "https://github.com/kbcx/gcr.io dir")

	flag.Parse()
}

func usage() {
	flag.Usage()
	os.Exit(-1)
}

func initPublicDir() {
	// init result
	if publicDir == "./public" {
		_exePath, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			return
		}
		exePath, _ := filepath.EvalSymlinks(filepath.Dir(_exePath))
		publicDir = path.Join(exePath, publicDir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if _, err := os.Stat(publicDir); err != nil {
		_ = os.MkdirAll(publicDir, 0777)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error

	initPublicDir()
	mirrorReposPath := path.Join(publicDir, MirrorRepos)
	mirrorImagesPath := path.Join(publicDir, MirrorImages)

	actions, err := github.ScanWorkflows(sourceDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mirrors := github.ParseMirrorAction(actions, sourceDir)
	mirrorsResponse, _ := json.Marshal(github.MirrorResponse{
		Data: mirrors,
	})
	err = os.WriteFile(mirrorReposPath, mirrorsResponse, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("write", string(mirrorsResponse), "to file", mirrorReposPath, "success.")

	imageMaps := github.ParseSourceImages(mirrors, sourceDir)
	imageMapsResponse, _ := json.Marshal(github.ImageMapResponse{
		Data: imageMaps,
	})
	err = os.WriteFile(mirrorImagesPath, imageMapsResponse, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("write", string(imageMapsResponse), "to file", mirrorReposPath, "success.")
}
