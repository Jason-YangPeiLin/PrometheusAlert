package controllers

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//华为云短信子程序
func PostHWmessage(text string,mobile string)(string)  {
	open:=beego.AppConfig.String("open-hwdx")
	if open=="0" {
		return "华为云短信接口未配置未开启状态,请先配置open-hwdx为1"
	}
	hwappkey:=beego.AppConfig.String("APP_Key")
	hwappsecret:=beego.AppConfig.String("APP_Secret")
	hwappurl:=beego.AppConfig.String("APP_Url")
	hwtplid:=beego.AppConfig.String("Templateid")
	hwsign:=beego.AppConfig.String("Signature")
	sender:=beego.AppConfig.String("Sender")
	//mobile格式:"15395105573,16619875573"
    //生成header
	now:=time.Now().Format("2006-01-02T15:04:05Z")
	nonce := "7226249334"
	digest := getSha256Code(nonce+now+hwappsecret)
	//digestBase64:=base64.URLEncoding.EncodeToString([]byte(digest))
	digestBase64:=base64.StdEncoding.EncodeToString([]byte(digest))
	xheader:=`"UsernameToken Username="`+hwappkey+`",PasswordDigest="`+digestBase64+`",Nonce="`+nonce+`",Created="`+now+`"`
	log.SetPrefix("[DEBUG 2]")
	log.Println(xheader)
	tr :=&http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", hwappurl + "/sms/batchSendSms/v1", strings.NewReader(url.Values{"from":{sender},"to":{mobile},"templateId":{hwtplid},"templateParas":{"["+text+"]"},"signature":{hwsign},"statusCallback":{""}}.Encode()))
	req.Header.Set("Authorization", `WSSE realm="SDP",profile="UsernameToken",type="Appkey"`)
	req.Header.Set("X-WSSE", xheader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	result,err:=ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	return string(result)
}