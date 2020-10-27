package main

import (
	"franklin/models"

	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/mmcdole/gofeed"
)

const webProtocol string = "http"
const webServerIP string = "127.0.0.1"
const webServerPort int = 8080

var serverURL string = fmt.Sprintf(
	"%v://%v:%d",
	webProtocol,
	webServerIP,
	webServerPort,
)

func main() {
	feeds := getFeedsData()

	for _, feed := range feeds.Data {
		processFeed(feed)
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func processFeed(feed models.Feed) {
	time_now := time.Now()
	if isTimeToSyncFeed(feed.PreviousDownloadTime, feed.DownloadFrequency, time_now) {
		fmt.Println(feed.Source)
		fmt.Println(feed.Title)
		fmt.Println("------------------")

		f := pullFeedFromSource(feed.URL)

		for _, item := range f.Items {
			createNewsItemIfNotPresent(feed.ID, item)
		}

		updateFeed(feed.ID, time_now)
	}
}

func getFeedsData() *models.FeedsGetResponse {
	u, err := url.Parse(serverURL)
	Check(err)

	u.Path = path.Join(
		u.Path,
		"feeds",
	)
	resource := u.String()

	resp, err := http.Get(resource)
	Check(err)

	body, err := ioutil.ReadAll(resp.Body)
	Check(err)

	feeds_resp := new(models.FeedsGetResponse)
	err = json.Unmarshal(body, &feeds_resp)
	Check(err)

	return feeds_resp
}

func isTimeToSyncFeed(prev_download_time time.Time, download_freq int, time_now time.Time) bool {
	if prev_download_time.IsZero() {
		fmt.Println("have not previously downloaded.  downloading...")
		return true
	} else if prev_download_time.AddDate(0, 0, download_freq).Before(time_now) {
		fmt.Println("download frequency expired.  downloading...")
		return true
	} else {
		return false
	}
}

func pullFeedFromSource(feed_url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(feed_url)
	return feed
}

func isNewsItemRegistered(news_item_hash string) bool {
	u, err := url.Parse(serverURL)
	Check(err)

	u.Path = path.Join(
		u.Path,
		"news_items",
		news_item_hash,
	)

	resp, err := http.Get(u.String())
	Check(err)

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

func registerNewsItem(body []uint8) {
	u, err := url.Parse(serverURL)
	Check(err)

	u.Path = path.Join(
		u.Path,
		"news_items",
	)
	request, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewBuffer(body),
	)
	request.Header.Set("Content-type", "application/json")
	Check(err)

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	_, err = client.Do(request)
	Check(err)
}

func createNewsItemIfNotPresent(feed_id uint, item *gofeed.Item) {
	guid_md5_sum := md5.Sum([]byte(item.GUID))
	guid_hash := hex.EncodeToString(guid_md5_sum[:])

	if isNewsItemRegistered(guid_hash) {
		// news item already registered
		return
	}

	var pub_date string
	if item.PublishedParsed != nil {
		pub_date = (*item.PublishedParsed).Format(time.RFC3339)
	}

	requestBody, err := json.Marshal(map[string]string{
		"title":       item.Title,
		"link":        item.Link,
		"description": item.Description,
		"guid":        item.GUID,
		"guid_md5":    guid_hash,
		"pub_date":    pub_date,
		"feed_id":     fmt.Sprint(feed_id),
	})
	Check(err)

	registerNewsItem(requestBody)
}

func updateFeed(feed_id uint, time_now time.Time) {
	requestBody, err := json.Marshal(map[string]string{
		"previous_download_time": time_now.Format(time.RFC3339),
	})
	Check(err)

	u, err := url.Parse(serverURL)
	Check(err)

	u.Path = path.Join(
		u.Path,
		"feeds",
		fmt.Sprint(feed_id),
	)
	request, err := http.NewRequest(
		"PATCH",
		u.String(),
		bytes.NewBuffer(requestBody),
	)
	request.Header.Set("Content-type", "application/json")
	Check(err)

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	_, err = client.Do(request)
	Check(err)
}
