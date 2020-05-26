package wechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

type AlertMsg struct {
	Receiver 		string			`receiver`
	Status			string			`status`
	Alerts			[]Alerts		`alerts`
	GroupLabels		GroupLabels		`groupLabels`
	CommonLabels 	CommonLabels	`commonLabels`
	ExternalURL 	string			`externalURL`
	Version			string			`version`
	GroupKey 		string			`groupKey`
}

type GroupLabels struct {
	Alertname 		string 			`alertname`
}

type CommonLabels struct {
	Alertname 		string 			`alertname`
	Instance 		string			`instance`
	Job				string			`job`
	Severity 		string			`severity`
	Team 			string 			`team`
}

type Alerts struct {
	Status 			string			`status`
	Labels 			CommonLabels	`labels`
	Annotations 	Annotations		`annotations`
	StartsAt 		string			`startsAt`
}

type Annotations struct {
	Summary 		string			`summary`
}

type Handler struct {
	CallbackUrl string
}

func (handler *Handler) HandleRequest(responseWriter http.ResponseWriter, request *http.Request) {
	defer func() {
		_, _ =responseWriter.Write([]byte(""))
	}()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Request read body error - ", err)
		return
	}

	alertMsg := new(AlertMsg)
	if json.Unmarshal(body, alertMsg) != nil {
		return
	}
	sendAlertMsg(handler.CallbackUrl, alertMsg)
}

func sendAlertMsg(callbackUrl string, alertMsg *AlertMsg)  {
	var instanceDetail = "服务报警<font color=\"warning\">[ " + alertMsg.CommonLabels.Job + " ] - [ " + alertMsg.CommonLabels.Alertname + " ]</font>，请相关同事注意。\n\n"
	for _, alert := range alertMsg.Alerts {
		instanceDetail += "> 实例:<font color=\"comment\">" + alert.Labels.Instance + "</font>\n"
		instanceDetail += "> 时间:<font color=\"comment\">" + alert.StartsAt + "</font>\n"
		instanceDetail += "> 原因:<font color=\"comment\">" +  alert.Annotations.Summary + "</font>\n\n"
	}

	param := make(map[string]interface{})
	param["msgtype"] = "markdown"
	param["markdown"] = map[string]string{"content" : instanceDetail}
	jsonStr, _ := json.Marshal(param)
	resp, err := http.Post(
		callbackUrl,
		"application/json",
		bytes.NewReader([]byte(jsonStr)))

	if err != nil {
		log.Println("Response error on wechat webhook: ", err)
		return
	}
	sendMsgBytes, _ := ioutil.ReadAll(resp.Body)
	log.Println("Respones on wechat webhook: ", string(sendMsgBytes))
	_ = resp.Close
}