package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	// "net"
	"net/http"
	"net/url"

	"coyote/setting"
)

type Client struct {
	API   string
	Email string
}

type Zones struct {
	Result []setting.Zone `json:"result"`
}

type Records struct {
	Result []setting.Record `json:"result"`
}

func createURL(p ...string) string {
	u, _ := url.ParseRequestURI("https://api.cloudflare.com")
	u.Path = "/client/v4/zones/"
	urlStr := fmt.Sprintf("%v", u)
	for _, v := range p {
		urlStr = fmt.Sprintf("%v%v/", urlStr, v)
	}

	return urlStr
}

func (c *Client) sendRequest(method string, urlStr string, data io.Reader) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, urlStr, data)
	if err != nil {
		log.Print(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", c.API)
	req.Header.Set("X-Auth-Email", c.Email)

	resp, reqErr := client.Do(req)

	if reqErr != nil || resp == nil {
		log.Print(reqErr)
		return nil
	}
	log.Printf("	%v	|	%v	|	%v\n", method, resp.Status, urlStr)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

func (c *Client) fetchRecord(z *setting.Zone) {
	u := createURL(z.ID, "dns_records")
	data := url.Values{}
	b := c.sendRequest("GET", u, bytes.NewBufferString(data.Encode()))
	if b == nil {
		return
	}
	var res Records
	json.Unmarshal(b, &res)

	for k1, v1 := range res.Result {
		for k2, v2 := range z.Records {
			if v1.Name == v2.Name {
				z.Records[k2] = res.Result[k1]
			}
		}
	}
}

func (c *Client) FetchAll() bool {
	u := createURL()

	data := url.Values{}
	data.Add("status", "active")
	data.Add("page", "1")
	data.Add("per_page", "20")
	data.Add("order", "status")
	data.Add("direction", "desc")
	data.Add("match", "all")

	b := c.sendRequest("GET", u, bytes.NewBufferString(data.Encode()))
	if b == nil {
		return false
	}
	var res Zones

	json.Unmarshal(b, &res)

	for _, v1 := range res.Result {
		for k, v2 := range setting.Cfg.Zones {
			if v2.Name == v1.Name {
				v2.ID = v1.ID
				c.fetchRecord(&v2)
				setting.Cfg.Zones[k] = v2
			}
		}
	}

	return true
}

func (c *Client) updateRec(zid string, rec *setting.Record, newIP string) {

	u := createURL(zid, "dns_records", rec.ID)
	oldIP := rec.Content
	rec.Content = newIP
	data, _ := json.Marshal(rec)

	b := c.sendRequest("PUT", u, bytes.NewBuffer(data))
	if b != nil {
		log.Printf("%v Changed from %v to %v", rec.Name, oldIP, newIP)
	}
}

func (c *Client) checkIP(ip string) {
	for k1, v1 := range setting.Cfg.Zones {
		for k2, v2 := range v1.Records {
			if v2.Content != ip {
				c.updateRec(v1.ID, &setting.Cfg.Zones[k1].Records[k2], ip)
			}
		}
	}
}

func (c *Client) Run() {
	data := url.Values{}
	b := c.sendRequest("GET", setting.Cfg.IPServer, bytes.NewBufferString(data.Encode()))

	c.checkIP(string(b))
}
