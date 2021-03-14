package containerimage

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

// TODO variable
var (
	logger = log.Log.WithName("container_actions")
)

func pullImage(src string) (v1.Image, error) {
	pull, err := crane.Pull(src)
	if err != nil {
		logger.Error(err, "Error is occur while pulling image %s ", src)
		return nil, err
	}

	return pull, nil
}

func pushImage(image v1.Image, target string) error {
	err := crane.Push(image, target, getAuthOptions())
	if err != nil {
		logger.Error(err, "Error is occur while pushing image %s", target)
		return err
	}

	return nil
}

// clone image from src to target Repo and return image source
func CloneImage(src string, targetRepo string) (string, error) {
	targetImage := ""
	image, err := pullImage(src)
	if err != nil {
		return "", err
	}
	if strings.Contains(src, "/") {
		targetImage = targetRepo + strings.SplitAfter(src, "/")[1]
	} else {
		targetImage = targetRepo + src
	}
	tag, err := name.NewTag(targetImage)
	if err != nil {
		logger.Error(err, "Error is occur while tagging image %s to %s", src, targetImage)
		return "", err
	}

	if _, err := daemon.Write(tag, image); err != nil {
		logger.Error(err, "Error is occur while writing  ta %s to image ", tag)
		return "", err
	}
	err = pushImage(image, targetImage)
	if err != nil {
		return "", err
	}
	return targetImage, nil
}

func getAuthOptions() crane.Option {
	accessToken := viper.GetString("ACCESS_TOKEN")
	user := viper.GetString("USERNAME")
	if user == "" {
		panic("docker user is not defined")
	} else if accessToken != "" {
		return crane.WithAuth(&authn.Basic{
			Username: user,
			Password: accessToken,
		})
	} else {
		pwd := viper.GetString("PASSWORD")
		if pwd != "" {
			return crane.WithAuth(&authn.Basic{
				Username: user,
				Password: pwd,
			})
		} else {
			panic("docker login credentials is not defined")
		}
	}
}
