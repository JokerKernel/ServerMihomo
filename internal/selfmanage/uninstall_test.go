package selfmanage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRunDownloadsScriptAndPassesUninstallArgument(t *testing.T) {
	client := testClient(http.StatusOK, "#!/bin/sh\nprintf 'action=%s\\n' \"$1\"\n")

	var output bytes.Buffer
	if err := run(context.Background(), client, "https://example.invalid/install.sh", strings.NewReader(""), &output, &output, "uninstall"); err != nil {
		t.Fatal(err)
	}
	if output.String() != "action=uninstall\n" {
		t.Fatalf("unexpected script output: %q", output.String())
	}
}

func TestRunDownloadsScriptWithoutArgumentForUpdate(t *testing.T) {
	client := testClient(http.StatusOK, "#!/bin/sh\nprintf 'arguments=%s\\n' \"$#\"\n")

	var output bytes.Buffer
	if err := run(context.Background(), client, "https://example.invalid/install.sh", strings.NewReader(""), &output, &output); err != nil {
		t.Fatal(err)
	}
	if output.String() != "arguments=0\n" {
		t.Fatalf("unexpected script output: %q", output.String())
	}
}

func TestRunRejectsUnexpectedContent(t *testing.T) {
	client := testClient(http.StatusOK, "not a shell script")

	err := run(context.Background(), client, "https://example.invalid/install.sh", strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "不是有效") {
		t.Fatalf("expected invalid script error, got: %v", err)
	}
}

func TestRunReportsHTTPFailure(t *testing.T) {
	client := testClient(http.StatusNotFound, "missing")

	err := run(context.Background(), client, "https://example.invalid/install.sh", strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected HTTP error, got: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func testClient(status int, body string) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})}
}
