package aliyun

import (
	"context"
	"errors"
	"fmt"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/ecodeclub/ekit"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

/**
   @author：biguanqun
   @since： 2023/8/20
   @desc：
**/

type Service struct {
	//appId    *string
	//signName *string
	client *sms.Client
	//limiter  ratelimit.Limiter
}

//func NewService(appId *string, signName *string, client *sms.Client, limiter ratelimit.Limiter) *Service {
//	return &Service{
//		appId:    appId,
//		signName: signName,
//		client:   client,
//		limiter:  limiter,
//	}
//}

func NewService(client *sms.Client) *Service {
	if client == nil {
		zap.L().Fatal("阿里云短信服务初始化失败", zap.Error(errors.New("阿里云短信服务初始化失败")))
	}
	return &Service{
		client: client,
	}
}

// SendSms 单次
func (s *Service) SendSms(ctx context.Context,
	signName, tplCode string, phone []string) error {
	phoneLen := len(phone)
	// phone1 phone2
	//    0     1
	for i := 0; i < phoneLen; i++ {
		phoneSignle := phone[i]

		// 1. 生成验证码
		code := fmt.Sprintf("%06v",
			rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		// 完全没有做成一个独立的发短信的实现。而是一个强耦合验证码的实现。
		bcode, _ := json.Marshal(map[string]interface{}{
			"code": code,
		})

		// 2. 初始化短信结构体
		smsRequest := &sms.SendSmsRequest{
			SignName:      ekit.ToPtr[string](signName),
			TemplateCode:  ekit.ToPtr[string](tplCode),
			PhoneNumbers:  ekit.ToPtr[string](phoneSignle),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
		}
		//fmt.Println("短信发送请求参数", smsRequest)
		//zap.L().Info("短信发送请求参数", zap.Any("smsRequest", smsRequest))

		// 3. 发送短信
		smsResponse, err := s.client.SendSms(smsRequest)
		if err != nil {
			fmt.Println("发送短信失败", err)
			zap.L().Fatal("阿里云短信服务发送失败", zap.Error(err))
			panic(err)
		}
		if *smsResponse.Body.Code == "OK" {
			fmt.Println(phoneSignle, string(bcode))
			fmt.Printf("发送手机号: %s 的短信成功,验证码为【%s】\n", phoneSignle, code)
		}
		fmt.Println(errors.New(*smsResponse.Body.Message))

	}
	return nil
}
