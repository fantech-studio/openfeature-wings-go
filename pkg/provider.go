package openfeatureproviderwings

import (
	"context"
	"net/http"

	"github.com/fantech-studio/openfeature-wings-go/internal"
	of "github.com/open-feature/go-sdk/openfeature"
)

type Provider struct {
	client *internal.Client
}

// NewProvider returns a new instance of the Provider for Wings implementing the Open Feature
func NewProvider() of.FeatureProvider {
	return &Provider{}
}

func (*Provider) Metadata() of.Metadata {
	return of.Metadata{
		Name: "wings",
	}
}

func (p *Provider) BooleanEvaluation(
	ctx context.Context, flag string, defaultValue bool, evalCtx of.FlattenedContext,
) of.BoolResolutionDetail {
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	res, err := p.client.Do(ctx, "/bool:evaluate", http.MethodPost, reqBody)
	if err != nil {
		return of.BoolResolutionDetail{}
	}
	return of.BoolResolutionDetail{
		Value: res.Bool.Value,
		ProviderResolutionDetail: of.ProviderResolutionDetail{
			Variant: res.Variant,
		},
	}
}

func (p *Provider) StringEvaluation(
	ctx context.Context, flag string, defaultValue string, evalCtx of.FlattenedContext,
) of.StringResolutionDetail {
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	res, err := p.client.Do(ctx, "/string:evaluate", http.MethodPost, reqBody)
	if err != nil {
		return of.StringResolutionDetail{}
	}
	return of.StringResolutionDetail{
		Value: res.String.Value,
		ProviderResolutionDetail: of.ProviderResolutionDetail{
			Variant: res.Variant,
		},
	}
}

func (*Provider) FloatEvaluation(
	ctx context.Context, flag string, defaultValue float64, evalCtx of.FlattenedContext,
) of.FloatResolutionDetail {
	// TODO: Implement
	return of.FloatResolutionDetail{}
}

func (p *Provider) IntEvaluation(
	ctx context.Context, flag string, defaultValue int64, evalCtx of.FlattenedContext,
) of.IntResolutionDetail {
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	res, err := p.client.Do(ctx, "/int:evaluate", http.MethodPost, reqBody)
	if err != nil {
		return of.IntResolutionDetail{}
	}
	return of.IntResolutionDetail{
		Value: res.Int.Value,
		ProviderResolutionDetail: of.ProviderResolutionDetail{
			Variant: res.Variant,
		},
	}
}

func (p *Provider) ObjectEvaluation(
	ctx context.Context, flag string, defaultValue interface{}, evalCtx of.FlattenedContext,
) of.InterfaceResolutionDetail {
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	res, err := p.client.Do(ctx, "/object:evaluate", http.MethodPost, reqBody)
	if err != nil {
		return of.InterfaceResolutionDetail{}
	}
	return of.InterfaceResolutionDetail{
		Value: res.Object.Value,
		ProviderResolutionDetail: of.ProviderResolutionDetail{
			Variant: res.Variant,
		},
	}
}

func (*Provider) Hooks() []of.Hook {
	return []of.Hook{}
}
