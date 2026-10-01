package transports

import (
	"audita-api-gateway/pkg/endpoint"
	model "audita-api-gateway/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	apperrors "audita-api-gateway/pkg/errors"

	"audita-api-gateway/pkg/ratelimit"

	"github.com/audita-bids/private-kit/decode"
	"github.com/audita-bids/private-kit/keys"
	"github.com/audita-bids/private-kit/look"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/automations"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/client"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/limits"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/pncp"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/whitelabel"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla/v2"
	"golang.org/x/time/rate"
)

var analysisLimiter = ratelimit.New(rate.Every(1*time.Minute), 5, ratelimit.ByToken)

var copilotLimiter = ratelimit.New(rate.Every(1*time.Minute), 15, ratelimit.ByToken)

func NewHTTPServer(endpoint endpoint.EndpointSetup, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	apmgorilla.Instrument(r)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), keys.CompanyContext, req.Header.Get("X-Company-Id"))))
		})
	})

	r.Methods(http.MethodGet).
		Path("/available-licenses").
		Handler(httptransport.NewServer(
			endpoint.GetAvailableLicenses,
			getAvailableLicensesDecodeHTTPRequest,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
				return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
					{
						Key:    keys.AuthTokenContext,
						Header: "Authorization",
						Value:  r.Header.Get("Authorization"),
					},
				})
			}),
		))

	r.Methods(http.MethodGet).
		Path("/license/{cnpj}/{year}/{sequence}").
		Handler(httptransport.NewServer(
			endpoint.GetLicense,
			getFindLicenseDecodeHTTPRequest,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
				return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
					{
						Key:    keys.AuthTokenContext,
						Header: "Authorization",
						Value:  r.Header.Get("Authorization"),
					},
				})
			}),
		))

	r.Methods(http.MethodGet).
		Path("/favorite-bids").
		Handler(httptransport.NewServer(
			endpoint.GetListFavoriteBid,
			decodeGetListFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/favorite-bid").
		Handler(httptransport.NewServer(
			endpoint.PostFavoriteBid,
			decodePostFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/{bid_id}/analysis").
		Handler(analysisLimiter.Middleware(httptransport.NewServer(
			endpoint.PostAnalysis,
			decodeAnalysisHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		)))

	r.Methods(http.MethodGet).
		Path("/analysis/runner/{runner_id}").
		Handler(httptransport.NewServer(
			endpoint.GetAsyncRunner,
			decodeGetAsyncRunnerHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/copilot").
		Handler(copilotLimiter.Middleware(httptransport.NewServer(
			endpoint.PostCopilot,
			decodeCopilotHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		)))

	r.Methods(http.MethodPost).
		Path("/holding-bid").
		Handler(httptransport.NewServer(
			endpoint.PostHoldingBid,
			decodeHoldingHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/holding-bids").
		Handler(httptransport.NewServer(
			endpoint.GetListHoldingBid,
			decodeGetListHoldingBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichListHoldingBidFilter)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.PostWhitelabel,
			decodeWhitelabelHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPatch).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.UpdateWhitelabel,
			decodeWhitelabelUpdateHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/whitelabel").
		Handler(httptransport.NewServer(
			endpoint.GetWhitelabel,
			decodeWhitelabelHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichListHoldingBidFilter)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids").
		Handler(httptransport.NewServer(
			endpoint.GetBids,
			decodeBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichListBidFilter)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids/{id}").
		Handler(httptransport.NewServer(
			endpoint.GetBid,
			decodeGetBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids/{id}/handles").
		Handler(httptransport.NewServer(
			endpoint.GetBidHandles,
			decodeGetBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/bids/{bid_id}/proposals").
		Handler(httptransport.NewServer(
			endpoint.PostBidProposal,
			decodeProposalHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/payment").
		Handler(httptransport.NewServer(
			endpoint.PostPayment,
			decodePaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/payment/{id}/cancel").
		Handler(httptransport.NewServer(
			endpoint.CancelPayment,
			decodeCancelPaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/billings/plans").
		Handler(httptransport.NewServer(
			endpoint.PostPlan,
			decodePlanRequestHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/billings/plans").
		Handler(httptransport.NewServer(
			endpoint.GetListPlans,
			decodePlanRequestHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/billings/payments").
		Handler(httptransport.NewServer(
			endpoint.GetListUserPayments,
			decodeListPaymentHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/clients/business").
		Handler(httptransport.NewServer(
			endpoint.GetListAllUsers,
			decodeListClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichListAllUsersFilter)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/clients/business").
		Handler(httptransport.NewServer(
			endpoint.PostCreateBusinessClient,
			decodeBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodDelete).
		Path("/clients/business/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteBusinessClient,
			decodeDeleteBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodPatch).
		Path("/clients/business/{id}").
		Handler(httptransport.NewServer(
			endpoint.PatchBusinessClient,
			decodeBusinessClientHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, request *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, request, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  request.Header.Get("Authorization"),
						},
					})
				}),
		))

	r.Methods(http.MethodGet).
		Path("/clients/scopes").
		Handler(httptransport.NewServer(
			endpoint.GetListAllScopes,
			decodeScopeHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/automations").
		Handler(httptransport.NewServer(
			endpoint.ListAutomations,
			decodeListAutomationsHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichAutomationFilter)
				},
			),
		))

	r.Methods(http.MethodDelete).
		Path("/bids/favorites/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteFavoriteBid,
			decodeDeleteFavoriteBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/certificates").
		Handler(httptransport.NewServer(
			endpoint.PostCertificate,
			decodePostCertificateHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/certificates").
		Handler(httptransport.NewServer(
			endpoint.ListCertificates,
			decodeListCertificatesHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/certificates/export").
		Handler(httptransport.NewServer(
			endpoint.PostExportCertificates,
			decodePostExportCertificatesHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/bids/export").
		Handler(httptransport.NewServer(
			endpoint.PostExportSearch,
			decodeBidHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, enrichListBidFilter)
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/bids/export/runner/{runner_id}").
		Handler(httptransport.NewServer(
			endpoint.GetSearchRunner,
			decodeGetAsyncRunnerHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/reports/export").
		Handler(httptransport.NewServer(
			endpoint.PostExportReports,
			decodePostExportReportsHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/reports/export/runner/{runner_id}").
		Handler(httptransport.NewServer(
			endpoint.GetReportsRunner,
			decodeGetAsyncRunnerHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/certificates/export/runner/{runner_id}").
		Handler(httptransport.NewServer(
			endpoint.GetCertificatesRunner,
			decodeGetAsyncRunnerHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/certificates/{id}").
		Handler(httptransport.NewServer(
			endpoint.GetCertificate,
			decodeCertificateByIdHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPatch).
		Path("/certificates/{id}").
		Handler(httptransport.NewServer(
			endpoint.PatchCertificate,
			decodePatchCertificateHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodDelete).
		Path("/certificates/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteCertificate,
			decodeCertificateByIdHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/coupons").
		Handler(httptransport.NewServer(
			endpoint.PostCoupon,
			decodePostCouponHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/coupons/search/{code}").
		Handler(httptransport.NewServer(
			endpoint.GetAndValidateCoupon,
			decodeCouponHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodPost).
		Path("/limits").
		Handler(httptransport.NewServer(
			endpoint.PostLimit,
			decodePostLimitHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/limits").
		Handler(httptransport.NewServer(
			endpoint.GetListLimits,
			decodeListLimitsHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.HTTPLook(ctx, r, nil)
				},
			),
		))

	r.Methods(http.MethodPatch).
		Path("/limits/{id}").
		Handler(httptransport.NewServer(
			endpoint.PatchLimit,
			decodePatchLimitHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodDelete).
		Path("/limits/{id}").
		Handler(httptransport.NewServer(
			endpoint.DeleteLimit,
			decodeLimitByIdHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	r.Methods(http.MethodGet).
		Path("/consumes").
		Handler(httptransport.NewServer(
			endpoint.GetListConsumes,
			decodeListConsumesHTTP,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(encodeError),
			httptransport.ServerBefore(
				func(ctx context.Context, r *http.Request) context.Context {
					return decode.InjectHeaderToContext(ctx, r, []decode.HeaderToContext{
						{
							Key:    keys.AuthTokenContext,
							Header: "Authorization",
							Value:  r.Header.Get("Authorization"),
						},
					})
				},
			),
		))

	return r
}

func getAvailableLicensesDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req pncp.PncpAvailableLicenseRequest

	q := r.URL.Query()

	req.DataInicial = q.Get("start_date")
	req.DataFinal = q.Get("final_date")

	if v := q.Get("city_ibge_code"); v != "" {
		req.CodigoMunicipioIbge = v
	}

	if v := q.Get("modality_code"); v != "" {
		req.CodigoModalidadeContratacao = v
	}

	if v := q.Get("uf"); v != "" {
		req.Uf = v
	}

	if v := q.Get("page"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("page inválido: %w", err)
		}
		req.Pagina = int32(parsed)
	}

	if v := q.Get("rows"); v != "" {
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("rows inválido: %w", err)
		}
		req.TamanhoPagina = int32(parsed)
	}

	return &req, nil
}

func getFindLicenseDecodeHTTPRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	var req pncp.PncpFindLicenseRequest

	req.Cnpj = vars["cnpj"]
	req.Ano = vars["year"]
	seqStr := vars["sequence"]

	seqInt, _ := strconv.Atoi(seqStr)
	req.Sequencial = int32(seqInt)

	if req.Cnpj == "" || req.Ano == "" || req.Sequencial == 0 {
		return nil, errors.New("invalid params")
	}

	return &req, nil
}

func decodeGetListFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest
	return &req, nil
}

func decodeGetListHoldingBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.HoldingRequest

	query := r.URL.Query()
	if value, ok := decode.QueryParam(query, "publication_month"); ok {
		converted, err := strconv.Atoi(value)

		if err != nil {
			return nil, nil
		}

		req.PublicationMonth = int32(converted)
	}

	if value, ok := decode.QueryParam(query, "publication_year"); ok {
		converted, err := strconv.Atoi(value)

		if err != nil {
			return nil, nil
		}

		req.PublicationYear = int32(converted)
	}

	return &req, nil
}

func decodePostFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeDeleteFavoriteBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.FavoriteBidRequest

	vars := mux.Vars(r)
	req.Id = vars["id"]

	return &req, nil
}

func decodeAnalysisHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AnalysisRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	req.BidId = mux.Vars(r)["bid_id"]

	return &req, nil
}

func decodeCopilotHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CopilotRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeProposalHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ProposalRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	req.BidId = mux.Vars(r)["bid_id"]

	return &req, nil
}

func decodeHoldingHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.HoldingRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeWhitelabelHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.WhitelabelRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeListPaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	return &req, nil
}

func decodePaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeCancelPaymentHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PaymentRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeListClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ClientRequest
	return &req, nil
}

func decodeBusinessClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BusinessClientRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodePlanRequestHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PlanRequest
	req.Decode(r)

	return &req, nil
}

func decodeWhitelabelUpdateHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.WhitelabelRequest

	if err := r.ParseMultipartForm(12 << 20); err != nil {
		return nil, err
	}

	// Text fields arrive as multipart form values.
	req.Id = r.FormValue("id")
	req.CompanyName = r.FormValue("company_name")
	req.FontFamily = r.FormValue("font_family")
	req.PrimaryColor = r.FormValue("primary_color")
	req.SecondaryColor = r.FormValue("secondary_color")
	req.AccentColor = r.FormValue("accent_color")
	req.ClearFields = r.Form["clear_fields"]

	if v := r.FormValue("theme"); v != "" {
		if n, convErr := strconv.Atoi(v); convErr == nil {
			req.Theme = whitelabel.Theme(n)
		}
	}

	files := map[string]*io.Reader{
		"logoImage":       &req.LogoImage,
		"mobileLogoImage": &req.MobileLogoImage,
		"backgroundImage": &req.BackgroundImage,
	}
	for key, target := range files {
		if file, _, ferr := r.FormFile(key); ferr == nil {
			*target = file
		}
	}

	return &req, nil
}

func decodeBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BidRequest

	return &req, nil
}

func decodeScopeHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.ScopeRequest

	return &req, nil
}

func decodeDeleteBusinessClientHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BusinessClientRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeGetAsyncRunnerHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.RunnerRequest
	vars := mux.Vars(r)

	req.RunnerId = vars["runner_id"]

	return &req, nil
}

func decodeGetBidHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.BidRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodePostCertificateHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CertificateRequest

	if err := r.ParseMultipartForm(12 << 20); err != nil {
		return nil, err
	}

	req.DocumentCnpj = r.FormValue("document_cnpj")
	req.Title = r.FormValue("title")
	req.Url = r.FormValue("url")
	req.ExpiresAt = r.FormValue("expires_at")

	if v := r.FormValue("type"); v != "" {
		if n, convErr := strconv.Atoi(v); convErr == nil {
			req.Type = n
		}
	}

	if file, header, ferr := r.FormFile("file"); ferr == nil {
		req.File = file
		req.FileName = header.Filename
		req.FileSize = header.Size
	}

	return &req, nil
}

func decodeCertificateByIdHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CertificateRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodePatchCertificateHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CertificateRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	req.Id = vars["id"]

	return &req, nil
}

func decodeListCertificatesHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CertificateRequest

	query := r.URL.Query()

	if value, ok := decode.QueryParam(query, "document_cnpj"); ok {
		req.DocumentCnpj = value
	}

	if value, ok := decode.QueryParam(query, "type"); ok {
		if n, convErr := strconv.Atoi(value); convErr == nil {
			req.Type = n
		}
	}

	return &req, nil
}

func decodePostExportReportsHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.HoldingRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodePostExportCertificatesHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CertificateRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodePostCouponHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.PostCouponRequest
	err = req.Decode(r)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func decodeCouponHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.CouponRequest
	vars := mux.Vars(r)

	req.Code = vars["code"]

	return &req, nil
}

func decodePostLimitHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.LimitRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	return &req, nil
}

func decodePatchLimitHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.LimitRequest

	if err := req.Decode(r); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	req.Id = vars["id"]

	return &req, nil
}

func decodeLimitByIdHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.LimitRequest
	vars := mux.Vars(r)

	req.Id = vars["id"]

	return &req, nil
}

func decodeListLimitsHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.LimitRequest

	query := r.URL.Query()

	if value, ok := decode.QueryParam(query, "organization_id"); ok {
		req.OrganizationId = value
	}

	if value, ok := decode.QueryParam(query, "company_id"); ok {
		req.CompanyId = value
	}

	if value, ok := decode.QueryParam(query, "payer_role"); ok {
		if n, convErr := strconv.Atoi(value); convErr == nil {
			req.PayerRole = client.PayerRole(n)
		}
	}

	return &req, nil
}

func decodeListConsumesHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.LimitRequest

	query := r.URL.Query()

	if value, ok := decode.QueryParam(query, "type"); ok {
		if n, convErr := strconv.Atoi(value); convErr == nil {
			req.Type = limits.ConsumeType(n)
		}
	}

	return &req, nil
}

func decodeListAutomationsHTTP(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req model.AutomationsRequest

	query := r.URL.Query()

	if value, ok := decode.QueryParam(query, "type"); ok {
		if converted, err := strconv.Atoi(value); err == nil {
			req.Type = automations.AutomationType(converted)
		}
	}

	if value, ok := decode.QueryParam(query, "status"); ok {
		if converted, err := strconv.Atoi(value); err == nil {
			req.Status = automations.AutomationStatus(converted)
		}
	}

	return &req, nil
}

func encodeHttpResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	writeError(w, apperrors.ParseError(err))
}

func writeError(w http.ResponseWriter, httpErr *apperrors.HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	json.NewEncoder(w).Encode(httpErr)
}

// enriches
func enrichListHoldingBidFilter(ctx context.Context, r *http.Request, filter *look.Params) {
	q := r.URL.Query()

	matchFields := []string{
		"process_id",
		"origin",
	}

	for _, field := range matchFields {
		if value, ok := decode.QueryParam(q, field); ok {
			filter.Conditions = append(filter.Conditions, look.Condition{
				Field:    field,
				Operator: "eq",
				Value:    value,
			})
		}
	}

	dateFields := []string{
		"proposal_opening_date",
		"proposal_closing_date",
	}

	for _, field := range dateFields {
		if startValue, ok := decode.QueryParam(q, field+"_start"); ok {
			filter.Conditions = append(filter.Conditions, look.Condition{
				Field:    field,
				Operator: "gte",
				Value:    startValue,
			})
		}

		if endValue, ok := decode.QueryParam(q, field+"_end"); ok {
			filter.Conditions = append(filter.Conditions, look.Condition{
				Field:    field,
				Operator: "lte",
				Value:    endValue,
			})
		}
	}
}

func enrichListAllUsersFilter(ctx context.Context, r *http.Request, filter *look.Params) {
	q := r.URL.Query()

	matchFields := []string{
		"status",
		"scope_id",
	}

	for _, field := range matchFields {
		if value, ok := decode.QueryParam(q, field); ok {
			filter.Conditions = append(filter.Conditions, look.Condition{
				Field:    field,
				Operator: "eq",
				Value:    value,
			})
		}
	}

	if _, ok := decode.QueryParam(q, "status"); !ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "status",
			Operator: "ne",
			Value:    strconv.Itoa(int(client.ClientStatus_SoftDeleted)),
		})
	}
}

func enrichListBidFilter(ctx context.Context, r *http.Request, filter *look.Params) {
	q := r.URL.Query()

	if value, ok := decode.QueryParam(q, "modalities"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "modality",
			Operator: "in",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "sphere"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "sphere",
			Operator: "eq",
			Value:    sphereLetterToCode(value),
		})
	}

	if value, ok := decode.QueryParam(q, "min_value"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "estimated_value",
			Operator: "gte",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "max_value"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "estimated_value",
			Operator: "lte",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "date_start"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "proposal_opening_date",
			Operator: "gte",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "date_end"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "proposal_opening_date",
			Operator: "lte",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "status"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "status",
			Operator: "in",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "uf"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "uf",
			Operator: "in",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "city_ibge"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "city_ibge",
			Operator: "in",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "proposal_closing_date_start"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "proposal_closing_date",
			Operator: "gte",
			Value:    value,
		})
	}
}

func enrichAutomationFilter(ctx context.Context, r *http.Request, filter *look.Params) {
	q := r.URL.Query()

	if value, ok := decode.QueryParam(q, "month"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "month",
			Operator: "eq",
			Value:    value,
		})
	}

	if value, ok := decode.QueryParam(q, "year"); ok {
		filter.Conditions = append(filter.Conditions, look.Condition{
			Field:    "year",
			Operator: "eq",
			Value:    value,
		})
	}
}

func sphereLetterToCode(letter string) string {
	switch letter {
	case "M":
		return "0"
	case "E":
		return "1"
	case "F":
		return "2"
	default:
		return "3"
	}
}
