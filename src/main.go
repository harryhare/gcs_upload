package main

import (
	"fmt"
	"google.golang.org/api/iterator"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
	"context"
	"cloud.google.com/go/storage"
	"log"
	"google.golang.org/api/option"
	"os"
	"net/http"
	"time"
	"io"
	"io/ioutil"
	"bytes"
	"strconv"
	"encoding/json"
	"gitlab.internal.unity3d.com/unity-connect/connect/server/shared"
	"strings"
	"regexp"
)

const project_id = "gcs-test-230118"
const bucket = "harryhare"

// explicit reads credentials from the specified path.
func explicit(jsonPath, projectID string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Buckets:")
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}

// implicit uses Application Default Credentials to authenticate.
func implicit() {
	ctx := context.Background()

	// For API packages whose import path is starting with "cloud.google.com/go",
	// such as cloud.google.com/go/storage in this case, if there are no credentials
	// provided, the client library will look for credentials in the environment.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	it := storageClient.Buckets(ctx, "gcs-test-230118")
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bucketAttrs.Name)
	}

	// For packages whose import path is starting with "google.golang.org/api",
	// such as google.golang.org/api/cloudkms/v1, use the
	// golang.org/x/oauth2/google package as shown below.
	oauthClient, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	kmsService, err := cloudkms.New(oauthClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(kmsService.Projects)
}

func GetPublicObject() {
	//bucket:="harryhare"
	url := "https://www.googleapis.com/storage/v1/b/harryhare/o/kitten.png"
	//storage.SignedURL("")
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}

func GetPrivateObject() {
	//bucket:="harryhare"
	url := "https://www.googleapis.com/storage/v1/b/harryhare/o/kitten2.png"

	token := GetToken(storage.ScopeReadOnly)
	//storage.SignedURL("")
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", "Bearer", token.AccessToken))
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}

type FileMeta struct {
	Name string `json:"name"`
}

func CreateResumableUpload(length int64) string {
	url := "https://www.googleapis.com/upload/storage/v1/b/harryhare/o?uploadType=resumable&name=test1.png"
	token := GetToken(storage.ScopeReadWrite)

	request, err := http.NewRequest(http.MethodPost, url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", "Bearer", token.AccessToken))
	request.Header.Add("X-Upload-Content-Type", "image/png")
	request.Header.Add("X-Upload-Content-Length", strconv.FormatInt(length, 10))
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Header.Get("Location"))
	io.Copy(os.Stdout, response.Body)
	return response.Header.Get("Location")
}

func CreateResumableUploadWithBody(length int64) string {

	url := "https://www.googleapis.com/upload/storage/v1/b/harryhare/o?uploadType=resumable&name=test/test24.jpg"
	token := GetToken(storage.ScopeReadWrite)

	meta := &FileMeta{
		Name: "test/test24.jpg",
	}
	shared.DumpJson(meta)
	bodyBytes, err := json.Marshal(meta)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(bodyBytes)

	request, err := http.NewRequest(http.MethodPost, url, body)
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", "Bearer", token.AccessToken))
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))
	request.Header.Add("X-Upload-Content-Type", "image/jpg")
	request.Header.Add("X-Upload-Content-Length", strconv.FormatInt(length, 10))
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Minute}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Header.Get("Location"))
	fmt.Printf("header:%s\n",response.Header)
	io.Copy(os.Stdout, response.Body)
	return response.Header.Get("Location")
}

var rangeReg = regexp.MustCompile(`bytes=(\d+)-(\d+)`)

func getRange(str string) (int64, int64) {
	if str == "" {
		return 0, 0
	}
	substrings := rangeReg.FindStringSubmatch(str)
	if len(substrings) != 3 {
		return 0, 0
	}
	begin, _ := strconv.ParseInt(substrings[1], 10, 64)
	end, _ := strconv.ParseInt(substrings[2], 10, 64)
	return begin, end + 1
}

func GetResumableStatus(url string) int64 {
	request, err := http.NewRequest(http.MethodPut, url, nil)
	request.Header.Add("Content-Length", "0")
	request.Header.Add("Content-Range", "bytes */*")
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Minute}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	shared.DumpJson(response.Header)
	fmt.Println(response.Header.Get("Content-Length"))
	fmt.Println(response.Header.Get("Range"))
	io.Copy(os.Stdout, response.Body)
	_, end := getRange(response.Header.Get("Range"))
	return end
}

func DeleteResumbelUpload(url string) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	request.Header.Add("Content-Length", "0")
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Minute}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	shared.DumpJson(response.Header)
	io.Copy(os.Stdout, response.Body)
}

func PutResumableUpload(url string, start, end, total int64, reader io.Reader) {
	//token := GetToken(storage.ScopeReadWrite)

	request, err := http.NewRequest(http.MethodPut, url, reader)
	//request.Header.Add("Authorization", fmt.Sprintf("%s %s", "Bearer", token.AccessToken))// don't need auth
	request.Header.Add("Content-Length", strconv.FormatInt(end-start, 10))
	request.Header.Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end-1, total))
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	client := http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Header.Get("Location"))
	io.Copy(os.Stdout, response.Body)
}

func SimpleUpload() {
	filePath := "/Users/unity/git/gcs_upload/kitten.png"
	url := "https://www.googleapis.com/upload/storage/v1/b/harryhare/o?uploadType=media&name=test/mykitten3.png"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(b)
	request, err := http.NewRequest(http.MethodPost, url, reader)
	token := GetToken(storage.ScopeReadWrite)
	if err != nil {
		fmt.Println("new quest error")
		panic(err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", "Bearer", token.AccessToken))
	request.Header.Add("Content-Type", "image/png")
	request.Header.Add("Content-Length", strconv.Itoa(len(b)))
	client := http.Client{Timeout: time.Hour}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("client.do error")
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}

func DeleteTest() {
	url := CreateResumableUploadWithBody(100000)
	GetResumableStatus(url)
	DeleteResumbelUpload(url)
}

func UploadTest() {
	filePath := "/Users/unity/git/gcs_upload/test.jpg"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	url := CreateResumableUploadWithBody(int64(len(b)))
	GetResumableStatus(url)
	l := int64(len(b))
	reader := bytes.NewReader(b)
	PutResumableUpload(url, 0, 262144, l, &io.LimitedReader{reader, 262144})
	start := GetResumableStatus(url)
	PutResumableUpload(url, start, l, l, reader)
}

func limitReaderTest() {
	str := "1234567890"
	reader := strings.NewReader(str)
	io.Copy(os.Stdout, &io.LimitedReader{reader, 3})
	io.Copy(os.Stdout, &io.LimitedReader{reader, 3})
	io.Copy(os.Stdout, &io.LimitedReader{reader, 3})
	io.Copy(os.Stdout, reader)
}

func intTest() {
	var x int
	x = 1 << 32
	fmt.Println(x)
	fmt.Println(1 << 32)
	x = 1 << 62
	fmt.Println(x)
	fmt.Println(1 << 62)
}

func getRangeTest() {
	fmt.Println(getRange("bytes=0-42"))
	fmt.Println(getRange("bytes=0-42/1234"))
}

func getTokenTest() {
	pre_token := ""
	c := 0
	for i := 0; i < 100; i++ {
		token := GetToken(storage.ScopeReadOnly).AccessToken
		if token != pre_token {
			fmt.Printf("%s %d\n", pre_token, c)
			c = 1
			pre_token = token
		}
		c++
	}
	fmt.Printf("%s %d\n", pre_token, c)
}
func main() {

	// todo try
	// storage.SignedURL()

	fmt.Println(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	//implicit()
	//explicit("/Users/unity/.gcp/gcs/gcs.json","gcs-test-230118")
	//SimpleUpload()
	//GetPrivateObject()
	//GetPublicObject()
	UploadTest()
	//DeleteTest()
}
