package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var (
	url_table_name    string = "urls_updated_tags"
	search_table_name string = "search"
	url_table_id      string = "gogiffy.gifs." + url_table_name
	search_table_id   string = "gogiffy.gifs." + search_table_name
	bigqueryClient    *bigquery.Client
	dataset           *bigquery.Dataset
	url_table         *bigquery.Table
	search_table      *bigquery.Table
)

type BigqueryURL struct {
	url              string
	contains_caption bool
	text             string
	tags             []string
	bucket_uid       string
}

type BigquerySearch struct {
	message_id  string
	query_time  time.Time
	token       string
	tags        []string
	index       int
	bucket_uids []string
}

type GifStats struct {
	total_gifs        int
	gifs_without_text int
	gifs_with_text    int
	unique_tags_count int
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

func (m *BigquerySearch) Save() (row map[string]bigquery.Value, insertID string, err error) {
	row = map[string]bigquery.Value{
		"message_id":  m.message_id,
		"query_time":  m.query_time.Format("2006-01-02 15:04:05"),
		"token":       m.token,
		"tags":        m.tags,
		"index":       m.index,
		"bucket_uids": m.bucket_uids,
	}
	return row, "", nil
}

func GetPublicURLFromUID(bucket_uid string) (url string) {
	return "https://storage.googleapis.com/go-giffy-gif-data/" + bucket_uid
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
	url_table = dataset.Table(url_table_name)
	search_table = dataset.Table(search_table_name)
}

func DoesUrlExist(url string) bool {
	query := bigqueryClient.Query(`
		SELECT url
		FROM ` + "`" + url_table_id + "`" + `
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
	FROM ` + "`" + url_table_id + "`" + `
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
	u := url_table.Inserter()
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

func GetUIDFromTags(tags []string) (bqUIDS []string, err error) {
	bqUIDS = []string{}

	query := bigqueryClient.Query(`
	SELECT *
	FROM ` + "`" + url_table_id + "`" + `
	WHERE EXISTS (
		SELECT 1
		FROM UNNEST(tags) AS element
		WHERE element in ('` + strings.Join(tags, "','") + `'))`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}

	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return bqUIDS, err
		}

		bqUIDS = append(bqUIDS, values[4].(string))
	}

	return bqUIDS, nil
}

func convertValuesToStrings(values []bigquery.Value) []string {
	strings := make([]string, len(values))
	for i, v := range values {

		strings[i] = fmt.Sprintf("%v", v)
		// // Assuming the value is a string
		// if str, ok := v.([]byte); ok {
		// 	strings[i] = string(str)
		// } else {

		// 	strings[i] = ""
		// }
	}
	return strings
}

func AddBigQuerySearch(bigQuerySearch BigquerySearch) error {
	bigQuerySearch.query_time = time.Now()
	u := search_table.Inserter()
	items := []*BigquerySearch{
		&bigQuerySearch,
	}
	err := u.Put(ctx, items)
	if err != nil {
		return err
	}
	fmt.Println("Search message stored in BigQuery")
	return nil
}

func GetBigQuerySearch(message_id string) (bigQuerySearch BigquerySearch, err error) {
	bigQuerySearch = BigquerySearch{}

	query := bigqueryClient.Query(`
	SELECT message_id, token, tags, index, bucket_uids
	FROM ` + "`" + search_table_id + "`" + `
	WHERE message_id = "` + message_id + `"
	ORDER BY query_time DESC LIMIT 1;`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}

	var values []bigquery.Value
	err = it.Next(&values)
	if err != nil {
		return bigQuerySearch, err
	}

	bigQuerySearch.message_id = values[0].(string)
	bigQuerySearch.token = values[1].(string)
	bigQuerySearch.tags = convertValuesToStrings(values[2].([]bigquery.Value))
	bigQuerySearch.index = int(values[3].(int64))
	bigQuerySearch.bucket_uids = convertValuesToStrings(values[4].([]bigquery.Value))
	return bigQuerySearch, nil
}

func GetUrlStats() (stats GifStats, err error) {
	stats = GifStats{}
	query := bigqueryClient.Query(`
	SELECT 
    	COUNT(*),
    	SUM(CASE WHEN text = "" THEN 1 ELSE 0 END),
    	SUM(CASE WHEN text != "" THEN 1 ELSE 0 END),
    	(SELECT COUNT(DISTINCT tag)
        	FROM ` + "`" + url_table_id + "`" + `,
     		UNNEST(tags) AS tag)
		FROM ` + "`" + url_table_id + "`" + `;`)

	it, err := query.Read(ctx)
	if err != nil {
		panic(err)
	}

	var values []bigquery.Value
	err = it.Next(&values)
	if err != nil {
		return stats, err
	}

	stats.total_gifs = int(values[0].(int64))
	stats.gifs_without_text = int(values[1].(int64))
	stats.gifs_with_text = int(values[2].(int64))
	stats.unique_tags_count = int(values[3].(int64))

	return stats, nil
}
