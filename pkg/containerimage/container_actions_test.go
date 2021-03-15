package containerimage

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

var (
	publicImageNames = []string{"nginx", "busybox", "alpine"}
	privateImageName = "kubermatico/private-image"
	publicGCRImage   = "gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0"
)

func TestCloneImage(t *testing.T) {

	containerRepository := getTestContainerRepository()
	// creating random tags of public image
	src := getRandomSourceOfPublicImage()
	cloneImage, err := containerRepository.CloneImage(src)
	if err != nil {
		t.Errorf("Expected no error returns err: %v", err)
	}
	if cloneImage != containerRepository.RepositoryPath+src {
		t.Errorf("Expected image is changed to %v but found %v", containerRepository.RepositoryPath+src, cloneImage)
	}
	image, err := getTooWrongTestContainerRepositoryWithBasicAuth().CloneImage(cloneImage)
	if err == nil {
		t.Errorf("Expected error but returns no error")
	}
	if image != "" {
		t.Errorf("Expected no image but return non empty image")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic situation that is not triggered")
		}
	}()
	getWrongTestContainerRepositoryWithBasicAuth().CloneImage(cloneImage)

}

func TestCloneImageWithWrongSource(t *testing.T) {

	containerRepository := getTestContainerRepository()
	// creating non avaliable source
	src := "nginx:" + strconv.Itoa(rand.Int())
	_, err := containerRepository.CloneImage(src)
	if err == nil {
		t.Errorf("Expected error but returns no error")
	}
}

func TestPullImageForPublic(t *testing.T) {
	src := getRandomSourceOfPublicImage()
	image, err := getTestContainerRepository().pullImage(src)
	if err != nil {
		t.Errorf("No error is expected but found %v", err)
	}
	if image == nil {
		t.Errorf("Expected non nil image but return nil image")
	}
}

func TestPullImageForGCRPublic(t *testing.T) {
	image, err := getTestContainerRepository().pullImage(publicGCRImage)
	if err != nil {
		t.Errorf("No error is expected but found %v", err)
	}
	if image == nil {
		t.Errorf("Expected non nil image but return nil image")
	}
}

func TestPullPrivateImageWithUnauthorizedUser(t *testing.T) {
	image, err := getEmptyContainerRepository().pullImage(privateImageName)
	if err == nil {
		t.Errorf("Error is expected but not found ")
	}
	if image != nil {
		t.Errorf("Expected nil image but return non nil image")
	}
}
func TestPullPrivateImageWithAuthorizedUser(t *testing.T) {
	src := getRandomSourceOfPublicImage()
	image, err := getTestContainerRepository().pullImage(src)
	if err != nil {
		t.Errorf("No error is expected but found %v", err)
	}
	if image == nil {
		t.Errorf("Expected non nil image but return nil image")
	}
}

func TestPushImage(t *testing.T) {
	src := getRandomSourceOfPublicImage()
	containerRepository := getTestContainerRepository()
	image, err := containerRepository.pullImage(src)
	validPushPath := containerRepository.RepositoryPath + "test-image:" + strconv.Itoa(rand.Int())
	err = containerRepository.tagImage(image, validPushPath)
	if err != nil {
		t.Errorf("No error is expected but found %v", err)
	}
	containerRepository.pushImage(image, validPushPath)
	pullImage, err := containerRepository.pullImage(validPushPath)

	if err != nil {
		t.Errorf("No error is expected but found %v", err)
	}
	if pullImage == nil {
		t.Errorf("Expected non nil image but return nil image")
	}
}

func TestPushImageWithUnauthorizedUser(t *testing.T) {
	src := getRandomSourceOfPublicImage()
	containerRepository := getTestContainerRepository()
	fmt.Println(src, "is pulling")
	image, err := containerRepository.pullImage(src)
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
	validPushPath := containerRepository.RepositoryPath + "test-image:" + strconv.Itoa(rand.Int())
	err = containerRepository.tagImage(image, validPushPath)
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic situation that is not triggered")
		}
	}()
	getEmptyContainerRepository().pushImage(image, validPushPath)

}

func TestTagImage(t *testing.T) {
	containerRepository := getTestContainerRepository()
	newTag := "tag" + strconv.Itoa(rand.Int())
	image, err := random.Image(1024, 5)
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
	err = containerRepository.tagImage(image, newTag)
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
}

func TestTagImageWithBadTag(t *testing.T) {
	containerRepository := getTestContainerRepository()
	image, err := random.Image(1024, 5)
	if err != nil {
		t.Errorf("No error expected but found %v", err)
	}
	newTag := "tag  with space"
	err = containerRepository.tagImage(image, newTag)
	if err == nil {
		t.Errorf("Error is expected but not found ")
	}
	newTag = "tagwithb@d'char"
	err = containerRepository.tagImage(image, newTag)
	if err == nil {
		t.Errorf("Error is expected but not found ")
	}

}

func TestAuthWithUsernameAndPwd(t *testing.T) {
	auth := getTestContainerRepositoryWithBasicAuth()
	options := auth.getAuthOptions()
	if options == nil {
		t.Errorf("Options is expected but not found")
	}

	tags, err := crane.ListTags(privateImageName, options)
	if err != nil {
		t.Errorf("Error is expected but not found ")
	}
	if tags == nil {
		t.Errorf("Expected tags but not found")
	}
}

func TestAuthWithAnonymous(t *testing.T) {
	auth := getEmptyContainerRepository()
	options := auth.getAuthOptions()
	if options == nil {
		t.Errorf("Options is expected but not found")
	}
	tags, err := crane.ListTags(privateImageName, options)
	if err == nil {
		t.Errorf("Error is not expected but found ")
	}
	if tags != nil {
		t.Errorf("Expected tags could not fetched but it exist")
	}
}

func TestAuthWithoutCredentials(t *testing.T) {
	auth := getContainerRepositoryWithoutCredentials()
	options := auth.getAuthOptions()
	if options == nil {
		t.Errorf("Options is expected but not found")
	}
	tags, err := crane.ListTags(privateImageName, options)
	if err == nil {
		t.Errorf("Error is not expected but found ")
	}
	if tags != nil {
		t.Errorf("Expected tags could not fetched but it exist")
	}
}

// could return unsupported MediaType that trigger below exception
// unsupported MediaType: "application/vnd.docker.distribution.manifest.v1+prettyjws", see https://github.com/google/go-containerregistry/issues/377
//func getRandomSourceOfPublicImage() string {
//	publicImage := publicImageNames[rand.Int()%len(publicImageNames)]
//	tags, _ := crane.ListTags(publicImage)
//	return publicImage + ":" + tags[rand.Int()%len(tags)]
//}

func getTestContainerRepository() *ContainerRepository {
	return &ContainerRepository{
		RepositoryPath: "kubermatico/",
		Username:       "kubermatico",
		Password:       "$pus@L&3G?!ewTK",
		AccessToken:    "a5c4c4a4-27f8-40af-9662-a0391bee1d6d",
	}
}

func getTestContainerRepositoryWithBasicAuth() *ContainerRepository {
	return &ContainerRepository{
		RepositoryPath: "kubermatico/",
		Username:       "kubermatico",
		Password:       "$pus@L&3G?!ewTK",
	}
}
func getWrongTestContainerRepositoryWithBasicAuth() *ContainerRepository {
	return &ContainerRepository{
		RepositoryPath: "kubermaticonot/",
		Username:       "kubermatico",
		Password:       "$pus@L&3G?!ewTK",
	}
}
func getTooWrongTestContainerRepositoryWithBasicAuth() *ContainerRepository {
	return &ContainerRepository{
		RepositoryPath: "kuberm@tico'/",
		Username:       "kubermatico",
		Password:       "$pus@L&3G?!ewTK",
	}
}

func getEmptyContainerRepository() *ContainerRepository {
	return &ContainerRepository{}
}
func getContainerRepositoryWithoutCredentials() *ContainerRepository {
	return &ContainerRepository{
		Username: "asda",
	}
}

func getRandomSourceOfPublicImage() string {
	tags := new(Tags)
	publicImageName := publicImageNames[rand.Int()%len(publicImageNames)]
	resp, err := http.Get("https://registry.hub.docker.com/v2/repositories/library/" + publicImageName + "/tags?page_size=100")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&tags)
	return publicImageName + ":" + tags.Results[rand.Int()%len(tags.Results)].Name
}

type Tags struct {
	Results []Result `json:"results,omitempty"`
	Count   int      `json:"count,omitempty"`
}

type Result struct {
	Name string `json:"name,omitempty"`
}
