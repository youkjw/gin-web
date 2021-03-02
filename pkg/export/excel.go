package export

import (
	"fmt"
	"gin-web/pkg/file"
	"gin-web/pkg/setting"
)

func GetExcelFullUrl(name string) string{
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	path := fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, GetExcelPath())
	if err := file.IsNotExistMkDir(path); err != nil {
		fmt.Printf("get excel path err:%v", err)
	}

	return path
}