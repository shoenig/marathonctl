package main

// All actions under command artifact

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// upload
type ArtifactUpload struct {
	client *Client
	format Formatter
}

func (a ArtifactUpload) Apply(args []string) {
	Check(len(args) == 2, "must supply path and file")

	path := args[0]
	filename := args[1]

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, e := writer.CreateFormFile("file", filepath.Base(path))
	Check(e == nil, "failed to create form file", e)
	f, e := os.Open(filename)
	Check(e == nil, "failed to open file "+filename, e)
	_, e = io.Copy(part, f)
	Check(e == nil, "failed to get file bytes", e)
	e = writer.Close()
	Check(e == nil, "failed to close file", e)

	request := a.client.POST("/v2/artifacts"+url.QueryEscape(path), body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)

	Check(response.StatusCode == 201, "unable to upload file", stringifyResponse(response))
	defer response.Body.Close()

	fmt.Println(a.format.Format(strings.NewReader(response.Header.Get("Location")), a.Humanize))
}

func (P ArtifactUpload) Humanize(body io.Reader) string {
	b, e := ioutil.ReadAll(body)
	Check(e == nil, "reading upload response failed", e)
	text := "LOCATION\n" + string(b)
	return Columnize(text)
}

// get
type ArtifactGet struct {
	client *Client
	format Formatter
}

func (a ArtifactGet) Apply(args []string) {
	Check(len(args) == 1, "must supply id")
	request := a.client.GET("/v2/artifacts" + url.QueryEscape(args[0]))
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	Check(response.StatusCode != 404, "artifact not found")
	Check(response.StatusCode == 200, "error downloading artifact", stringifyResponse(response))
	defer response.Body.Close()

	b, e := ioutil.ReadAll(response.Body)
	os.Stdout.Write(b)
}

// delete
type ArtifactDelete struct {
	client *Client
	format Formatter
}

func (a ArtifactDelete) Apply(args []string) {
	Check(len(args) == 1, "must supply id")
	request := a.client.DELETE("/v2/artifacts" + url.QueryEscape(args[0]))
	response, e := a.client.Do(request)
	Check(e == nil, "failed to delete artifact", e)
	Check(response.StatusCode == 200, "failed to delete artifact", stringifyResponse(response))
	defer response.Body.Close()

	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a ArtifactDelete) Humanize(body io.Reader) string {
	return "DELETED"
}

func stringifyResponse(res *http.Response) string {
	cnt, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return string(cnt)
}
