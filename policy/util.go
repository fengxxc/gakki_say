package policy

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ImageType int

const (
	Unknown ImageType = iota
	Png
	Jpg
	Jpeg
	Gif
)

func getImg(fileName string) (image.Image, ImageType, error) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return nil, Unknown, err
	}
	defer f.Close()
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var img image.Image
	var imgType ImageType
	// 先用最土的方式判断类型
	switch fileExt {
	case ".jpg":
		imgType = Jpg
		img, err = jpeg.Decode(f)
	case ".jpeg":
		imgType = Jpeg
		img, err = jpeg.Decode(f)
	case ".png":
		imgType = Png
		img, err = png.Decode(f)
	case ".gif":
		imgType = Gif
		img, err = gif.Decode(f)
	default:
		log.Println("Unsupported file format: " + fileExt)
	}
	if err != nil {
		log.Println(err)
		return nil, Unknown, err
	}
	return img, imgType, nil
}

func imgToBytes(img image.Image, imgType ImageType) []byte {
	buf := new(bytes.Buffer)
	switch imgType {
	case Jpg:
		jpeg.Encode(buf, img, nil)
	case Jpeg:
		jpeg.Encode(buf, img, nil)
	case Png:
		png.Encode(buf, img)
	case Gif:
		gif.Encode(buf, img, nil)
	}
	return buf.Bytes()
}
