package policy

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
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

type RGBA struct {
	R int
	G int
	B int
	A int
}

func (rgba *RGBA) getGgStyleRGBA() (float64, float64, float64, float64) {
	return GetGgStyleRGBA(rgba.R, rgba.G, rgba.B, rgba.A)
}

type DrawStringConfig struct {
	ax          float64
	ay          float64
	fontFamily  string
	textBgColor *RGBA
}

func imgWriteText(fileName string, text string, drawStringConfig DrawStringConfig) (image.Image, error) {
	// img, err := gg.LoadPNG(fileName)
	img, err := gg.LoadImage(fileName)
	if err != nil {
		log.Printf("Error loading %s", fileName)
		return nil, err
	}
	size := img.Bounds().Size()
	dc := gg.NewContext(size.X, size.Y)
	dc.DrawImage(img, 0, 0)

	fontFamily := "simhei.ttf"
	if drawStringConfig.fontFamily != "" {
		fontFamily = drawStringConfig.fontFamily
	}
	err = dc.LoadFontFace("./font/"+fontFamily, 100)
	if err != nil {
		log.Printf("Error loading font face %s", "simhei.ttf")
		return nil, err
	}

	textWidth, textHeight := dc.MeasureString(text)

	// 文字底色
	if drawStringConfig.textBgColor != nil {
		dc.DrawRectangle(float64(size.X)/2-textWidth/2, float64(size.Y)/10*8-textHeight/2, textWidth+20, textHeight+40)
		// dc.SetRGBA(1, 0.8, 1, 0.5)
		dc.SetRGBA(drawStringConfig.textBgColor.getGgStyleRGBA())
		dc.Fill()
	}

	// 写文字
	dc.SetRGB(GetGgStyleRGB(255, 255, 255))
	dc.DrawStringAnchored(text, float64(size.X)/2, float64(size.Y)/10*8, drawStringConfig.ax, drawStringConfig.ay)

	// for test
	// dc.SetRGB(GetGgStyleRGB(0, 255, 0))
	// dc.DrawCircle(float64(size.X/2), float64(size.Y/2), 5)
	// dc.Fill()
	return dc.Image(), nil
}

func GetGgStyleRGBA(r int, g int, b int, a int) (float64, float64, float64, float64) {
	return float64(r) / 255, float64(g) / 255, float64(b) / 255, float64(a) / 255
}

func GetGgStyleRGB(r int, g int, b int) (float64, float64, float64) {
	return float64(r) / 255, float64(g) / 255, float64(b) / 255
}

func keysValueToMap(keys []string, value string, mapRef map[string]mapset.Set) {
	for _, k := range keys {
		if _, exist := mapRef[k]; !exist {
			mapRef[k] = mapset.NewSet()
		}
		mapRef[k].Add(value)
	}
}

func mapContainsKey(maps map[string]mapset.Set, key string) bool {
	if _, exist := maps[key]; exist {
		return true
	}
	return false
}

func getSetRandom(set mapset.Set) interface{} {
	size := set.Cardinality()
	if size == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix())
	randomIdx := rand.Intn(size) // [0,size)的随机值，返回值为int
	var res interface{} = set.ToSlice()[randomIdx]
	return res
}
