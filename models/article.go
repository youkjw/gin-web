package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model
	TagID int `json:"tag_id" gorm:"index"`
	Tag Tag `json:"tag"`
	Title string `json:"title"`
	Desc string `json:"desc"`
	State int `json:"state"`
}

func GetArticle (maps interface{}) (*Article, error){
	var article Article
	db["go"].Where(maps).First(&article)
	err := db["go"].Model(&article).Related(&article.Tag).Error

	return &article, err
}

func GetArticles (pageNum int, pageSize int, maps interface{}) (articles []Article){
	db["go"].Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func GetArticleTotal (maps interface{}) (count int) {
	db["go"].Model(&Article{}).Where(maps).Count(&count)

	return
}

func ExistsArticleByID (id int) bool {
	var article Article
	db["go"].Select("id").Where("id = ?", id).First(&article)
	if article.ID > 0 {
		return true
	}

	return false
}

func AddArticle (data map[string]interface{}) bool {
	db["go"].Create(&Article{
		Title: data["title"].(string),
		Desc: data["desc"].(string),
		State: data["state"].(int),
		TagID: data["tag_id"].(int),
	})

	return true
}

func EditArticle (id int, data interface{}) bool {
	db["go"].Model(&Article{}).Where("id = ?", id).Update(data)

	return true
}

func DeteleArticle (id int) bool {
	db["go"].Where("id = ?", id).Delete(&Article{})

	return true
}

func (article *Article) BeforeCreate (scope *gorm.Scope) error {
	scope.SetColumn("CreateAt", time.Now().Unix())
	scope.SetColumn("UpdateAt", time.Now().Unix())

	return nil
}


func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdateAt", time.Now().Unix())

	return nil
}