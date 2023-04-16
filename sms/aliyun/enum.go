package aliyun

import "context"

type NewClientConfig struct {
	Ctx             context.Context
	AccessKeyID     string
	AccessKeySecret string
}

type SendMessagesInput struct {
	Ctx           context.Context
	PhoneNumbers  string
	SignName      string
	TemplateName  string
	TemplateParam string
}
