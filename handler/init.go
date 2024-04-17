package handler

import (
	"net/http"
	"net/url"

	co "github.com/McaxDev/Back/config"
	ai "github.com/sashabaranov/go-openai"
)

// 初始化函数
func HandlerInit() {
	conf := ai.DefaultConfig(co.Config.GptToken)
	proxyurl, err := url.Parse(co.Config.ProxyAddr)
	if err != nil {
		co.SysLog("ERROR", "GPT网络代理启动失败")
	} else {
		transport := &http.Transport{Proxy: http.ProxyURL(proxyurl)}
		conf.HTTPClient = &http.Client{Transport: transport}
	}
	cli = ai.NewClientWithConfig(conf)
}
