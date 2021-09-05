package fcm

import (
	"github.com/reiwav/x/mlog"

	fcm "github.com/NaySoftware/go-fcm"
)

var logFCM = mlog.NewTagLog("FCM")

const (
	RESPONSE_FAIL = "fail"
)

type FcmClient struct {
	*fcm.FcmClient
}

type FmcMessage struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type FcmMessageData struct {
	Data interface{} `json:"data,omitempty"`
	FmcMessage
}

func NewFCM(serverKey string) *FcmClient {
	return &FcmClient{
		FcmClient: fcm.NewFcmClient(serverKey),
	}
}

func (f *FcmClient) SendToMany(ids []string, data FmcMessage, platform string) (error, string) {
	var noti = fcm.NotificationPayload{
		Title: data.Title,
		Body:  data.Body,
	}
	f.NewFcmRegIdsMsg(ids, data)
	if platform != "android" {
		f.SetNotificationPayload(&noti)
	}
	status, err := f.Send()
	if err != nil {
		logFCM.Debugln(0, err)
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOne(id string, data FmcMessage, platform string) (error, string) {
	return f.SendToMany([]string{id}, data, platform)
}

func (f *FcmClient) SendToManyData(ids []string, data FcmMessageData, platform string) (error, string) {
	f.NewFcmRegIdsMsg(ids, data.Data)

	if platform != "android" {
		var noti = fcm.NotificationPayload{
			Title: data.Title,
			Body:  data.Body,
		}
		f.SetNotificationPayload(&noti)
	}
	status, err := f.Send()
	if err != nil {
		logFCM.Debugln(0, err)
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOneData(id string, data FcmMessageData, platform string) (error, string) {
	return f.SendToManyData([]string{id}, data, platform)
}
