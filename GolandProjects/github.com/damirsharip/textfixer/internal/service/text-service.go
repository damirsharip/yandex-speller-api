package service

import (
	"strings"

	"textfixer/internal/models"

	"github.com/gin-gonic/gin"
)

//newTexts := make([][]string, len(input.Texts))
//for i := range input.Texts {
//	strs := strings.Split(input.Texts[i], " ")
//	n := make([]string, len(input.Texts[i]))
//	n = strs
//	newTexts[i] = n
//}

//for i := range yr[0] {
//	//yr[0][i].Word = yr[0][i].S[0]
//	newTexts[i][yr[0][i].Pos-1] = yr[0][i].S[0]
//}

//correctorTexts := make([]string, len(input.Texts))
//for i, text := range input.Texts {
//for _, misspell := range yr[i] {
//if len(misspell.S) > 0 {
//text = strings.Replace(text, misspell.Word, misspell.S[0], -1)
//}
//}
//correctorTexts[i] = text
//}

func (s *Service) CheckSpell(ctx *gin.Context, req models.TextsInput) (models.CorrectorResponse, error) {
	res := models.CorrectorResponse{}
	spellerResp, err := s.Client.CheckTexts(ctx, req.Texts)
	if err != nil {
		return res, err
	}

	correctorTexts := make([]string, len(req.Texts))
	for i, text := range req.Texts {
		for _, misspell := range spellerResp[i] {
			if len(misspell.Suggestions) > 0 {
				text = strings.Replace(text, misspell.Word, misspell.Suggestions[0], -1)
			}
		}
		correctorTexts[i] = text
	}
	res.Texts = correctorTexts

	return res, nil
}
