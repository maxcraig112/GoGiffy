package main

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var (
	bigqueryClient *bigquery.Client
	dataset        *bigquery.Dataset
	table          *bigquery.Table
)

type BigqueryURL struct {
	url              string
	contains_caption bool
	text             string
	tags             []string
	bucket_uid       string
}

func (m *BigqueryURL) Save() (row map[string]bigquery.Value, insertID string, err error) {

	row = map[string]bigquery.Value{
		"url":              m.url,
		"contains_caption": m.contains_caption,
		"text":             m.text,
		"tags":             m.tags,
		"bucket_uid":       m.bucket_uid,
	}
	return row, "", nil
}

func (m *BigqueryURL) GetPublicURL() (url string) {
	return "https://storage.googleapis.com/go-giffy-gif-data/" + m.bucket_uid
}

func init() {
	ctx = context.Background()
	bigqueryClient, err = bigquery.NewClient(ctx, project_ID)
	if err != nil {
		panic(err)
	}

	dataset = bigqueryClient.Dataset("gifs")
	table = dataset.Table("urls")

}

func DoesUrlExist(url string) bool {
	query := bigqueryClient.Query(`
		SELECT url
		FROM ` + "`gogiffy.gifs.urls`" + `
		WHERE url = "` + url + `";`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		return err != iterator.Done
	}
}

func GetBigQueryURL(url string) (bigQueryURL BigqueryURL, err error) {
	bigQueryURL = BigqueryURL{}

	query := bigqueryClient.Query(`
	SELECT *
	FROM ` + "`gogiffy.gifs.urls`" + `
	WHERE url = "` + url + `";`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}

	var values []bigquery.Value
	err = it.Next(&values)
	if err != nil {
		return bigQueryURL, err
	}

	bigQueryURL.url = values[0].(string)
	bigQueryURL.contains_caption = values[1].(bool)
	bigQueryURL.text = values[2].(string)
	bigQueryURL.tags = convertValuesToStrings(values[3].([]bigquery.Value))
	bigQueryURL.bucket_uid = values[4].(string)

	return bigQueryURL, nil

}

func AddBigQueryURL(bigQueryURL BigqueryURL) error {
	u := table.Inserter()
	items := []*BigqueryURL{
		&bigQueryURL,
	}
	err := u.Put(ctx, items)
	if err != nil {
		return err
	}
	fmt.Println("URL stored in BigQuery")
	return nil
}

func GetUrlsFromTag(tags []string) (bqURLs []BigqueryURL, err error) {
	bqURLs = []BigqueryURL{}

	query := bigqueryClient.Query(`
	SELECT *
	FROM ` + "`gogiffy.gifs.urls`" + `
	WHERE EXISTS (
		SELECT 1
		FROM UNNEST(tags) AS element
		WHERE element in ('` + strings.Join(tags, "','") + `'))`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}

	for {
		newBQUrl := BigqueryURL{}
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return bqURLs, err
		}
		newBQUrl.url = values[0].(string)
		newBQUrl.contains_caption = values[1].(bool)
		newBQUrl.text = values[2].(string)
		newBQUrl.tags = convertValuesToStrings(values[3].([]bigquery.Value))
		newBQUrl.bucket_uid = values[4].(string)

		bqURLs = append(bqURLs, newBQUrl)
	}

	return bqURLs, nil
}

func convertValuesToStrings(values []bigquery.Value) []string {
	strings := make([]string, len(values))
	for i, v := range values {
		// Assuming the value is a string
		if str, ok := v.([]byte); ok {
			strings[i] = string(str)
		} else {
			strings[i] = ""
		}
	}
	return strings
}
