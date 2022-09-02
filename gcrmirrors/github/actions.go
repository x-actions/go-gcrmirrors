package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Name string         `json:"name"`
	On   On             `json:"on"`
	Jobs map[string]Job `json:"-"`
}

type Schedule struct {
	Cron string `json:"cron"`
}

type Push struct {
	Branches []string `json:"branches"`
}

type On struct {
	Schedule []Schedule `json:"schedule"`
	Push     Push       `json:"push"`
}

type With struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Token      string `json:"token"`
	Path       string `json:"path"`
}

type Steps struct {
	Name string                 `json:"name"`
	Uses string                 `json:"uses"`
	With With                   `json:"with,omitempty"`
	Env  map[string]interface{} `json:"env,omitempty"`
}

type Job struct {
	RunsOn string  `json:"runs-on"`
	Steps  []Steps `json:"steps"`
}

func parseAction(yamlPath string) (*Action, error) {
	in, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	action := &Action{}
	err = yaml.Unmarshal(in, &action)
	if err != nil {
		log.Println(err)
		return action, err
	}

	return action, nil
}

func ScanWorkflows(sourceDir string) ([]*Action, error) {
	workflowsDir := fmt.Sprintf("%s/.github/workflows", sourceDir)
	dirs, err := os.ReadDir(workflowsDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil, nil
	}

	var actions []*Action
	for _, obj := range dirs {
		filePath := path.Join(workflowsDir, obj.Name())
		if obj.IsDir() {
			fmt.Println(obj.Name(), "is Dir, skip.")
		} else {
			action, err := parseAction(filePath)
			if err != nil {
				_ = fmt.Sprintf("parse action %s err %v\n", filePath, err.Error())
			}
			actions = append(actions, action)
		}
	}

	return actions, nil
}

type Mirror struct {
	Name               string `json:"name"`
	SrcRepo            string `json:"src_repo"`
	DestRepo           string `json:"dest_repo"`
	ShortDestRepo      string `json:"short_dest_repo"`
	ImageCount         int    `json:"image_count"`
	SrcImageListURL    string `json:"src_image_list_url"`
	RawSrcImageListURL string `json:"raw_src_image_list_url"`
	ActionName         string `json:"action_name"`
	BadgeURL           string `json:"badge_url"`
	WorkflowURL        string `json:"workflow_url"`
}

type MirrorResponse struct {
	Data []*Mirror `json:"data"`
}

func parseImageFile(imageFile string) (string, int, error) {
	bytes, err := os.ReadFile(imageFile)
	if err != nil {
		fmt.Println("read file", imageFile, err.Error())
		return "", 0, err
	}

	lines := strings.Split(string(bytes), "\n")
	items := strings.Split(lines[0], "/")
	var srcRepo string
	if len(items) == 2 {
		srcRepo = items[0]
	} else if len(items) == 3 {
		srcRepo = strings.Join(items[:2], "/")
	} else if len(items) > 3 {
		srcRepo = strings.Join(items[:len(items)-3], "/")
	}

	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "#") == false {
			count++
		}
	}

	return srcRepo, count, nil
}

func ParseMirrorAction(actions []*Action, sourceDir string) []*Mirror {
	var mirrors []*Mirror
	for _, action := range actions {
		for name, job := range action.Jobs {
			for _, step := range job.Steps {
				if strings.HasPrefix(step.Uses, "x-actions/python3-cisctl") {
					destRepo := step.Env["DEST_REPO"].(string)
					destRepoArray := strings.Split(destRepo, "/")
					shortDestRepo := destRepoArray[len(destRepoArray)-1]
					rawSrcImageListURL := step.Env["SRC_IMAGE_LIST_URL"].(string)
					fileRelativePath := parseTxtRelativePath(rawSrcImageListURL)
					srcRepo, imageCount, _ := parseImageFile(path.Join(sourceDir, fileRelativePath))

					srcImageListURL := fmt.Sprintf(
						"https://github.com/%s/%s/blob/main/%s", step.Env["GIT_ORG"], step.Env["GIT_REPO"],
						fileRelativePath)
					workflowFileName := strings.ReplaceAll(action.Name, "/", "-")
					badgeURL := fmt.Sprintf(
						"https://github.com/x-mirrors/gcr.io/actions/workflows/%s.yml/badge.svg", workflowFileName)
					workflowUrl :=
						fmt.Sprintf("https://github.com/x-mirrors/gcr.io/actions/workflows/%s.yml", workflowFileName)

					mirror := Mirror{
						Name:               name,
						SrcRepo:            srcRepo,
						DestRepo:           destRepo,
						ShortDestRepo:      shortDestRepo,
						SrcImageListURL:    srcImageListURL,
						RawSrcImageListURL: rawSrcImageListURL,
						ImageCount:         imageCount,
						ActionName:         action.Name,
						BadgeURL:           badgeURL,
						WorkflowURL:        workflowUrl,
					}
					mirrors = append(mirrors, &mirror)
				}
			}
		}
	}

	return mirrors
}

type ImageMap struct {
	SrcImage    string `json:"src_image"`
	MirrorImage string `json:"mirror_image"`
}

type ImageMapResponse struct {
	Data []*ImageMap `json:"data"`
}

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

func mirrorImageName(srcImage string, mirror *Mirror) string {
	items := strings.Split(srcImage, "/")
	if len(items) == 2 {
		return fmt.Sprintf("%s/%s", mirror.ShortDestRepo, items[len(items)-1])
	} else if len(items) == 3 {
		return fmt.Sprintf("%s/%s", mirror.ShortDestRepo, items[len(items)-1])
	} else if len(items) > 3 {
		length := len(items)
		return fmt.Sprintf("%s/%s-%s", mirror.ShortDestRepo, items[length-3], items[length-1])
	}

	return ""
}

func ParseSourceImages(mirrors []*Mirror, sourceDir string) []*ImageMap {
	var imageMaps []*ImageMap
	for _, mirror := range mirrors {
		srcImageListFilePath := path.Join(sourceDir, parseTxtRelativePath(mirror.SrcImageListURL))
		srcImageList, _ := readFile(srcImageListFilePath)

		for _, srcImage := range srcImageList {
			if strings.HasPrefix(srcImage, "#") == false && strings.Trim(srcImage, " ") != "" {
				im := ImageMap{
					SrcImage:    srcImage,
					MirrorImage: mirrorImageName(srcImage, mirror),
				}
				imageMaps = append(imageMaps, &im)
			}
		}
	}

	return imageMaps
}
