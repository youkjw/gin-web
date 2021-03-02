package v1

import (
	"encoding/json"
	"fmt"
	"gin-web/pkg/app"
	error2 "gin-web/pkg/error"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type TeacherRank struct {
	Code int
	Msg string
	Data struct {
		List [] struct {
			TeacherUid string `json:"teacher_uid"`
			VoteNumber int `json:"vote_number"`
			Rank int `json:"rank"`
			Nickname string `json:"nickname"`
			Avatar string `json:"avatar"`
	} `json:"list"`
	} `json:"data"`
}

func HttpGet(c *gin.Context) {
	appG := app.Gin{c}
	uri := "https://sandbox-m.yiqiwen.cn/api/user_vote/teacher_rank"

	Params := url.Values{}
	Url, err := url.Parse(uri)
	if err != nil {
		appG.Response(http.StatusOK, error2.ERROR, nil)
		return
	}

	Params.Set("page", strconv.Itoa(1))
	Params.Set("size", strconv.Itoa(10))
	Url.RawQuery = Params.Encode()

	UriPath := Url.String()
	resp, err := http.Get(UriPath)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	var p TeacherRank

	err = json.Unmarshal([]byte(body), &p)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(p)

	appG.Response(http.StatusOK, error2.SUCCESS, p)
}