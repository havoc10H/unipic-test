package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sysdevguru/unipic/cache"
	"github.com/sysdevguru/unipic/config"

	"github.com/gomodule/redigo/redis"
)

// Picture represents today's picture information
type Picture struct {
	Date        string `json:"date"`
	URL         string `json:"url"`
	CopyRight   string `json:"copyright"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	HDURL       string `json:"hdurl"`
}

// GetPic returns picture information
func (p *Picture) GetPic() error {
	// get cache connection from pool
	conn, err := cache.GetConn()
	if err != nil {
		log.Printf("unipic: failed to get connection to cache: %v\n", err)
		// get response from Nasa
		return p.fetchPicture(p.Date, conn)
	}

	// get from cache
	s, err := cache.Get(conn, p.Date)
	if err != nil {
		log.Printf("unipic: failed to get from cache: %v\n", err)
		// get response from Nasa
		return p.fetchPicture(p.Date, conn)
	}

	// unmarshal response from cache
	err = json.Unmarshal([]byte(s), &p)
	if err != nil {
		return err
	}

	// close redis connection
	conn.Close()

	return nil
}

// fetchPicture fetches picture from NASA service
func (p *Picture) fetchPicture(date string, conn redis.Conn) error {
	defer conn.Close()
	client := &http.Client{}

	// request to Nasa api
	req, err := http.NewRequest("GET", config.Global.Config.NasaURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("date", date)
	q.Add("api_key", os.Getenv("NASA_API_KEY"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// unmarshal response from Nasa
	err = json.Unmarshal(respBody, &p)
	if err != nil {
		return err
	}

	// store in cache
	err = cache.Set(conn, date, string(respBody))
	if err != nil {
		log.Printf("unipic: failed to store in cache: %v\n", err)
	}

	return nil
}
