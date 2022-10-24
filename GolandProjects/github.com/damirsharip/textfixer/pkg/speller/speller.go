package speller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"textfixer/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	apiUrl   = "https://speller.yandex.net"
	resource = "/services/spellservice.json/checkTexts"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CheckTexts(ctx *gin.Context, texts []string) (models.Response, error) {
	res, err := c.send(ctx, texts)

	return res, err
}

//
//func (r *Client) send(ctx *gin.Context, postData []string) (models.Response, error) {
//	//logger := logger.LoggerFromGinContext(ctx)
//	var yr models.Response
//
//	//resp, err := http.PostForm(serviceURL, url.Values{
//	//	"text":   postData.Texts,
//	//	"lang":   {postData.Lang},
//	//	"format": {postData.Format},
//	//})
//	//if err != nil {
//	//	return
//	//}
//
//	data := url.Values{}
//	for _, v := range postData {
//		data.Set("text", v)
//	}
//
//	u, err := url.ParseRequestURI(apiUrl)
//	if err != nil {
//		//c.logger.Error(err)
//		//ctx.JSON(http.StatusInternalServerError, gin.H{
//		//	"status": 500,
//		//	"error":  "internal server error",
//		//})
//		return yr, err
//	}
//	u.Path = resource
//	u.RawQuery = data.Encode()
//
//	client := &http.Client{}
//	req, err := http.NewRequest(http.MethodGet, u.String(), nil) // URL-encoded payload
//	if err != nil {
//		//h.logger.Error(err)
//		//ctx.JSON(http.StatusInternalServerError, gin.H{
//		//	"status": 500,
//		//	"error":  "internal server error while creating request",
//		//})
//		return yr, err
//	}
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	resp, err := client.Do(req)
//
//	if err != nil {
//		//h.logger.Error(err)
//		//ctx.JSON(http.StatusInternalServerError, gin.H{
//		//	"status": 500,
//		//	"error":  "internal server error while doing request",
//		//})
//		return yr, err
//	}
//
//	//body, err := ioutil.ReadAll(resp.Body)
//	//if err != nil {
//	//	return
//	//}
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		//h.logger.Error(err)
//		//ctx.JSON(http.StatusInternalServerError, gin.H{
//		//	"status": 500,
//		//	"error":  "internal server error while doing request",
//		//})
//		return yr, err
//	}
//
//	err = resp.Body.Close()
//	if err != nil {
//		return yr, err
//	}
//
//	if err = json.Unmarshal(body, &yr); err != nil {
//		return yr, err
//	}
//	//if err != nil {
//	//	h.logger.Error(err)
//	//	ctx.JSON(http.StatusInternalServerError, gin.H{
//	//		"status": 500,
//	//		"error":  "internal server error while doing request",
//	//	})
//	//	return
//	//}
//
//	if resp.StatusCode != http.StatusOK {
//		//logger.Info(string(body))
//		return yr, errors.New(fmt.Sprint("Response status: ", resp.Status))
//	}
//
//	return yr, err
//}

func (r *Client) send(ctx *gin.Context, postData []string) (yr models.Response, err error) {
	//logger := log.LoggerFromGinContext(ctx)
	//
	resp, err := http.PostForm("https://speller.yandex.net/services/spellservice.json/checkTexts", url.Values{
		"text": postData,
	})
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &yr); err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		//logger.Info(string(body))
		return yr, errors.New(fmt.Sprint("Response status: ", resp.Status))
	}

	return
}