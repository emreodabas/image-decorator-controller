package containerimage

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/gcrane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

var (
	logger = log.Log.WithName("container_actions")
)

type ContainerRepository struct {
	RepositoryPath string
	Username       string
	Password       string
	AccessToken    string
}

func (cr *ContainerRepository) pullImage(src string) (v1.Image, error) {

	if strings.Contains(src, "gcr.io") {
		pull, err := crane.Pull(src, crane.WithAuthFromKeychain(gcrane.Keychain))
		if err != nil {
			logger.Error(err, "Error is occur while pulling image %s ", src)
			return nil, err
		}
		return pull, nil
	}
	pull, err := crane.Pull(src, cr.getAuthOptions())
	if err != nil {
		logger.Error(err, "Error is occur while pulling image %s ", src)
		return nil, err
	}

	return pull, nil
}

func (cr *ContainerRepository) pushImage(image v1.Image, target string) {

	err := crane.Push(image, target, cr.getAuthOptions())
	if err != nil {
		logger.Error(err, "Error is occur while pushing image %s", target)
		panic("Unable to push image to backup repository")
	}
}

// clone image from src to target Repo and return image source
func (cr *ContainerRepository) CloneImage(src string) (string, error) {
	targetImage := ""
	image, err := cr.pullImage(src)
	if err != nil {
		return "", err
	}
	if strings.Contains(src, "/") {
		after := strings.SplitAfter(src, "/")
		targetImage = cr.RepositoryPath + after[len(after)-1]
	} else {
		targetImage = cr.RepositoryPath + src
	}
	err = cr.tagImage(image, targetImage)
	if err != nil {
		return "", err
	}
	cr.pushImage(image, targetImage)
	return targetImage, nil
}

func (cr *ContainerRepository) tagImage(image v1.Image, targetImage string) error {

	tag, err := name.NewTag(targetImage)
	if err != nil {
		logger.Error(err, "Error is occur while tagging image to %s", targetImage)
		return err
	}
	daemon.Write(tag, image)
	return nil
}

func (cr *ContainerRepository) getAuthOptions() crane.Option {

	if cr.Username == "" {
		// no auth options
		return crane.WithAuth(authn.Anonymous)
	} else if cr.AccessToken != "" {
		return crane.WithAuth(&authn.Basic{
			Username: cr.Username,
			Password: cr.AccessToken,
		})
	} else {
		if cr.Password != "" {
			return crane.WithAuth(&authn.Basic{
				Username: cr.Username,
				Password: cr.Password,
			})
		} else {
			//no auth option
			return crane.WithAuth(authn.Anonymous)
		}
	}
}
