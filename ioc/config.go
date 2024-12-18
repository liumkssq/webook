package ioc

import (
	"context"
	"github.com/spf13/viper"
)

// Configer interface
type Configer interface {
	GetString(ctx context.Context, key string) (string, error)
	MustGetString(ctx context.Context, key string) string
	GetStringOrDefault(ctx context.Context, key string, def string) string
}

type ViperConfigAdapter struct {
	v *viper.Viper
}

type myConfiger struct {
}

func (m *myConfiger) GetString(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *myConfiger) MustGetString(ctx context.Context, key string) string {
	str, err := m.GetString(ctx, key)
	if err != nil {
		panic(err)
	}
	return str
}

func (m *myConfiger) GetStringOrDefault(ctx context.Context, key string, def string) string {
	str, err := m.GetString(ctx, key)
	if err != nil {
		return def
	}
	return str
}
