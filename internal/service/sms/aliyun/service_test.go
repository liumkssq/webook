package aliyun

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ecodeclub/ekit"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

/**
   @author：biguanqun
   @since： 2023/8/20
   @desc：
**/

//func TestSender(t *testing.T) {
//
//	keyId := ""
//	keySecret := ""
//
//	config := &openapi.Config{
//		AccessKeyId:     ekit.ToPtr[string](keyId),
//		AccessKeySecret: ekit.ToPtr[string](keySecret),
//	}
//	client, err := sms.NewClient(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//	service := NewService(client)
//
//	testCases := []struct {
//		signName string
//		tplCode  string
//		phone    string
//		wantErr  error
//	}{
//		{
//			signName: "webook",
//			tplCode:  "SMS_462745194",
//			phone:    "",
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.signName, func(t *testing.T) {
//			er := service.SendSms(context.Background(), tc.signName, tc.tplCode, tc.phone)
//			assert.Equal(t, tc.wantErr, er)
//		})
//	}
//}

func TestService_SendSms(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	keyId := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	keySecret := os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")

	config := &openapi.Config{
		AccessKeyId:     ekit.ToPtr[string](keyId),
		AccessKeySecret: ekit.ToPtr[string](keySecret),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	client, err := sms.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		panic("client is nil")
	}
	service := NewService(client)

	tests := []struct {
		signName string
		tplCode  string
		phone    []string
		wantErr  error
	}{
		{
			signName: "个人学习开发自用",
			tplCode:  "SMS_475815554",
			phone:    []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.signName, func(t *testing.T) {
			er := service.SendSms(context.Background(), tt.signName, tt.tplCode, tt.phone)
			assert.Equal(t, tt.wantErr, er)
		})
	}
}
