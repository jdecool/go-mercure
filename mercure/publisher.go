package mercure

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Publisher struct {
	url    *url.URL
	jwt    string
	client *http.Client
}

type Update struct {
	Topics  []string
	Data    []byte
	Targets []string
	Id      string
	Type    string
	Retry   int
}

func NewPublisher(url *url.URL, jwt string, httpClient *http.Client) Publisher {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	p := Publisher{
		url,
		jwt,
		httpClient,
	}

	return p
}

func (c *Publisher) Publish(u Update) (string, error) {
	d, err := getData(u)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", c.url.String(), strings.NewReader(d.Encode()))
	if err != nil {
		return "", err
	}

	r.Header.Set("User-Agent", "go-mercure/dev")
	r.Header.Set("Authorization", "Bearer "+c.jwt)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Content-Length", strconv.Itoa(len(d.Encode())))

	h := &http.Client{}
	resp, err := h.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getData(u Update) (url.Values, error) {
	d := url.Values{}

	if len(u.Topics) == 0 {
		return d, errors.New("Missing topic")
	}

	if len(u.Data) == 0 {
		return d, errors.New("Missing data")
	}

	for _, topic := range u.Topics {
		d.Add("topic", topic)
	}

	d.Add("data", string(u.Data))

	for _, target := range u.Targets {
		d.Add("target", target)
	}

	if u.Id != "" {
		d.Add("id", u.Id)
	}

	if u.Type != "" {
		d.Add("type", u.Type)
	}

	if u.Retry > 0 {
		d.Add("retry", string(u.Retry))
	}

	return d, nil
}
