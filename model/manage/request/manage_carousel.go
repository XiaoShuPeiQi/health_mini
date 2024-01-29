package request

import (
	"main.go/model/common/request"
	"main.go/model/manage"
)

type MallCarouselSearch struct {
	manage.MallCarousel
	request.PageInfo
}

type MallCarouselAddParam struct {
	CarouselUrl  string `json:"carouselUrl"`  //轮播图图片地址
	RedirectUrl  string `json:"redirectUrl"`  //点击后要跳转的地址
	CarouselRank string `json:"carouselRank"` //排序值，表示轮播图显示顺序
}

type MallCarouselUpdateParam struct {
	CarouselId   int    `json:"carouselId"`
	CarouselUrl  string `json:"carouselUrl"`
	RedirectUrl  string `json:"redirectUrl"`
	CarouselRank string `json:"carouselRank" `
}
