package qrcode

import (
	"gin-web/pkg/file"
	"gin-web/pkg/setting"
	"gin-web/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
)

type QrCode struct {
	URL string
	Width int
	Height int
	Ext string
	Level qr.ErrorCorrectionLevel
	Mode qr.Encoding
}

const (
	EXT_JPG = ".jpg"
	WIDTH = 200
	HEIGH = 200
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Ext:    EXT_JPG,
		Level:  level,
		Mode:   mode,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()

	if file.CheckNotExist(src) == true {
		return false
	}

	return true
}

func (q *QrCode) Encode(path string) (string, error) {
	filename := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + filename

	if file.CheckNotExist(src) == true {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", err
		}

		f, err := file.Open(src)
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}