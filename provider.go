package openfeaturewings

import (
	"context"
	"errors"
	"net/http"

	of "github.com/open-feature/go-sdk/openfeature"

	"github.com/fantech-studio/openfeature-wings-go/internal"
)

var _ of.FeatureProvider = (*Provider)(nil)

type Provider struct {
	client *internal.Client
}

// NewProvider returns a new instance of the Provider for Wings implementing the Open Feature
func NewProvider(host string) of.FeatureProvider {
	return &Provider{
		client: &internal.Client{
			Cli:  new(http.Client),
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
	const path = "/bool:evaluate"

	if flag == "" {
		return of.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: of.ProviderResolutionDetail{
				FlagMetadata:    of.FlagMetadata(evalCtx),
				ResolutionError: of.NewGeneralResolutionError("flag must be non-empty"),
			},
		}
	}

	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}

	res, err := p.client.Do(ctx, path, http.MethodPost, reqBody)
	if err != nil {
		var e of.ResolutionError
		if errors.As(err, &e) {
			resolutionDetail.ResolutionError = e
			return of.BoolResolutionDetail{ProviderResolutionDetail: resolutionDetail}
		}
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.BoolResolutionDetail{ProviderResolutionDetail: resolutionDetail}
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
	const path = "/string:evaluate"

	if flag == "" {
		return of.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: of.ProviderResolutionDetail{
				FlagMetadata:    of.FlagMetadata(evalCtx),
				ResolutionError: of.NewGeneralResolutionError("flag must be non-empty"),
			},
		}
	}

	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}

	res, err := p.client.Do(ctx, path, http.MethodPost, reqBody)
	if err != nil {
		var e of.ResolutionError
		if errors.As(err, &e) {
			resolutionDetail.ResolutionError = e
			return of.StringResolutionDetail{ProviderResolutionDetail: resolutionDetail}
		}
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.StringResolutionDetail{ProviderResolutionDetail: resolutionDetail}
	}

	resolutionDetail.Variant = res.Variant
	return of.StringResolutionDetail{
		Value:                    res.String.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (p *Provider) FloatEvaluation(
	ctx context.Context, flag string, defaultValue float64, evalCtx of.FlattenedContext,
) of.FloatResolutionDetail {
	const path = "/float:evaluate"

	if flag == "" {
		return of.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: of.ProviderResolutionDetail{
				FlagMetadata:    of.FlagMetadata(evalCtx),
				ResolutionError: of.NewGeneralResolutionError("flag must be non-empty"),
			},
		}
	}

	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}

	res, err := p.client.Do(ctx, path, http.MethodPost, reqBody)
	if err != nil {
		var e of.ResolutionError
		if errors.As(err, &e) {
			resolutionDetail.ResolutionError = e
			return of.FloatResolutionDetail{ProviderResolutionDetail: resolutionDetail}
		}
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.FloatResolutionDetail{ProviderResolutionDetail: resolutionDetail}
	}

	resolutionDetail.Variant = res.Variant
	return of.FloatResolutionDetail{
		Value:                    res.Float.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (p *Provider) IntEvaluation(
	ctx context.Context, flag string, defaultValue int64, evalCtx of.FlattenedContext,
) of.IntResolutionDetail {
	const path = "/int:evaluate"

	if flag == "" {
		return of.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: of.ProviderResolutionDetail{
				FlagMetadata:    of.FlagMetadata(evalCtx),
				ResolutionError: of.NewGeneralResolutionError("flag must be non-empty"),
			},
		}
	}

	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}

	res, err := p.client.Do(ctx, path, http.MethodPost, reqBody)
	if err != nil {
		var e of.ResolutionError
		if errors.As(err, &e) {
			resolutionDetail.ResolutionError = e
			return of.IntResolutionDetail{ProviderResolutionDetail: resolutionDetail}
		}
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.IntResolutionDetail{ProviderResolutionDetail: resolutionDetail}
	}

	resolutionDetail.Variant = res.Variant
	return of.IntResolutionDetail{
		Value:                    res.Int.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (p *Provider) ObjectEvaluation(
	ctx context.Context, flag string, defaultValue any, evalCtx of.FlattenedContext,
) of.InterfaceResolutionDetail {
	const path = "/object:evaluate"

	if flag == "" {
		return of.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: of.ProviderResolutionDetail{
				FlagMetadata:    of.FlagMetadata(evalCtx),
				ResolutionError: of.NewGeneralResolutionError("flag must be non-empty"),
			},
		}
	}

	reqBody := &internal.EvalRequest{
		ID:   flag,
		Meta: evalCtx,
	}
	resolutionDetail := of.ProviderResolutionDetail{
		FlagMetadata: of.FlagMetadata(evalCtx),
	}

	res, err := p.client.Do(ctx, path, http.MethodPost, reqBody)
	if err != nil {
		var e of.ResolutionError
		if errors.As(err, &e) {
			resolutionDetail.ResolutionError = e
			return of.InterfaceResolutionDetail{ProviderResolutionDetail: resolutionDetail}
		}
		resolutionDetail.ResolutionError = of.NewGeneralResolutionError(err.Error())
		return of.InterfaceResolutionDetail{ProviderResolutionDetail: resolutionDetail}
	}

	resolutionDetail.Variant = res.Variant
	return of.InterfaceResolutionDetail{
		Value:                    res.Object.Value,
		ProviderResolutionDetail: resolutionDetail,
	}
}

func (*Provider) Hooks() []of.Hook {
	return make([]of.Hook, 0, 0)
}
