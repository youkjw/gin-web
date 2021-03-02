package service

import (
	"encoding/json"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/cache"
	"gin-web/pkg/export"
	"gin-web/pkg/logging"
	"gin-web/pkg/redis"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

func (t *Tag) GetTagsKey() string {
 	keys := []string{
		cache.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}

func (t *Tag) GetAll () ([]models.Tag, error){
	var (
		tags, cacheTags []models.Tag
	)

	key := t.GetTagsKey()
	if redis.Exists(key) {
		data, err := redis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetAllTags(t.getMaps())
	if err != nil {
		return nil, err
	}

	redis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{}{
	maps := make(map[string]interface{})

	if t.Name != "" {
		maps["name"] = t.Name
	}

	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	sheetName := "Sheet1"
	f := excelize.NewFile()

	data := map[string]string{"A1":"ID", "B1":"名称", "C1":"状态", "D1":"创建时间", "E1":"修改时间"}

	for k, tag := range tags {
		data["A" + strconv.Itoa(k + 2)] = strconv.Itoa(tag.ID)
		data["B" + strconv.Itoa(k + 2)] = tag.Name

		var state string
		switch tag.State {
		case 1:
			state = "启用"
		default:
			state = "禁用"
		}
		data["C" + strconv.Itoa(k + 2)] = state

		data["D" + strconv.Itoa(k + 2)] = time.Unix(int64(tag.CreateAt), 0).Format("2006-01-02 15:04:05")
		data["E" + strconv.Itoa(k + 2)] = time.Unix(int64(tag.UpdateAt), 0).Format("2006-01-02 15:04:05")

	}

	for k, v := range data {
		f.SetCellValue(sheetName, k, v)
	}

	etime := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + etime + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	err = f.SaveAs(fullPath)
	if err != nil {
		return "", err
	}

	return filename, nil
}

//func (t *Tag) Export() (string, error) {
//	tags, err := t.GetAll()
//	if err != nil {
//		return "", err
//	}
//
//	file := xlsx.NewFile()
//	sheet, err := file.AddSheet("标签信息")
//	if err != nil {
//		return "", err
//	}
//
//	titles := []string{"ID", "名称", "状态","创建时间", "修改时间"}
//	row := sheet.AddRow()
//
//	var cell *xlsx.Cell
//	for _, title := range titles {
//		cell = row.AddCell()
//		cell.Value = title
//	}
//
//	for _, v := range tags {
//		var state string
//		switch v.State {
//		case 1:
//			state = "启用"
//		default:
//			state = "禁用"
//		}
//
//		values := []string{
//			strconv.Itoa(v.ID),
//			v.Name,
//			state,
//			time.Unix(int64(v.CreateAt), 0).Format("2006-01-02 15:04:05"),
//			time.Unix(int64(v.UpdateAt), 0).Format("2006-01-02 15:04:05"),
//		}
//
//		row = sheet.AddRow()
//		for _, value := range values {
//			cell = row.AddCell()
//			cell.Value = value
//		}
//	}
//
//	etime := strconv.Itoa(int(time.Now().Unix()))
//	filename := "tags-" + etime + ".xlsx"
//
//	fullpath := export.GetExcelFullPath() + filename
//	err = file.Save(fullpath)
//	if err != nil {
//		return "", err
//	}
//
//	return filename, nil
//}

func (t *Tag) ExportIn(r io.Reader) error {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return fmt.Errorf("open file err:%v", err)
	}

	rows, err := f.GetRows("tag1")
	if err != nil {
		return fmt.Errorf("get excel rows err:%v", err)
	}

	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for k, cell := range row {
				if k == 2 {
					switch cell {
					case "启用":
						cell = strconv.Itoa(1)
					default:
						cell = strconv.Itoa(0)
					}
				}
				data = append(data, cell)
			}

			if len(data) > 0 {
				models.AddTag(data[1], 1)
			}
		}
	}

	return nil
}