package helper

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func getResponseBody(res *http.Response) []byte {

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func ExecRequest(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

func GetRequest(ctx context.Context, uri string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	return req
}

func GetCtxWithTimout(t time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), t)

	return ctx, cancel
}

func GetResponseApi(url string, timout time.Duration) []byte {

	ctx, cancel := GetCtxWithTimout(timout)
	defer cancel()

	res := ExecRequest(GetRequest(ctx, url))
	defer res.Body.Close()

	return getResponseBody(res)

}
