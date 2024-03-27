package main

import (
	"context"
	"fmt"
	gf "image/gif"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
	bucketHandle  *storage.BucketHandle

	bucketName string = "go-giffy-gif-data"
)

func init() {
	ctx = context.Background()
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	bucketHandle = storageClient.Bucket(bucketName)
}

func StoreGifInBucket(gif Gif) error {
	obj := bucketHandle.Object(gif.UID)

	w := obj.NewWriter(ctx)

	err := gf.EncodeAll(w, gif.data)
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	gif.bucketURI = fmt.Sprintf("gs://%s/%s", bucketName, gif.UID)

	fmt.Println("Object " + gif.UID + " stored in bucket")
	return nil
}

func GetGifFromBucket(UID string) (gif Gif, err error) {
	gif = Gif{}
	obj := bucketHandle.Object(UID)
	r, err := obj.NewReader(ctx)
	if err != nil {
		fmt.Println("Error with creating bucket reader")
		return gif, err
	}

	gif.data, err = gf.DecodeAll(r)
	if err != nil {
		fmt.Println("Error with decoding gif")
		return gif, err
	}

	if err := r.Close(); err != nil {
		return gif, err
	}

	fmt.Println("Object " + UID + " retrieved in bucket")
	gif.bucketURI = fmt.Sprintf("gs://%s/%s", bucketName, UID)
	gif.UID = UID
	return gif, nil
}
