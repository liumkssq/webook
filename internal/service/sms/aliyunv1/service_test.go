package aliyunv1

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/liumkssq/webook/internal/service/sms"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	type args struct {
		c        *dysmsapi.Client
		signName string
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.c, tt.args.signName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Send(t *testing.T) {
	type fields struct {
		client   *dysmsapi.Client
		signName string
	}
	type args struct {
		ctx     context.Context
		tplId   string
		args    []string
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				client:   tt.fields.client,
				signName: tt.fields.signName,
			}
			if err := s.Send(tt.args.ctx, tt.args.tplId, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SendOrigin(t *testing.T) {
	type fields struct {
		client   *dysmsapi.Client
		signName string
	}
	type args struct {
		ctx     context.Context
		tplId   string
		args    map[string]string
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				client:   tt.fields.client,
				signName: tt.fields.signName,
			}
			if err := s.SendOrigin(tt.args.ctx, tt.args.tplId, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("SendOrigin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_SendV1(t *testing.T) {
	type fields struct {
		client   *dysmsapi.Client
		signName string
	}
	type args struct {
		ctx     context.Context
		tplId   string
		args    []sms.NamedArg
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				client:   tt.fields.client,
				signName: tt.fields.signName,
			}
			if err := s.SendV1(tt.args.ctx, tt.args.tplId, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("SendV1() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
