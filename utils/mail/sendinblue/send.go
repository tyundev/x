package sendinblue

import (
	"context"
	"fmt"
	"sync"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
)

type RequestParam struct {
	TemplateName string     `json:"template_name"`
	TemplateID   int64      `json:"template_id"`
	Values       []MapValue `json:"values"`
	EmailTo      []string   `json:"emails"`
	CC           []string   `json:"cc"`
	BCC          []string   `json:"bcc"`
}

type MapValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ClientEmail struct {
	Client *sendinblue.APIClient
	m      sync.Mutex
}

var client *ClientEmail

func GetClient(apiKey string) (*ClientEmail, error) {
	if client == nil {
		err := initClient(apiKey)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func initClient(apiKey string) error {
	var ctx context.Context

	cfg := sendinblue.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	clientEm := sendinblue.NewAPIClient(cfg)
	_, _, err := clientEm.AccountApi.GetAccount(ctx)
	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		return err
	}

	client = &ClientEmail{
		Client: clientEm,
	}
	return nil
}

func GetList(ctx *gin.Context, limit, offset int64) (interface{}, error) {
	client.m.Lock()
	defer client.m.Unlock()
	var tps, _, err = client.Client.TransactionalEmailsApi.GetSmtpTemplates(ctx, &sendinblue.TransactionalEmailsApiGetSmtpTemplatesOpts{
		Limit:  optional.NewInt64(limit),
		Offset: optional.NewInt64(offset),
	})
	if err != nil {
		return nil, err
	}
	return tps, err
}

func SendEmail(ctx *gin.Context, pr *RequestParam) (interface{}, error) {
	client.m.Lock()
	defer client.m.Unlock()

	// var tps, _, err = client.Client.TransactionalEmailsApi.GetSmtpTemplate(ctx, pr.TemplateID)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("======== DEN DAY", err)
	if len(pr.EmailTo) == 0 {
		return nil, fmt.Errorf("mail to not exist")
	}
	var emailTos = []sendinblue.SendSmtpEmailTo{}
	for _, val := range pr.EmailTo {
		em := sendinblue.SendSmtpEmailTo{
			Email: val,
		}
		emailTos = append(emailTos, em)
	}
	var emailBccs = []sendinblue.SendSmtpEmailBcc{}
	for _, val := range pr.BCC {
		em := sendinblue.SendSmtpEmailBcc{
			Email: val,
		}
		emailBccs = append(emailBccs, em)
	}
	var emailCcs = []sendinblue.SendSmtpEmailCc{}
	for _, val := range pr.CC {
		em := sendinblue.SendSmtpEmailCc{
			Email: val,
		}
		emailCcs = append(emailCcs, em)
	}
	if len(emailBccs) == 0 {
		emailBccs = nil
	}
	if len(emailCcs) == 0 {
		emailCcs = nil
	}
	var dataPrs interface{}
	var prs = make(map[string]string)
	for _, val := range pr.Values {
		prs[val.Key] = val.Value
	}
	if len(prs) > 0 {
		dataPrs = prs
	}

	var res, a, err1 = client.Client.TransactionalEmailsApi.SendTransacEmail(ctx, sendinblue.SendSmtpEmail{
		To:         emailTos,
		Bcc:        emailBccs,
		Cc:         emailCcs,
		TemplateId: pr.TemplateID,
		Params:     &dataPrs,
	})
	fmt.Println("ERROR: ", err1, a)
	return res, err1

}
