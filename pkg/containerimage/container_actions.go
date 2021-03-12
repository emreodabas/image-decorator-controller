package containerimage

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"strings"
)

func PullImage(src string) (v1.Image, error) {

	pull, err := crane.Pull(src)

	if err != nil {
		return nil, err
	}

	return pull, nil
}

func PushImage(image v1.Image, target string) error {

	err := crane.Push(image, target)

	if err != nil {
		return err
	}

	return nil
}

// clone image from src to target Repo and return image source
func CloneImage(src string, targetRepo string) (string, error) {
	targetImage := ""
	if !strings.HasPrefix(src, targetRepo) {
		image, err := PullImage(src)
		if err != nil {
			//log.Info("Pulling image", "container source", src)
			return "", err
		}
		if strings.Contains(src, "/") {
			targetImage = targetRepo + strings.SplitAfter(src, "/")[1]
		} else {
			targetImage = targetRepo + src
		}
		tag, err := name.NewTag(targetImage)
		if err != nil {
			//log.Info("Tagging image", "%s container tag as %s ", src, targetImage)
			return "", err
		}

		if s, err := daemon.Write(tag, image); err != nil {
			return "", err
		} else {
			fmt.Println(s)
		}

		err = PushImage(image, targetImage)
		if err != nil {
			//log.Info("Tagging image", "%s container tag as %s ", src, targetImage)
			return "", err
		}
		return targetImage, nil
	} else {
		return src, nil
	}

}
