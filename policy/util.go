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

	"github.com/fogleman/gg"
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

func getImgTypeByFileName(fileName string) ImageType {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var imgType ImageType
	// 先用最土的方式判断类型
	switch fileExt {
	case ".jpg":
		imgType = Jpg
	case ".jpeg":
		imgType = Jpeg
	case ".png":
		imgType = Png
	case ".gif":
		imgType = Gif
	default:
		imgType = Unknown
	}
	return imgType
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

func imgWriteText(fileName string, text string) (image.Image, error) {
	// img, err := gg.LoadPNG(fileName)
	img, err := gg.LoadImage(fileName)
	if err != nil {
		log.Printf("Error loading %s", fileName)
		return nil, err
	}
	size := img.Bounds().Size()
	dc := gg.NewContext(size.X, size.Y)
	dc.DrawImage(img, 0, 0)

	err = dc.LoadFontFace("./font/simhei.ttf", 120)
	if err != nil {
		log.Printf("Error loading font face %s", "simhei.ttf")
		return nil, err
	}

	textWidth, textHeight := dc.MeasureString(text)

	// 文字底色
	dc.DrawRectangle(float64(size.X)/2-textWidth/2, float64(size.Y)/2-textHeight/2, textWidth, textHeight+20)
	dc.SetRGBA(1, 0.8, 1, 0.5)
	dc.Fill()
	// dc.DrawString(text, (float64(size.X)-textWidth)/2, (float64(size.Y))/2)
	dc.SetRGB(getGgStyleRGB(255, 255, 255))
	dc.DrawStringAnchored(text, float64(size.X)/2, float64(size.Y)/2, 0.5, 0.5)

	// for test
	dc.SetRGB(getGgStyleRGB(0, 255, 0))
	dc.DrawCircle(float64(size.X/2), float64(size.Y/2), 5)
	dc.Fill()
	return dc.Image(), nil
}

func getGgStyleRGB(r int, g int, b int) (float64, float64, float64) {
	return float64(r / 255), float64(g / 255), float64(b / 255)
}
