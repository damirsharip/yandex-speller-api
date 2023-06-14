package handler

import (
	"fmt"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
)

type ShortLinkRequest struct {
	LongUrl string `json:"long_url"`
}

func (r ShortLinkRequest) Validate() error {
	if r.LongUrl == "" {
		return fmt.Errorf("invalid long url")
	}

	longUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return err
	}

	if longUrl.Scheme == "" {
		return fmt.Errorf("schema should be provided in long url")
	}

	return nil
}

type ShortLinkResponse struct {
	ShortUrl  string `json:"short_url"`
	SecretKey string `json:"secret_key"`
}

type InfoResponse struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
	Clicks   int    `json:"clicks"`
}

func (r *InfoResponse) FromLink(link *db.Link, baseUrl string) {
	if link == nil {
		return
	}

	r.ShortUrl = shortLinkFunc(baseUrl, link.ShortSuffix)
	r.LongUrl = link.Link
	r.Clicks = link.Clicks
}

type VipLinkRequest struct {
	LongUrl string  `json:"long_url"`
	VipKey  *string `json:"vip_key"`
	Ttl     *int    `json:"ttl"`
	TtlUnit *string `json:"ttl_unit"`
}

func (r *VipLinkRequest) Validate() error {
	if r.LongUrl == "" {
		//return fmt.Errorf("invalid long url")
		r.LongUrl = "https://yandex.ru"
	}

	longUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return err
	}

	if longUrl.Scheme == "" {
		return fmt.Errorf("schema should be provided in long url")
	}

	if r.Ttl == nil || *r.Ttl == 0 {
		t := 10
		r.Ttl = &t
	}
	if r.TtlUnit == nil || *r.TtlUnit == "" {
		s := "HOURS"
		r.TtlUnit = &s
	}

	if *r.TtlUnit == "SECONDS" && *r.Ttl > 172800 {
		return fmt.Errorf("link durability should not exceed 2 days")
	}

	if *r.TtlUnit == "MINUTES" && *r.Ttl > 2880 {
		return fmt.Errorf("link durability should not exceed 2 days")
	}

	if *r.TtlUnit == "HOURS" && *r.Ttl > 48 {
		return fmt.Errorf("link durability should not exceed 2 days")
	}

	if *r.TtlUnit == "DAYS" && *r.Ttl > 2 {

		return fmt.Errorf("link durability should not exceed 2 days")
	}

	return nil
}
