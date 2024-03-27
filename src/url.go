package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func UrlIsGif(url string) bool {
	//checks if the suffix is .gif or if it has tenor in it
	return !strings.HasPrefix(url, "https://storage.googleapis.com/go-giffy-gif-data/") && (strings.HasSuffix(url, ".gif") || strings.Contains(url, "tenor"))
}

func convertTenorUrl(url string) (convertedUrl string, err error) {
	// Fetch HTML content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read HTML content
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Define regular expression pattern
	regexPattern := `src="https://media1\.tenor\.com/([^"]+?\.gif)"`

	// Compile regular expression
	re := regexp.MustCompile(regexPattern)

	// Find matches in HTML content
	match := re.FindStringSubmatch(string(htmlContent))

	// Extract text from the first match
	if len(match) > 1 {
		return `https://media1.tenor.com/` + match[1], nil
	} else {
		return "", errors.New("tenor gif can't be converted")
	}
}

func convertUrl(url string) (convertedUrl string, err error) {
	//Check if the gif is a tenor gif, and of the form
	if strings.HasPrefix(url, "https://tenor.com/view/") {
		convertedUrl, err = convertTenorUrl(url)
		if err != nil {
			return "", nil
		}

	} else {
		convertedUrl = url
	}
	return convertedUrl, nil
}

func ProcessUrls(urls []string) (gif []Gif, err error) {
	gifList := []Gif{}
	for i := 0; i < len(urls); i++ {
		gif := Gif{}
		urls[i], err = convertUrl(urls[i])
		if !DoesUrlExist(urls[i]) {
			if err != nil {
				return gifList, err
			}
			fmt.Println(urls[i] + " does not exist")

			gif, err = CreateGifFromURL(urls[i])
			if err != nil {
				fmt.Println(err.Error())
				return gifList, err
			}

			if err = StoreGifInBucket(gif); err != nil {
				return gifList, err
			}

			text, err := GetTextFromImageFromBucket(gif)
			if err != nil {
				return gifList, err
			}

			bqURL := BigqueryURL{
				url:              urls[i],
				contains_caption: gif.IsCaptionGif(),
				text:             text,
				tags:             strings.Fields(text),
				bucket_uid:       gif.UID,
			}

			if err = AddBigQueryURL(bqURL); err != nil {
				return gifList, err
			}
			gifList = append(gifList, gif)

		} else {
			fmt.Println(urls[i] + " DOES EXIST")
			bqURL, err := GetBigQueryURL(urls[i])
			if err != nil {
				fmt.Println("error with getting BigQueryURL")
				return gifList, err
			}

			gif, err = GetGifFromBucket(bqURL.bucket_uid)
			if err != nil {
				fmt.Println("error with getting Gif from bucket")
				return gifList, err
			}

			gifList = append(gifList, gif)
		}

	}

	return gifList, nil
}
