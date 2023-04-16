package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	aliutil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/jau1jz/cornus/commons/log"
)

type Client struct {
	aliyunSmsClient *dysmsapi20170525.Client
}

func (c *Client) Send(input SendMessagesInput) error {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String("18683031655"),
		SignName:      tea.String("天才助理"),
		TemplateCode:  tea.String("SMS_262305037"),
		TemplateParam: tea.String(`{"code":"123456"}`),
	}
	runtime := &aliutil.RuntimeOptions{}
	// 复制代码运行请自行打印 API 的返回值
	_, err := c.aliyunSmsClient.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		log.Slog.ErrorF(input.Ctx, "send aliyun sms error %s", err)
		return err
	}
	return nil
}

func NewAliyunSMSClient(config NewClientConfig) (*Client, error) {
	openApiConfig := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &config.AccessKeyID,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &config.AccessKeySecret,
	}
	// 访问的域名
	openApiConfig.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result := &dysmsapi20170525.Client{}
	_result, err := dysmsapi20170525.NewClient(openApiConfig)
	if err != nil {
		log.Slog.ErrorF(config.Ctx, "new aliyun sms client error %s", err)
	}
	return &Client{aliyunSmsClient: _result}, nil
}
