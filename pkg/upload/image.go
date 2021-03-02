package upload

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	gf "gin-web/pkg/file"
	"gin-web/pkg/logging"
	"gin-web/pkg/setting"
	"gin-web/pkg/util"
)

const (
	EXT_JPG = ".jpg"
)

type Image struct {
	Name string
	Src string
	Width int
	Height int
	Rect Rect
	Pt Pt
}

type Rect struct {
	Name string
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type Pt struct {
	X int
	Y int
}

func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	if ext == "" {
		ext = EXT_JPG
	}

	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := gf.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := gf.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = gf.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := gf.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func (i *Image) MergeImage(name string, bg string) (string, error) {
	imagePath := GetImageFullPath()
	filename := GetImageName(name)

	if err := CheckImage(imagePath + filename); err != nil {
		return "", err
	}

	mergedF, err := gf.Open(imagePath + filename)
	if err != nil {
		return "", err
	}
	defer mergedF.Close()

	bgF, err := gf.Open(bg)
	if err != nil {
		return "", err
	}
	defer bgF.Close()

	distF, err := gf.Open(i.Src)
	if err != nil {
		return "", err
	}
	defer distF.Close()

	bgImage, err := jpeg.Decode(bgF)
	if err != nil {
		return "", err
	}
	distImage, err := jpeg.Decode(distF)
	if err != nil {
		return "", err
	}

	jpg := image.NewRGBA(image.Rect(i.Rect.X0, i.Rect.Y0, i.Rect.X1, i.Rect.Y1))
	draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
	draw.Draw(jpg, jpg.Bounds(), distImage, distImage.Bounds().Min.Sub(image.Pt(i.Pt.X, i.Pt.Y)), draw.Over)

	jpeg.Encode(mergedF, jpg, nil)

	return filename, nil
}