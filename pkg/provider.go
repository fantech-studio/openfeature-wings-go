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
func NewProvider(host string) of.FeatureProvider {
	return &Provider{
		client: &internal.Client{
			Cli:  &http.Client{},
			Host: host,
		},
	}
}

func (*Provider) Metadata() of.Metadata {
	return of.Metadata{
		Name: "wings",
	}
}

func (p *Provider) BooleanEvaluation(
	ctx context.Context, flag string, defaultValue bool, evalCtx of.FlattenedContext,
) of.BoolResolutionDetail {
	if flag == "" {
		return of.BoolResolutionDetail{}
	}
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}
	res, statusCode, err := p.client.Do(ctx, "/bool:evaluate", http.MethodPost, reqBody)
	if err != nil {
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.BoolResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	if resolutionErrFunc := p.resolveStatusCode(statusCode); resolutionErrFunc != nil {
		resolutionDetail.ResolutionError = resolutionErrFunc(err.Error())
		return of.BoolResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	resolutionDetail.Variant = res.Variant

	return of.BoolResolutionDetail{
		Value:                    res.Bool.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (p *Provider) StringEvaluation(
	ctx context.Context, flag string, defaultValue string, evalCtx of.FlattenedContext,
) of.StringResolutionDetail {
	if flag == "" {
		return of.StringResolutionDetail{}
	}
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}
	res, statusCode, err := p.client.Do(ctx, "/string:evaluate", http.MethodPost, reqBody)
	if err != nil {
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.StringResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	if resolutionErrFunc := p.resolveStatusCode(statusCode); resolutionErrFunc != nil {
		resolutionDetail.ResolutionError = resolutionErrFunc(err.Error())
		return of.StringResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	resolutionDetail.Variant = res.Variant

	return of.StringResolutionDetail{
		Value:                    res.String.Value,
		ProviderResolutionDetail: resolutionDetail,
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
	if flag == "" {
		return of.IntResolutionDetail{}
	}
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}
	res, statusCode, err := p.client.Do(ctx, "/int:evaluate", http.MethodPost, reqBody)
	if err != nil {
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.IntResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	if resolutionErrFunc := p.resolveStatusCode(statusCode); resolutionErrFunc != nil {
		resolutionDetail.ResolutionError = resolutionErrFunc(err.Error())
		return of.IntResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	resolutionDetail.Variant = res.Variant

	return of.IntResolutionDetail{
		Value:                    res.Int.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (p *Provider) ObjectEvaluation(
	ctx context.Context, flag string, defaultValue interface{}, evalCtx of.FlattenedContext,
) of.InterfaceResolutionDetail {
	if flag == "" {
		return of.InterfaceResolutionDetail{}
	}
	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}
	res, statusCode, err := p.client.Do(ctx, "/object:evaluate", http.MethodPost, reqBody)
	if err != nil {
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.InterfaceResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	if resolutionErrFunc := p.resolveStatusCode(statusCode); resolutionErrFunc != nil {
		resolutionDetail.ResolutionError = resolutionErrFunc(err.Error())
		return of.InterfaceResolutionDetail{
			ProviderResolutionDetail: resolutionDetail,
		}
	}
	resolutionDetail.Variant = res.Variant

	return of.InterfaceResolutionDetail{
		Value:                    res.Int.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (*Provider) Hooks() []of.Hook {
	return []of.Hook{}
}

func (*Provider) resolveStatusCode(statusCode int) func(string) of.ResolutionError {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return of.NewInvalidContextResolutionError
	case http.StatusNotFound:
		return of.NewFlagNotFoundResolutionError
	default:
		return of.NewGeneralResolutionError
	}
}
