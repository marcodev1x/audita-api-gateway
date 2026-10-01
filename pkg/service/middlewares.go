package service

import (
	apperrors "audita-api-gateway/pkg/errors"
	model "audita-api-gateway/request"
	"context"
	"slices"
	"time"

	"github.com/Oudwins/zog"
	"github.com/audita-bids/private-kit/decode"
	"github.com/audita-bids/private-kit/middlewares"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/agents"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/automations"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/bids"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/billings"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/certificates"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/client"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/coupons"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/limits"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/pncp"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/whitelabel"
	"github.com/go-kit/log"
	"google.golang.org/protobuf/types/known/structpb"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *loggingMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetAvailableLicenses", "status", "completed")
	}()

	mw.logger.Log("method", "GetAvailableLicenses", "status", "started")
	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *loggingMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetLicense", "status", "completed")
	}()

	mw.logger.Log("method", "GetLicense", "status", "started")
	return mw.next.GetLicense(ctx, request)
}

func (mw *loggingMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "PostFavoriteBid", "status", "completed")
	}()

	mw.logger.Log("method", "PostFavoriteBid", "status", "started")
	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *loggingMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListFavoriteBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetListFavoriteBid", "status", "started")
	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *loggingMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostAnalysis", "status", "completed")
	}()

	mw.logger.Log("method", "PostAnalysis", "status", "started")
	return mw.next.PostAnalysis(ctx, request)
}

func (mw *loggingMiddleware) GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*agents.AgentsAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetAsyncRunner", "status", "completed")
	}()

	mw.logger.Log("method", "GetAsyncRunner", "status", "started")
	return mw.next.GetAsyncRunner(ctx, request)
}

func (mw *loggingMiddleware) PostCopilot(ctx context.Context, request *model.CopilotRequest) (*agents.AgentsComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostCopilot", "status", "completed")
	}()

	mw.logger.Log("method", "PostCopilot", "status", "started")
	return mw.next.PostCopilot(ctx, request)
}

func (mw *loggingMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostHoldingBid", "status", "completed")
	}()

	mw.logger.Log("method", "PostHoldingBid", "status", "started")
	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *loggingMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListHoldingBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetListHoldingBid", "status", "started")
	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *loggingMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "PostWhitelabel", "status", "started")
	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "GetWhitelabel", "status", "started")
	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	defer func() {
		mw.logger.Log("method", "UpdateWhitelabel", "status", "completed")
	}()

	mw.logger.Log("method", "UpdateWhitelabel", "status", "started")
	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *loggingMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetBids", "status", "completed")
	}()

	mw.logger.Log("method", "GetBids", "status", "started")
	return mw.next.GetBids(ctx, request)
}

func (mw *loggingMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetBid", "status", "completed")
	}()

	mw.logger.Log("method", "GetBid", "status", "started")
	return mw.next.GetBid(ctx, request)
}

func (mw *loggingMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetBidHandles", "status", "completed")
	}()

	mw.logger.Log("method", "GetBidHandles", "status", "started")
	return mw.next.GetBidHandles(ctx, request)
}

func (mw *loggingMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostBidProposal", "status", "completed")
	}()

	mw.logger.Log("method", "PostBidProposal", "status", "started")
	return mw.next.PostBidProposal(ctx, request)
}

func (mw *loggingMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	defer func() {
		mw.logger.Log("method", "UpdateBidProposal", "status", "completed")
	}()

	mw.logger.Log("method", "UpdateBidProposal", "status", "started")
	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *loggingMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	defer func() {
		mw.logger.Log("method", "PostPayment", "status", "completed")
	}()

	mw.logger.Log("method", "PostPayment", "status", "started")
	return mw.next.PostPayment(ctx, request)
}

func (mw *loggingMiddleware) CancelPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostCancelPaymentResponse, error) {
	defer func() {
		mw.logger.Log("method", "CancelPayment", "status", "completed")
	}()

	mw.logger.Log("method", "CancelPayment", "status", "started")
	return mw.next.CancelPayment(ctx, request)
}

func (mw *loggingMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostPlan", "status", "completed")
	}()

	mw.logger.Log("method", "PostPlan", "status", "started")
	return mw.next.PostPlan(ctx, request)
}

func (mw *loggingMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListPlans", "status", "completed")
	}()

	mw.logger.Log("method", "GetListPlans", "status", "started")
	return mw.next.GetListPlans(ctx, request)
}

func (mw *loggingMiddleware) GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListUserPayments", "status", "completed")
	}()

	mw.logger.Log("method", "GetListUserPayments", "status", "started")
	return mw.next.GetListUserPayments(ctx, request)
}

func (mw *loggingMiddleware) GetListAllUsers(ctx context.Context, request *model.ClientRequest) (*client.ListAllClientsResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListAllUsers", "status", "completed")
	}()

	mw.logger.Log("method", "GetListAllUsers", "status", "started")
	return mw.next.GetListAllUsers(ctx, request)
}

func (mw *loggingMiddleware) PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostCreateBusinessClient", "status", "completed")
	}()

	mw.logger.Log("method", "PostCreateBusinessClient", "status", "started")
	return mw.next.PostCreateBusinessClient(ctx, request)
}

func (mw *loggingMiddleware) GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListAllScopes", "status", "completed")
	}()

	mw.logger.Log("method", "GetListAllScopes", "status", "started")
	return mw.next.GetListAllScopes(ctx, request)
}

func (mw *loggingMiddleware) DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*structpb.Struct, error) {
	defer func() {
		mw.logger.Log("method", "DeleteBusinessClient", "status", "completed")
	}()

	mw.logger.Log("method", "DeleteBusinessClient", "status", "started")
	return mw.next.DeleteBusinessClient(ctx, request)
}

func (mw *loggingMiddleware) PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	defer func() {
		mw.logger.Log("method", "PatchBusinessClient", "status", "completed")
	}()

	mw.logger.Log("method", "PatchBusinessClient", "status", "started")
	return mw.next.PatchBusinessClient(ctx, request)
}

func (mw *loggingMiddleware) ListAutomations(ctx context.Context, request *model.AutomationsRequest) (*automations.ListAutomationsResponse, error) {
	defer func() {
		mw.logger.Log("method", "ListAutomations", "status", "completed")
	}()

	mw.logger.Log("method", "ListAutomations", "status", "started")
	return mw.next.ListAutomations(ctx, request)
}

func (mw *loggingMiddleware) PostCoupon(ctx context.Context, request *model.PostCouponRequest) (*coupons.CouponComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostCoupon", "status", "completed")
	}()

	mw.logger.Log("method", "PostCoupon", "status", "started", "code", request.Code)
	return mw.next.PostCoupon(ctx, request)
}

func (mw *loggingMiddleware) PostCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostCertificate", "status", "completed")
	}()

	mw.logger.Log("method", "PostCertificate", "status", "started")
	return mw.next.PostCertificate(ctx, request)
}

func (mw *loggingMiddleware) GetCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetCertificate", "status", "completed")
	}()

	mw.logger.Log("method", "GetCertificate", "status", "started")
	return mw.next.GetCertificate(ctx, request)
}

func (mw *loggingMiddleware) PatchCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	defer func() {
		mw.logger.Log("method", "PatchCertificate", "status", "completed")
	}()

	mw.logger.Log("method", "PatchCertificate", "status", "started")
	return mw.next.PatchCertificate(ctx, request)
}

func (mw *loggingMiddleware) DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.DeleteCertificateResponse, error) {
	defer func() {
		mw.logger.Log("method", "DeleteCertificate", "status", "completed")
	}()

	mw.logger.Log("method", "DeleteCertificate", "status", "started")
	return mw.next.DeleteCertificate(ctx, request)
}

func (mw *loggingMiddleware) ListCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.ListCertificatesResponse, error) {
	defer func() {
		mw.logger.Log("method", "ListCertificates", "status", "completed")
	}()

	mw.logger.Log("method", "ListCertificates", "status", "started")
	return mw.next.ListCertificates(ctx, request)
}

func (mw *loggingMiddleware) PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostExportCertificates", "status", "completed")
	}()

	mw.logger.Log("method", "PostExportCertificates", "status", "started")
	return mw.next.PostExportCertificates(ctx, request)
}

func (mw *loggingMiddleware) GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetCertificatesAsyncRunner", "status", "completed")
	}()

	mw.logger.Log("method", "GetCertificatesAsyncRunner", "status", "started")
	return mw.next.GetCertificatesAsyncRunner(ctx, request)
}

func (mw *loggingMiddleware) PostExportReports(ctx context.Context, request *model.HoldingRequest) (*bids.BidsAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostExportReports", "status", "completed")
	}()

	mw.logger.Log("method", "PostExportReports", "status", "started")
	return mw.next.PostExportReports(ctx, request)
}

func (mw *loggingMiddleware) GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetReportsAsyncRunner", "status", "completed")
	}()

	mw.logger.Log("method", "GetReportsAsyncRunner", "status", "started")
	return mw.next.GetReportsAsyncRunner(ctx, request)
}

func (mw *loggingMiddleware) PostExportSearch(ctx context.Context, request *model.BidRequest) (*bids.BidsAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostExportSearch", "status", "completed")
	}()

	mw.logger.Log("method", "PostExportSearch", "status", "started")
	return mw.next.PostExportSearch(ctx, request)
}

func (mw *loggingMiddleware) GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetSearchAsyncRunner", "status", "completed")
	}()

	mw.logger.Log("method", "GetSearchAsyncRunner", "status", "started")
	return mw.next.GetSearchAsyncRunner(ctx, request)
}

func (mw *loggingMiddleware) GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (*coupons.CouponComplete, error) {
	defer func() {
		mw.logger.Log("method", "GetAndValidateCoupon", "status", "completed")
	}()

	mw.logger.Log("method", "GetAndValidateCoupon", "status", "started")
	return mw.next.GetAndValidateCoupon(ctx, request)
}

func (mw *loggingMiddleware) DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.DeleteFavoriteBidResponse, error) {
	defer func() {
		mw.logger.Log("method", "DeleteFavoriteBid", "status", "completed")
	}()

	mw.logger.Log("method", "DeleteFavoriteBid", "status", "started")
	return mw.next.DeleteFavoriteBid(ctx, request)
}

func (mw *loggingMiddleware) PostLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	defer func() {
		mw.logger.Log("method", "PostLimit", "status", "completed")
	}()

	mw.logger.Log("method", "PostLimit", "status", "started")
	return mw.next.PostLimit(ctx, request)
}

func (mw *loggingMiddleware) PatchLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	defer func() {
		mw.logger.Log("method", "PatchLimit", "status", "completed")
	}()

	mw.logger.Log("method", "PatchLimit", "status", "started")
	return mw.next.PatchLimit(ctx, request)
}

func (mw *loggingMiddleware) GetListLimits(ctx context.Context, request *model.LimitRequest) (*limits.GetListLimitsResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListLimits", "status", "completed")
	}()

	mw.logger.Log("method", "GetListLimits", "status", "started")
	return mw.next.GetListLimits(ctx, request)
}

func (mw *loggingMiddleware) DeleteLimit(ctx context.Context, request *model.LimitRequest) (*limits.DeleteLimitResponse, error) {
	defer func() {
		mw.logger.Log("method", "DeleteLimit", "status", "completed")
	}()

	mw.logger.Log("method", "DeleteLimit", "status", "started")
	return mw.next.DeleteLimit(ctx, request)
}

func (mw *loggingMiddleware) GetListConsumes(ctx context.Context, request *model.LimitRequest) (*limits.GetListConsumesResponse, error) {
	defer func() {
		mw.logger.Log("method", "GetListConsumes", "status", "completed")
	}()

	mw.logger.Log("method", "GetListConsumes", "status", "started")
	return mw.next.GetListConsumes(ctx, request)
}

func RecoveryMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &recoveryMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type recoveryMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *recoveryMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (resp *pncp.PncpAvailableLicenseResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAvailableLicenses", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *recoveryMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (resp *pncp.PncpFindLicenseResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetLicense", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetLicense(ctx, request)
}

func (mw *recoveryMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (resp *bids.PostFavoriteBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostFavoriteBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (resp *bids.GetListFavoriteBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListFavoriteBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (resp *agents.AgentsComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostAnalysis", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostAnalysis(ctx, request)
}

func (mw *recoveryMiddleware) GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (resp *agents.AgentsAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAsyncRunner", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetAsyncRunner(ctx, request)
}

func (mw *recoveryMiddleware) PostCopilot(ctx context.Context, request *model.CopilotRequest) (resp *agents.AgentsComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostCopilot", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostCopilot(ctx, request)
}

func (mw *recoveryMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (resp *bids.HoldingBidComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostHoldingBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *recoveryMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (resp *bids.GetListHoldingBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListHoldingBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *recoveryMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (resp *whitelabel.WhitelabelComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "UpdateWhitelabel", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *recoveryMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (resp *bids.GetListBidsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBids", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBids(ctx, request)
}

func (mw *recoveryMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (resp *bids.BidComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBid(ctx, request)
}

func (mw *recoveryMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (resp *bids.GetBidClientHandlesResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetBidHandles", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *recoveryMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (resp *bids.ProposalComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostBidProposal", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *recoveryMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (resp *bids.ProposalComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "UpdateBidProposal", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *recoveryMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (resp *billings.PostPaymentResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostPayment", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostPayment(ctx, request)
}

func (mw *recoveryMiddleware) CancelPayment(ctx context.Context, request *model.PaymentRequest) (resp *billings.PostCancelPaymentResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "CancelPayment", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.CancelPayment(ctx, request)
}

func (mw *recoveryMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (resp *billings.PlanComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostPlan", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostPlan(ctx, request)
}

func (mw *recoveryMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (resp *billings.GetListPlansResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListPlans", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListPlans(ctx, request)
}

func (mw *recoveryMiddleware) GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (resp *billings.GetListUserPaymentsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListUserPayments", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListUserPayments(ctx, request)
}

func (mw *recoveryMiddleware) GetListAllUsers(ctx context.Context, request *model.ClientRequest) (resp *client.ListAllClientsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListAllUsers", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListAllUsers(ctx, request)
}

func (mw *recoveryMiddleware) PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (resp *client.ClientComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostCreateBusinessClient", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostCreateBusinessClient(ctx, request)
}

func (mw *recoveryMiddleware) GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListAllScopes", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetListAllScopes(ctx, request)
}

func (mw *recoveryMiddleware) DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (resp *structpb.Struct, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "DeleteBusinessClient", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.DeleteBusinessClient(ctx, request)
}

func (mw *recoveryMiddleware) PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (resp *client.ClientComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PatchBusinessClient", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PatchBusinessClient(ctx, request)
}

func (mw *recoveryMiddleware) ListAutomations(ctx context.Context, request *model.AutomationsRequest) (resp *automations.ListAutomationsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "ListAutomations", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.ListAutomations(ctx, request)
}

func (mw *recoveryMiddleware) PostCoupon(ctx context.Context, request *model.PostCouponRequest) (result *coupons.CouponComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostCoupon", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostCoupon(ctx, request)
}

func (mw *recoveryMiddleware) PostCertificate(ctx context.Context, request *model.CertificateRequest) (result *certificates.CertificateComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostCertificate", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostCertificate(ctx, request)
}

func (mw *recoveryMiddleware) GetCertificate(ctx context.Context, request *model.CertificateRequest) (result *certificates.CertificateComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetCertificate", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetCertificate(ctx, request)
}

func (mw *recoveryMiddleware) PatchCertificate(ctx context.Context, request *model.CertificateRequest) (result *certificates.CertificateComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PatchCertificate", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PatchCertificate(ctx, request)
}

func (mw *recoveryMiddleware) DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (result *certificates.DeleteCertificateResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "DeleteCertificate", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.DeleteCertificate(ctx, request)
}

func (mw *recoveryMiddleware) ListCertificates(ctx context.Context, request *model.CertificateRequest) (result *certificates.ListCertificatesResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "ListCertificates", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.ListCertificates(ctx, request)
}

func (mw *recoveryMiddleware) PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (result *certificates.CertificatesAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostExportCertificates", "status", "recovered", "error", r)
		}
	}()

	return mw.next.PostExportCertificates(ctx, request)
}

func (mw *recoveryMiddleware) GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (result *certificates.CertificatesAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetCertificatesAsyncRunner", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetCertificatesAsyncRunner(ctx, request)
}

func (mw *recoveryMiddleware) PostExportReports(ctx context.Context, request *model.HoldingRequest) (result *bids.BidsAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostExportReports", "status", "recovered", "error", r)
		}
	}()

	return mw.next.PostExportReports(ctx, request)
}

func (mw *recoveryMiddleware) GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (result *bids.BidsAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetReportsAsyncRunner", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetReportsAsyncRunner(ctx, request)
}

func (mw *recoveryMiddleware) PostExportSearch(ctx context.Context, request *model.BidRequest) (result *bids.BidsAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostExportSearch", "status", "recovered", "error", r)
		}
	}()

	return mw.next.PostExportSearch(ctx, request)
}

func (mw *recoveryMiddleware) GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (result *bids.BidsAsyncRunnerComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetSearchAsyncRunner", "status", "recovered", "error", r)
		}
	}()

	return mw.next.GetSearchAsyncRunner(ctx, request)
}

func (mw *recoveryMiddleware) GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (result *coupons.CouponComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetAndValidateCoupon", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetAndValidateCoupon(ctx, request)
}

func (mw *recoveryMiddleware) DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (result *bids.DeleteFavoriteBidResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "DeleteFavoriteBid", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.DeleteFavoriteBid(ctx, request)
}

func (mw *recoveryMiddleware) PostLimit(ctx context.Context, request *model.LimitRequest) (result *limits.LimitsComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PostLimit", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PostLimit(ctx, request)
}

func (mw *recoveryMiddleware) PatchLimit(ctx context.Context, request *model.LimitRequest) (result *limits.LimitsComplete, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "PatchLimit", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.PatchLimit(ctx, request)
}

func (mw *recoveryMiddleware) GetListLimits(ctx context.Context, request *model.LimitRequest) (result *limits.GetListLimitsResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListLimits", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListLimits(ctx, request)
}

func (mw *recoveryMiddleware) DeleteLimit(ctx context.Context, request *model.LimitRequest) (result *limits.DeleteLimitResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "DeleteLimit", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.DeleteLimit(ctx, request)
}

func (mw *recoveryMiddleware) GetListConsumes(ctx context.Context, request *model.LimitRequest) (result *limits.GetListConsumesResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			mw.logger.Log("method", "GetListConsumes", "status", "recovered", "error", r)
			err = apperrors.Internal("internal server error")
		}
	}()

	return mw.next.GetListConsumes(ctx, request)
}

type validationMiddleware struct {
	next   Service
	logger log.Logger
}

func ValidationMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &validationMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw *validationMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *validationMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	schema := zog.Struct(zog.Shape{
		"cnpj":       zog.String().Required(),
		"ano":        zog.String().Required(),
		"sequencial": zog.Int32().Required(),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetLicense(ctx, request)
}

func (mw *validationMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *validationMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *validationMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostAnalysis(ctx, request)
}

func (mw *validationMiddleware) GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*agents.AgentsAsyncRunnerComplete, error) {
	schema := zog.Struct(zog.Shape{
		"RunnerId": zog.String().Required(zog.Message("Runner ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetAsyncRunner(ctx, request)
}

func (mw *validationMiddleware) PostCopilot(ctx context.Context, request *model.CopilotRequest) (*agents.AgentsComplete, error) {
	schema := zog.Struct(zog.Shape{
		"Message": zog.String().
			Required(zog.Message("Message is required")).
			Max(500, zog.Message("Message is too long")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostCopilot(ctx, request)
}

func (mw *validationMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId": zog.String().Required(zog.Message("Bid ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *validationMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	schema := zog.Struct(zog.Shape{
		"PublicationMonth": zog.Int32(),
		"PublicationYear":  zog.Int32(),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *validationMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *validationMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *validationMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *validationMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	return mw.next.GetBids(ctx, request)
}

func (mw *validationMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().Required(zog.Message("ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetBid(ctx, request)
}

func (mw *validationMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().Required(zog.Message("ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *validationMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	schema := zog.Struct(zog.Shape{
		"BidId":       zog.String().Required(zog.Message("Bid ID is required")),
		"value":       zog.Float64().Optional().GTE(0),
		"observation": zog.String().Optional().Max(500),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *validationMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *validationMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	schema := zog.Struct(zog.Shape{
		"PaymentMethod": zog.CustomFunc(func(pm *billings.PaymentMethod, ctx zog.Ctx) bool {
			switch *pm {
			case billings.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
				billings.PaymentMethod_PAYMENT_METHOD_PIX,
				billings.PaymentMethod_PAYMENT_METHOD_BOLETO:
				return true
			default:
				return false
			}
		}, zog.Message("Invalid payment method")),
		"PlanId": zog.String().Required(zog.Message("Plan ID is required")),
		"Payer": zog.Ptr(zog.Struct(zog.Shape{
			"Email": zog.String().Required(zog.Message("Payer email is required")),
		})).NotNil(zog.Message("Payer is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostPayment(ctx, request)
}

func (mw *validationMiddleware) CancelPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostCancelPaymentResponse, error) {
	schema := zog.Struct(zog.Shape{
		"Id": zog.String().Required(zog.Message("Payment ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.CancelPayment(ctx, request)
}

func (mw *validationMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	return mw.next.GetListPlans(ctx, request)
}

func (mw *validationMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	schema := zog.Struct(zog.Shape{
		"Name":         zog.String().Required(zog.Message("Plan name is required")),
		"Description":  zog.String().Optional().Max(500),
		"PriceInCents": zog.Int64().GTE(0),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PostPlan(ctx, request)
}

func (mw *validationMiddleware) GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error) {
	return mw.next.GetListUserPayments(ctx, request)
}

func (mw *validationMiddleware) GetListAllUsers(ctx context.Context, request *model.ClientRequest) (*client.ListAllClientsResponse, error) {
	return mw.next.GetListAllUsers(ctx, request)
}

func (mw *validationMiddleware) PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	schema := zog.Struct(zog.Shape{
		"Name": zog.String().
			Min(2, zog.Message("Name must be at least 5 characters long")).Required(),
		"Email": zog.String().
			Email(zog.Message("Email must be a valid email address")).Required(),
	})

	if err := schema.Validate(&request.CreateBusinessClientRequest); err != nil {
		return nil, decode.ErrorFields(err)
	}

	if !request.GenerateRandomicPassword && len(request.Password) < 8 {
		return nil, apperrors.BadRequest("validation_error", "Password must be at least 8 characters long")
	}

	return mw.next.PostCreateBusinessClient(ctx, request)
}

func (mw *validationMiddleware) GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	return mw.next.GetListAllScopes(ctx, request)
}

func (mw *validationMiddleware) DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*structpb.Struct, error) {
	schema := zog.Struct(zog.Shape{
		"Id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.DeleteBusinessClient(ctx, request)
}

func (mw *validationMiddleware) PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	schema := zog.Struct(zog.Shape{
		"Id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.PatchBusinessClient(ctx, request)
}

func (mw *validationMiddleware) ListAutomations(ctx context.Context, request *model.AutomationsRequest) (*automations.ListAutomationsResponse, error) {
	return mw.next.ListAutomations(ctx, request)
}

func (mw *validationMiddleware) PostCoupon(ctx context.Context, request *model.PostCouponRequest) (*coupons.CouponComplete, error) {
	// zog resolve a chave capitalizando a primeira letra, entao ela tem que ser
	// o nome do campo em Go e nao o do JSON: "start_at" vira "Start_at", que
	// nao existe na struct, e a validacao entra em panic.
	schema := zog.Struct(zog.Shape{
		"code": zog.String().
			Required(zog.Message("Code is required")),
		"discountPercentage": zog.Int().
			Required(zog.Message("Discount percentage is required")).
			GT(0, zog.Message("Discount percentage must be greater than 0")).
			LTE(100, zog.Message("Discount percentage must be at most 100")),
		"startAt": zog.String().
			Required(zog.Message("Start at is required")),
		"endAt": zog.String().
			Required(zog.Message("End at is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	// The coupons service parses these with `time.Parse(time.RFC3339, ...)` and
	// a bad format comes back as a gRPC Unknown, which the gateway can only
	// render as a 500. Parsing here turns "2027-01-01" into a 400 that says so.
	for _, v := range []string{request.StartAt, request.EndAt} {
		if _, err := time.Parse(time.RFC3339, v); err != nil {
			return nil, apperrors.BadRequest("invalid_date", "start_at and end_at must be RFC 3339 with an offset, e.g. 2026-09-04T00:00:00-03:00")
		}
	}

	return mw.next.PostCoupon(ctx, request)
}

func (mw *validationMiddleware) PostCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	schema := zog.Struct(zog.Shape{
		"title": zog.String().
			Required(zog.Message("Title is required")),
		"expiresAt": zog.String().
			Required(zog.Message("Expires at is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	if _, err := time.Parse(time.RFC3339, request.ExpiresAt); err != nil {
		return nil, apperrors.BadRequest("invalid_date", "expires_at must be RFC 3339 with an offset, e.g. 2027-01-01T00:00:00-03:00")
	}

	if request.File == nil && request.Url == "" {
		return nil, apperrors.BadRequest("file_required", "send a file or an url")
	}

	return mw.next.PostCertificate(ctx, request)
}

func (mw *validationMiddleware) GetCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetCertificate(ctx, request)
}

func (mw *validationMiddleware) PatchCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	if len(request.UpdateMask) == 0 {
		return nil, apperrors.BadRequest("update_mask_required", "update_mask must name at least one field")
	}

	if slices.Contains(request.UpdateMask, "expires_at") {
		if _, err := time.Parse(time.RFC3339, request.ExpiresAt); err != nil {
			return nil, apperrors.BadRequest("invalid_date", "expires_at must be RFC 3339 with an offset, e.g. 2027-01-01T00:00:00-03:00")
		}
	}

	return mw.next.PatchCertificate(ctx, request)
}

func (mw *validationMiddleware) DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.DeleteCertificateResponse, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.DeleteCertificate(ctx, request)
}

func (mw *validationMiddleware) ListCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.ListCertificatesResponse, error) {
	return mw.next.ListCertificates(ctx, request)
}

func (mw *validationMiddleware) PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	return mw.next.PostExportCertificates(ctx, request)
}

func (mw *validationMiddleware) GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	schema := zog.Struct(zog.Shape{
		"RunnerId": zog.String().Required(zog.Message("Runner ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetCertificatesAsyncRunner(ctx, request)
}

func (mw *validationMiddleware) PostExportReports(ctx context.Context, request *model.HoldingRequest) (*bids.BidsAsyncRunnerComplete, error) {
	if request.PublicationYear <= 0 {
		return nil, apperrors.BadRequest("publication_year_required", "publication_year is required")
	}

	if request.PublicationMonth < 0 || request.PublicationMonth > 12 {
		return nil, apperrors.BadRequest("invalid_publication_month", "publication_month must be between 1 and 12, or 0 for the whole year")
	}

	return mw.next.PostExportReports(ctx, request)
}

func (mw *validationMiddleware) GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	schema := zog.Struct(zog.Shape{
		"RunnerId": zog.String().Required(zog.Message("Runner ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetReportsAsyncRunner(ctx, request)
}

func (mw *validationMiddleware) PostExportSearch(ctx context.Context, request *model.BidRequest) (*bids.BidsAsyncRunnerComplete, error) {
	return mw.next.PostExportSearch(ctx, request)
}

func (mw *validationMiddleware) GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	schema := zog.Struct(zog.Shape{
		"RunnerId": zog.String().Required(zog.Message("Runner ID is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetSearchAsyncRunner(ctx, request)
}

func (mw *validationMiddleware) GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (*coupons.CouponComplete, error) {
	schema := zog.Struct(zog.Shape{
		"code": zog.String().
			Required(zog.Message("Code is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.GetAndValidateCoupon(ctx, request)
}

func (mw *validationMiddleware) DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.DeleteFavoriteBidResponse, error) {
	schema := zog.Struct(zog.Shape{
		"Id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.DeleteFavoriteBid(ctx, request)
}

func (mw *validationMiddleware) PostLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	if _, ok := limits.ConsumeType_name[int32(request.Type)]; !ok || request.Type == limits.ConsumeType_CONSUME_TYPE_UNSPECIFIED {
		return nil, apperrors.BadRequest("type_required", "type must be a valid consume type")
	}

	if _, ok := limits.LimitPeriod_name[int32(request.Period)]; !ok {
		return nil, apperrors.BadRequest("invalid_period", "period must be a valid limit period")
	}

	if _, ok := client.PayerRole_name[int32(request.PayerRole)]; !ok {
		return nil, apperrors.BadRequest("invalid_payer_role", "payer_role must be a valid payer role")
	}

	if !request.Unlimited && request.Quantity < 0 {
		return nil, apperrors.BadRequest("invalid_quantity", "quantity must be zero or greater")
	}

	return mw.next.PostLimit(ctx, request)
}

func (mw *validationMiddleware) PatchLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	if len(request.UpdateMask) == 0 {
		return nil, apperrors.BadRequest("update_mask_required", "update_mask must name at least one field")
	}

	if slices.Contains(request.UpdateMask, "quantity") && request.Quantity < 0 {
		return nil, apperrors.BadRequest("invalid_quantity", "quantity must be zero or greater")
	}

	if _, ok := limits.LimitPeriod_name[int32(request.Period)]; !ok {
		return nil, apperrors.BadRequest("invalid_period", "period must be a valid limit period")
	}

	return mw.next.PatchLimit(ctx, request)
}

func (mw *validationMiddleware) GetListLimits(ctx context.Context, request *model.LimitRequest) (*limits.GetListLimitsResponse, error) {
	return mw.next.GetListLimits(ctx, request)
}

func (mw *validationMiddleware) DeleteLimit(ctx context.Context, request *model.LimitRequest) (*limits.DeleteLimitResponse, error) {
	schema := zog.Struct(zog.Shape{
		"id": zog.String().
			Required(zog.Message("Id is required")),
	})

	if err := schema.Validate(request); err != nil {
		return nil, decode.ErrorFields(err)
	}

	return mw.next.DeleteLimit(ctx, request)
}

func (mw *validationMiddleware) GetListConsumes(ctx context.Context, request *model.LimitRequest) (*limits.GetListConsumesResponse, error) {
	return mw.next.GetListConsumes(ctx, request)
}

func AuthenticationMiddleware(logger log.Logger, clients client.ClientServiceClient) Middleware {
	return func(next Service) Service {
		return &authenticationMiddleware{
			next:    next,
			logger:  logger,
			clients: clients,
		}
	}
}

var payingPayers = []client.PayerRole{
	client.PayerRole_Basic,
	client.PayerRole_Pro,
	client.PayerRole_Enterprise,
	client.PayerRole_FreeTasting,
}

var genericPayers = []client.PayerRole{
	client.PayerRole_Pro,
	client.PayerRole_Enterprise,
	client.PayerRole_FreeTasting,
}

var copilotPayers = []client.PayerRole{
	client.PayerRole_Pro,
	client.PayerRole_Enterprise,
}

type authenticationMiddleware struct {
	next    Service
	logger  log.Logger
	clients client.ClientServiceClient
	cache   *middlewares.AuthCache
}

func (mw *authenticationMiddleware) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"licenses:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetAvailableLicenses(ctx, request)
}

func (mw *authenticationMiddleware) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"licenses:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetLicense(ctx, request)
}

func (mw *authenticationMiddleware) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostFavoriteBid(ctx, request)
}

func (mw *authenticationMiddleware) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListFavoriteBid(ctx, request)
}

func (mw *authenticationMiddleware) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"ai:write"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostAnalysis(ctx, request)
}

func (mw *authenticationMiddleware) GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*agents.AgentsAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"ai:write"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetAsyncRunner(ctx, request)
}

func (mw *authenticationMiddleware) PostCopilot(ctx context.Context, request *model.CopilotRequest) (*agents.AgentsComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"ai:write"},
		PayerRole: copilotPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostCopilot(ctx, request)
}

func (mw *authenticationMiddleware) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"hold_bids:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostHoldingBid(ctx, request)
}

func (mw *authenticationMiddleware) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"hold_bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListHoldingBid(ctx, request)
}

func (mw *authenticationMiddleware) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"whitelabel:write"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"whitelabel:read"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"whitelabel:write"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.UpdateWhitelabel(ctx, request)
}

func (mw *authenticationMiddleware) GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBids(ctx, request)
}

func (mw *authenticationMiddleware) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBid(ctx, request)
}

func (mw *authenticationMiddleware) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetBidHandles(ctx, request)
}

func (mw *authenticationMiddleware) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostBidProposal(ctx, request)
}

func (mw *authenticationMiddleware) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.UpdateBidProposal(ctx, request)
}

func (mw *authenticationMiddleware) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"payments:write"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostPayment(ctx, request)
}

func (mw *authenticationMiddleware) CancelPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostCancelPaymentResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"payments:write"},
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.CancelPayment(ctx, request)
}

func (mw *authenticationMiddleware) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"plans:write"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostPlan(ctx, request)
}

func (mw *authenticationMiddleware) GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"plans:read"},
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_SuperAdmin,
			client.ClientRole_Admin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListPlans(ctx, request)
}

func (mw *authenticationMiddleware) GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"payments:read"},
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListUserPayments(ctx, request)
}

func (mw *authenticationMiddleware) GetListAllUsers(ctx context.Context, request *model.ClientRequest) (*client.ListAllClientsResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"clients:read"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListAllUsers(ctx, request)
}

func (mw *authenticationMiddleware) PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"client:create"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostCreateBusinessClient(ctx, request)
}

func (mw *authenticationMiddleware) GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"scopes:read"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListAllScopes(ctx, request)
}

func (mw *authenticationMiddleware) DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*structpb.Struct, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"client:create"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.DeleteBusinessClient(ctx, request)
}

func (mw *authenticationMiddleware) PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"client:update"},
		PayerRole: genericPayers,
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PatchBusinessClient(ctx, request)
}

func (mw *authenticationMiddleware) ListAutomations(ctx context.Context, request *model.AutomationsRequest) (*automations.ListAutomationsResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.ListAutomations(ctx, request)
}

func (mw *authenticationMiddleware) PostCoupon(ctx context.Context, request *model.PostCouponRequest) (*coupons.CouponComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"coupons:write"},
		Roles: []client.ClientRole{
			client.ClientRole_Organization,
			client.ClientRole_Business,
			client.ClientRole_Admin,
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostCoupon(ctx, request)
}

func (mw *authenticationMiddleware) PostCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostCertificate(ctx, request)
}

func (mw *authenticationMiddleware) GetCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetCertificate(ctx, request)
}

func (mw *authenticationMiddleware) PatchCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:write"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PatchCertificate(ctx, request)
}

func (mw *authenticationMiddleware) DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.DeleteCertificateResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:delete"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.DeleteCertificate(ctx, request)
}

func (mw *authenticationMiddleware) ListCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.ListCertificatesResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.ListCertificates(ctx, request)
}

func (mw *authenticationMiddleware) PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostExportCertificates(ctx, request)
}

func (mw *authenticationMiddleware) GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"certificates:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetCertificatesAsyncRunner(ctx, request)
}

func (mw *authenticationMiddleware) PostExportReports(ctx context.Context, request *model.HoldingRequest) (*bids.BidsAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"hold_bids:read"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostExportReports(ctx, request)
}

func (mw *authenticationMiddleware) GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"hold_bids:read"},
		PayerRole: genericPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetReportsAsyncRunner(ctx, request)
}

func (mw *authenticationMiddleware) PostExportSearch(ctx context.Context, request *model.BidRequest) (*bids.BidsAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostExportSearch(ctx, request)
}

func (mw *authenticationMiddleware) GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:read"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetSearchAsyncRunner(ctx, request)
}

func (mw *authenticationMiddleware) GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (*coupons.CouponComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"coupons:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetAndValidateCoupon(ctx, request)
}

func (mw *authenticationMiddleware) DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.DeleteFavoriteBidResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes:    []string{"bids:delete"},
		PayerRole: payingPayers,
	})

	if err != nil {
		return nil, err
	}

	return mw.next.DeleteFavoriteBid(ctx, request)
}

func (mw *authenticationMiddleware) PostLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"limits:write"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PostLimit(ctx, request)
}

func (mw *authenticationMiddleware) PatchLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"limits:write"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.PatchLimit(ctx, request)
}

func (mw *authenticationMiddleware) GetListLimits(ctx context.Context, request *model.LimitRequest) (*limits.GetListLimitsResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"limits:read"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListLimits(ctx, request)
}

func (mw *authenticationMiddleware) DeleteLimit(ctx context.Context, request *model.LimitRequest) (*limits.DeleteLimitResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"limits:write"},
		Roles: []client.ClientRole{
			client.ClientRole_SuperAdmin,
		},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.DeleteLimit(ctx, request)
}

func (mw *authenticationMiddleware) GetListConsumes(ctx context.Context, request *model.LimitRequest) (*limits.GetListConsumesResponse, error) {
	user, ctx, err := middlewares.ValidateAuthCached(ctx, mw.clients, mw.cache)

	if err != nil {
		return nil, err
	}

	ctx, user, err = middlewares.ActiveCompany(ctx, user)

	if err != nil {
		return nil, err
	}

	err = middlewares.ValidateScopes(user, &middlewares.Scoping{
		Scopes: []string{"client:read"},
	})

	if err != nil {
		return nil, err
	}

	return mw.next.GetListConsumes(ctx, request)
}
