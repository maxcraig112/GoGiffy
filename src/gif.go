package main

import (
	"fmt"
	gf "image/gif"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Gif struct {
	data      *gf.GIF
	bucketURI string
	UID       string
}

func (m *Gif) GetPublicURL() (url string) {
	return "https://storage.googleapis.com/go-giffy-gif-data/" + m.UID
}

func CreateGifFromURL(url string) (gif Gif, err error) {
	//Get the raw GET data from the URL
	gif = Gif{}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("here")
		return gif, err
	}
	defer response.Body.Close()

	fmt.Println("☁️ Gif downloaded from URL")
	//use the image gif package to decode it into a data type
	gif.data, err = gf.DecodeAll(response.Body)
	if err != nil {
		fmt.Println("there")
		return gif, err
	}

	gif.UID = uuid.New().String() + ".gif"
	gif.bucketURI = fmt.Sprintf("gs://%s/%s", bucketName, gif.UID)
	fmt.Println("✍🏻 Gif decoded into variable")
	return
}

func (gif *Gif) SaveGifAsFileFromData(fileName string) (err error) {
	//Create file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//encode the file will the gif data type
	err = gf.EncodeAll(file, gif.data)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	fmt.Println("👨‍💻 File opened and gif saved")
	return nil
}

func (gif *Gif) IsCaptionGif() bool {
	firstFrame := gif.data.Image[0]
	colour := firstFrame.At(10, 10)
	r, g, b, _ := colour.RGBA()
	minVal := uint32(230)

	return (255 <= r && r <= minVal) && (255 <= g && g <= minVal) && (255 <= b && b <= minVal)

}
