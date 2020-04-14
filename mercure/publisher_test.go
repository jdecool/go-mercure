package mercure

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPublisher(t *testing.T) {
	u, _ := url.Parse("http://localhost:3001/.well-known/mercure")
	jwt := "jwt-token"

	p := NewPublisher(u, jwt, nil)

	assert.IsType(t, &http.Client{}, p.client)
}

func TestNewPublisherWithCustomHttpClient(t *testing.T) {
	u, _ := url.Parse("http://localhost:3001/.well-known/mercure")
	jwt := "jwt-token"
	client := &http.Client{}

	p := NewPublisher(u, jwt, client)

	assert.IsType(t, &http.Client{}, p.client)
	assert.Equal(t, client, p.client)
}

func TestPublishMinimalUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "go-mercure/dev", req.Header.Get("User-Agent"))
		assert.Equal(t, "Bearer jwt-token", req.Header.Get("Authorization"))
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		assert.Equal(t, "52", req.Header.Get("Content-Length"))

		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, "data=Hello+World&topic=http%3A%2F%2Flocalhost%2Ftest", string(b))

		rw.Write([]byte("uuid"))
	}))
	defer server.Close()

	u := Update{
		Topics: []string{
			"http://localhost/test",
		},
		Data: []byte("Hello World"),
	}
	url, _ := url.Parse(server.URL)

	p := NewPublisher(url, "jwt-token", server.Client())
	r, err := p.Publish(u)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", r)
}

func TestPublishCompleteUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "go-mercure/dev", req.Header.Get("User-Agent"))
		assert.Equal(t, "Bearer jwt-token", req.Header.Get("Authorization"))
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		assert.Equal(t, "93", req.Header.Get("Content-Length"))

		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, "data=Hello+World&id=id&retry=%05&target=target1&topic=http%3A%2F%2Flocalhost%2Ftest&type=type", string(b))

		rw.Write([]byte("uuid"))
	}))
	defer server.Close()

	u := Update{
		Topics: []string{
			"http://localhost/test",
		},
		Data: []byte("Hello World"),
		Targets: []string{
			"target1",
		},
		Id:    "id",
		Type:  "type",
		Retry: 5,
	}
	url, _ := url.Parse(server.URL)

	p := NewPublisher(url, "jwt-token", server.Client())
	r, err := p.Publish(u)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", r)
}

func TestPublishUpdateWithMultipleTopics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "go-mercure/dev", req.Header.Get("User-Agent"))
		assert.Equal(t, "Bearer jwt-token", req.Header.Get("Authorization"))
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		assert.Equal(t, "127", req.Header.Get("Content-Length"))

		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, "data=Hello+World&topic=http%3A%2F%2Flocalhost%2Ftest1&topic=http%3A%2F%2Flocalhost%2Ftest2&topic=http%3A%2F%2Flocalhost%2Ftest3", string(b))

		rw.Write([]byte("uuid"))
	}))
	defer server.Close()

	u := Update{
		Topics: []string{
			"http://localhost/test1",
			"http://localhost/test2",
			"http://localhost/test3",
		},
		Data: []byte("Hello World"),
	}
	url, _ := url.Parse(server.URL)

	p := NewPublisher(url, "jwt-token", server.Client())
	r, err := p.Publish(u)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", r)
}

func TestPublishUpdateWithMultipleTargets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "go-mercure/dev", req.Header.Get("User-Agent"))
		assert.Equal(t, "Bearer jwt-token", req.Header.Get("Authorization"))
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
		assert.Equal(t, "97", req.Header.Get("Content-Length"))

		b, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, "data=Hello+World&target=target1&target=target2&target=target3&topic=http%3A%2F%2Flocalhost%2Ftest", string(b))

		rw.Write([]byte("uuid"))
	}))
	defer server.Close()

	u := Update{
		Topics: []string{
			"http://localhost/test",
		},
		Data: []byte("Hello World"),
		Targets: []string{
			"target1",
			"target2",
			"target3",
		},
	}
	url, _ := url.Parse(server.URL)

	p := NewPublisher(url, "jwt-token", server.Client())
	r, err := p.Publish(u)

	assert.Nil(t, err)
	assert.Equal(t, "uuid", r)
}

func TestErrorOccuredWhenPublishUpdateWithoutTopic(t *testing.T) {
	u := Update{
		Data: []byte("Hello World"),
	}
	url, _ := url.Parse("http://localhost")

	p := NewPublisher(url, "jwt-token", nil)
	r, err := p.Publish(u)

	assert.Empty(t, r)
	assert.Error(t, err)
	assert.Equal(t, errors.New("Missing topic"), err)
}

func TestErrorOccuredWhenPublishUpdateWithoutData(t *testing.T) {
	u := Update{
		Topics: []string{
			"http://localhost/test",
		},
	}
	url, _ := url.Parse("http://localhost")

	p := NewPublisher(url, "jwt-token", nil)
	r, err := p.Publish(u)

	assert.Empty(t, r)
	assert.Error(t, err)
	assert.Equal(t, errors.New("Missing data"), err)
}
