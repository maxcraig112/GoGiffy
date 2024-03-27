package main

import (
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

var (
	visionClient *vision.ImageAnnotatorClient
)

func init() {
	// apiKey, err := GetToken("visionAPIKey.txt")
	// if err != nil {
	// 	return
	// }

	// visionClient, err = vision.NewImageAnnotatorClient(ctx, option.WithAPIKey(apiKey))
	// if err != nil {
	// 	return
	// }

	visionClient, err = vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return
	}
}

func GetTextFromImageFromBucket(gif Gif) (text string, err error) {

	image := vision.NewImageFromURI(gif.bucketURI)
	annotations, err := visionClient.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return text, err
	}

	// Print detected labels
	if len(annotations) > 0 {
		return strings.ReplaceAll(annotations[0].Description, "\n", " "), nil
	}
	return "", nil
}
