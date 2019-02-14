package main

import (
	"golang.org/x/oauth2"
	"os"
	"io/ioutil"
	"golang.org/x/oauth2/google"
	"context"
)

func GetToken(scope string) *oauth2.Token {
	ctx := context.Background()
	//google.CredentialsFromJSON(ct)
	//scopes := []string{"https://www.googleapis.com/auth/devstorage.read_only"}
	//credentials,err:=google.FindDefaultCredentials(ctx,scopes...)
	//if err!=nil{
	//	panic(err)
	//}
	file, err := os.Open("/Users/unity/.gcp/gcs/gcs.json")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	credentials, err := google.CredentialsFromJSON(ctx, b, scope)
	token, err := credentials.TokenSource.Token()
	if err != nil {
		panic(err)
	}
	//fmt.Println(token)
	//fmt.Println(token.AccessToken)
	//fmt.Println(token.TokenType)
	return token
}
