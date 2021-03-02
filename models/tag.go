package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model
	Name string `json:"name"`
	State int `json:"state"`
}

func GetAllTags (maps interface{}) (tags []Tag, err error){
	err = db["go"].Where(maps).Find(&tags).Error
	return
}

func GetTags (pageNum int, pageSize int, maps interface{}) (tags []Tag, err error){
	err = db["go"].Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	return
}

func GetTagTotal (maps interface{}) (count int) {
	db["go"].Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistsTagByName (name string) bool {
	var tag Tag
	db["go"].Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func ExistsTagByID(id int) bool {
	var tag Tag
	db["go"].Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func AddTag(name string, state int) bool {
	db["go"].Create(&Tag{
		Model: Model{},
		Name:  name,
		State: state,
	})

	return true
}

func EditTag (id int, data interface{}) bool {
	db["go"].Model(&Tag{}).Where("id = ?", id).Update(data)

	return true
}

func DeteleTag (id int) bool {
	db["go"].Where("id = ?", id).Delete(&Tag{})

	return true
}

func (tag *Tag) BeforeCreate (scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	scope.SetColumn("UpdateAt", time.Now().Unix())

	return nil
}


func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())

	return nil
}