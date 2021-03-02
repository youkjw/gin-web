package service

import (
	"encoding/json"
	"gin-web/models"
	"gin-web/pkg/cache"
	"gin-web/pkg/logging"
	"gin-web/pkg/redis"
	"strconv"
	"strings"
)

type Article struct {
	ID int
	TagID int
	State int

	PageNum int
	PageSize int
}

func (a *Article) GetArticleKey() string {
	return cache.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *Article) GetArticlesKey() string {
	keys := []string {
		cache.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}

	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}

func (a *Article) GetArticle() (*models.Article, error){
	var cacheArticle *models.Article

	cachePack := Article{ID:a.ID}
	cacheKey := cachePack.GetArticleKey()
	if redis.Exists(cacheKey) {
		data, err := redis.Get(cacheKey)
		if err != nil {
			logging.Info(err)
		}

		json.Unmarshal(data, &cacheArticle)
		return cacheArticle, nil
	}

	maps := make(map[string]interface{})
	maps["id"] = a.ID
	article, err := models.GetArticle(maps)
	if err != nil {
		return nil, err
	}

	redis.Set(cacheKey, article, 3600)
	return article, nil
}