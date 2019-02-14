package main

import (
	"io/ioutil"
	"gitlab.internal.unity3d.com/unity-connect/connect/server/shared/cloud/gcs_restful"
	"bytes"
	"context"
	"fmt"
	"gitlab.internal.unity3d.com/unity-connect/connect/server/shared"
	"io"
)

func testToken() {
	ctx := context.Background()
	token := gcs_restful.GetToken(ctx)
	fmt.Println(token.AccessToken)
	return
}

func testCreateSimpleUpload() {
	ctx := context.Background()

	filePath := "/Users/unity/git/gcs_upload/kitten.png"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	upload := &gcs_restful.ResumableUpload{
		FileKey: "kitten101.png",
		Length:  int64(len(b)),
		Private: false,
		Meta: &gcs_restful.FileMeta{
			Name: "kitten101.png",
		},
		Reader: bytes.NewReader(b),
	}

	err = gcs_restful.SimpleUpload(ctx, upload)
	if err != nil {
		panic(err)
	}
}

func testCreateResumableUpload() {
	filePath := "/Users/unity/git/gcs_upload/kitten.png"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	upload := &gcs_restful.ResumableUpload{
		FileKey: "kitten102.png",
		Length:  int64(len(b)),
		Private: false,
		Meta: &gcs_restful.FileMeta{
			Name: "kitten102.png",
		},
		Reader: bytes.NewReader(b),
	}

	upload, err = gcs_restful.CreateResumableUpload(ctx, upload)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", upload)
}

func testResumableUpload() {
	filePath := "/Users/unity/git/gcs_upload/test.jpg"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	upload := &gcs_restful.ResumableUpload{
		FileKey: "test102.jpg",
		Length:  int64(len(b)),
		Private: false,
		Meta: &gcs_restful.FileMeta{
			Name: "test102.jpg",
		},
		Reader: bytes.NewReader(b),
	}

	upload, err = gcs_restful.CreateResumableUpload(ctx, upload)
	if err != nil {
		panic(err)
	}
	shared.DumpJson(upload)

	url := upload.Location
	start, err := gcs_restful.GetResumableStatus(ctx, url)
	fmt.Printf("status:%d,%v\n", start, err)

	err = gcs_restful.PutResumableUpload(ctx, url, 0, 262144, upload.Length, &io.LimitedReader{upload.Reader, 262144})
	if err != nil {
		panic(err)
	}

	start, err = gcs_restful.GetResumableStatus(ctx, url)
	fmt.Printf("status:%d,%v\n", start, err)
	if err != nil {
		panic(err)
	}

	err = gcs_restful.PutResumableUpload(ctx, url, start, upload.Length, upload.Length, upload.Reader)
	shared.DumpJson(err)
	if err != nil {
		panic(err)
	}

	err = gcs_restful.MakeObjectPublic(ctx, upload.FileKey)
	if err != nil {
		panic(err)
	}
}

func testDeleteResumableUpload() {
	filePath := "/Users/unity/git/gcs_upload/kitten.png"
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	upload := &gcs_restful.ResumableUpload{
		FileKey: "kitten102.png",
		Length:  int64(len(b)),
		Private: false,
		Meta: &gcs_restful.FileMeta{
			Name: "kitten102.png",
		},
		Reader: bytes.NewReader(b),
	}

	upload, err = gcs_restful.CreateResumableUpload(ctx, upload)
	if err != nil {
		panic(err)
	}
	shared.DumpJson(upload)

	url := upload.Location
	err = gcs_restful.DeleteResumbleUpload(ctx, url)
	shared.DumpJson(err)
}

func testACL() {
	ctx := context.Background()
	key := "mykitten.png"
	acl := &gcs_restful.ObjectACL{
		Entity: "allUsers",
		Role:   "READER",
	}
	err := gcs_restful.PostObjectACL(ctx, key, acl)
	if err != nil {
		panic(err)
	}
}
func main() {
	//testToken()
	testResumableUpload()
	//testDeleteResumableUpload()
	//testACL()

	fmt.Println()
}
