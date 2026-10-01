package service

import (
	model "audita-api-gateway/request"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/audita-bids/private-kit/connectors"
	"github.com/audita-bids/private-kit/decode"
	"github.com/audita-bids/private-kit/keys"
	"github.com/audita-bids/private-kit/look"
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
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"resty.dev/v3"
)

type Service interface {
	GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error)
	GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error)
	PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error)
	DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.DeleteFavoriteBidResponse, error)
	GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error)
	PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error)
	GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*agents.AgentsAsyncRunnerComplete, error)
	PostCopilot(ctx context.Context, request *model.CopilotRequest) (*agents.AgentsComplete, error)
	PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error)
	GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error)
	PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	GetWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error)
	GetBids(ctx context.Context, request *model.BidRequest) (*bids.GetListBidsResponse, error)
	GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error)
	GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error)
	PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error)
	UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error)
	PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error)
	CancelPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostCancelPaymentResponse, error)
	PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error)
	GetListPlans(ctx context.Context, request *model.PlanRequest) (*billings.GetListPlansResponse, error)
	GetListUserPayments(ctx context.Context, request *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error)
	GetListAllUsers(ctx context.Context, request *model.ClientRequest) (*client.ListAllClientsResponse, error)
	PostCreateBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error)
	GetListAllScopes(ctx context.Context, request *model.ScopeRequest) (*client.ListAllScopesResponse, error)
	DeleteBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*structpb.Struct, error)
	PatchBusinessClient(ctx context.Context, request *model.BusinessClientRequest) (*client.ClientComplete, error)
	ListAutomations(ctx context.Context, request *model.AutomationsRequest) (*automations.ListAutomationsResponse, error)
	GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (*coupons.CouponComplete, error)
	PostCoupon(ctx context.Context, request *model.PostCouponRequest) (*coupons.CouponComplete, error)
	PostCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error)
	GetCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error)
	PatchCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error)
	DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.DeleteCertificateResponse, error)
	ListCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.ListCertificatesResponse, error)
	PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificatesAsyncRunnerComplete, error)
	GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*certificates.CertificatesAsyncRunnerComplete, error)
	PostExportReports(ctx context.Context, request *model.HoldingRequest) (*bids.BidsAsyncRunnerComplete, error)
	GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error)
	PostExportSearch(ctx context.Context, request *model.BidRequest) (*bids.BidsAsyncRunnerComplete, error)
	GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error)
	PostLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error)
	PatchLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error)
	GetListLimits(ctx context.Context, request *model.LimitRequest) (*limits.GetListLimitsResponse, error)
	DeleteLimit(ctx context.Context, request *model.LimitRequest) (*limits.DeleteLimitResponse, error)
	GetListConsumes(ctx context.Context, request *model.LimitRequest) (*limits.GetListConsumesResponse, error)
}

const maxCertificateBytes = 12 << 20

var allowedCertificateTypes = []string{
	"application/pdf",
	"image/png",
	"image/jpeg",
}

var (
	ErrCertificateTooLarge       = errors.New("certificate file is too large")
	ErrCertificateEmptyFile      = errors.New("certificate file is empty")
	ErrCertificateTypeNotAllowed = errors.New("certificate file type is not allowed")
	ErrCertificateUploadFailed   = errors.New("certificate upload failed")
)

type service struct {
	cdnApi *resty.Client
	logger log.Logger

	pncp         pncp.PncpServiceClient
	clients      client.ClientServiceClient
	bids         bids.BidsServiceClient
	agents       agents.AgentsServiceClient
	whitelabel   whitelabel.WhitelabelServiceClient
	billings     billings.BillingServiceClient
	automations  automations.AutomationsServiceClient
	coupons      coupons.CouponsServiceClient
	certificates certificates.CertificatesServiceClient
	limits       limits.LimitsServiceClient
}

func NewService(logger log.Logger) Service {
	clients := client.NewClientServiceClient(connectors.Client())
	cdn := os.Getenv("CDN_HOST")

	var svc Service
	{
		svc = &service{
			cdnApi:       resty.New().SetBaseURL(cdn).SetRetryCount(2), // will do 2 retries
			logger:       logger,
			pncp:         pncp.NewPncpServiceClient(connectors.Pncp()),
			bids:         bids.NewBidsServiceClient(connectors.Bids()),
			clients:      clients,
			agents:       agents.NewAgentsServiceClient(connectors.Agents()),
			whitelabel:   whitelabel.NewWhitelabelServiceClient(connectors.Whitelabel()),
			billings:     billings.NewBillingServiceClient(connectors.Billings()),
			automations:  automations.NewAutomationsServiceClient(connectors.Automations()),
			coupons:      coupons.NewCouponsServiceClient(connectors.Coupons()),
			certificates: certificates.NewCertificatesServiceClient(connectors.Certificates()),
			limits:       limits.NewLimitsServiceClient(connectors.Limits()),
		}
		svc = LoggingMiddleware(logger)(svc)
		svc = RecoveryMiddleware(logger)(svc)
		svc = ValidationMiddleware(logger)(svc)
		svc = AuthenticationMiddleware(logger, clients)(svc)
	}

	return svc
}

func (s *service) GetAvailableLicenses(ctx context.Context, request *pncp.PncpAvailableLicenseRequest) (*pncp.PncpAvailableLicenseResponse, error) {
	return s.pncp.AvailableLicenses(ctx, request)
}

func (s *service) GetLicense(ctx context.Context, request *pncp.PncpFindLicenseRequest) (*pncp.PncpFindLicenseResponse, error) {
	return s.pncp.FindLicense(ctx, request)
}

func (s *service) PostFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.PostFavoriteBidResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	request.UserId = user.Id

	companyId, organizationId := ownership(ctx)

	return s.bids.PostFavoriteBid(ctx, &bids.PostFavoriteBidRequest{
		UserId:         request.UserId,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		BidId:          request.BidId,
	})
}

func (s *service) DeleteFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.DeleteFavoriteBidResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.DeleteFavoriteBid(ctx, &bids.DeleteFavoriteBidRequest{
		Id:             request.Id,
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) GetListFavoriteBid(ctx context.Context, request *model.FavoriteBidRequest) (*bids.GetListFavoriteBidResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	companyId, organizationId := ownership(ctx)

	return s.bids.GetListFavoriteBid(filter.Outgoing(ctx), &bids.GetListFavoriteBidRequest{
		UserId:         request.UserId,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostAnalysis(ctx context.Context, request *model.AnalysisRequest) (*agents.AgentsComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	bid, err := s.bids.GetBid(ctx, &bids.GetBidRequest{Id: request.BidId})
	if err != nil {
		return nil, err
	}

	edital := findEdital(bid.GetFiles())
	if edital == nil {
		return nil, fmt.Errorf("no edital identified among the bid files")
	}

	raw, err := download(ctx, edital.Url)
	if err != nil {
		return nil, err
	}

	companyId, organizationId := ownership(ctx)

	return s.agents.PostAnalysis(ctx, &agents.PostAnalysisRequest{
		UserId:         request.UserId,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		Base64:         base64.StdEncoding.EncodeToString(raw),
		ProcessId:      request.BidId,
	})
}

func (s *service) GetAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*agents.AgentsAsyncRunnerComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.agents.GetAsyncRunner(ctx, &agents.GetAsyncRunnerRequest{
		RunnerId:       request.RunnerId,
		RunnerType:     agents.RunnerType_ANALYSIS_RUNNER,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostCopilot(ctx context.Context, request *model.CopilotRequest) (*agents.AgentsComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	companyId, organizationId := ownership(ctx)

	return s.agents.PostCopilot(ctx, &agents.PostCopilotRequest{
		UserId:         request.UserId,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		Message:        request.Message,
	})
}

func (s *service) PostHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.HoldingBidComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	companyId, organizationId := ownership(ctx)

	return s.bids.PostHoldingBid(ctx, &bids.PostHoldingBidRequest{
		UserId:         request.UserId,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		BidId:          request.BidId,
	})
}

// GetListHoldingBid lists the holdings of the whole company, not only the caller's: a holding belongs to the company.
func (s *service) GetListHoldingBid(ctx context.Context, request *model.HoldingRequest) (*bids.GetListHoldingBidResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.GetListHoldingBid(filter.Outgoing(ctx), &bids.GetListHoldingBidRequest{
		CompanyId:        companyId,
		OrganizationId:   organizationId,
		PublicationMonth: request.PublicationMonth,
		PublicationYear:  request.PublicationYear,
	})
}

func (s *service) PostWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.whitelabel.CreateWhitelabel(ctx, &whitelabel.WhitelabelComplete{
		Id:             request.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) GetWhitelabel(ctx context.Context, _ *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	companyId, _ := ownership(ctx)

	return s.whitelabel.FindWhitelabel(ctx, &whitelabel.FindWhitelabelRequest{
		CompanyId: companyId,
	})
}

func (s *service) UpdateWhitelabel(ctx context.Context, request *model.WhitelabelRequest) (*whitelabel.WhitelabelComplete, error) {
	companyId, organizationId := ownership(ctx)

	if request.MobileLogoImage != nil {
		url, err := s.uploadFile(request.MobileLogoImage, "mobileLogoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.MobileLogoUri = url
		}
	}

	if request.LogoImage != nil {
		url, err := s.uploadFile(request.LogoImage, "logoImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.LogoUri = url
		}
	}

	if request.BackgroundImage != nil {
		url, err := s.uploadFile(request.BackgroundImage, "backgroundImage")
		if err != nil {
			return nil, err
		}

		if url != "" {
			request.BackgroundUri = url
		}
	}

	return s.whitelabel.UpdateWhitelabel(ctx, &whitelabel.WhitelabelComplete{
		Id:              request.Id,
		CompanyId:       companyId,
		OrganizationId:  organizationId,
		ClearFields:     request.ClearFields,
		CompanyName:     request.CompanyName,
		LogoUri:         request.LogoUri,
		MobileLogoUri:   request.MobileLogoUri,
		BackgroundUri:   request.BackgroundUri,
		Theme:           request.Theme,
		FontFamily:      request.FontFamily,
		FontUri:         request.FontUri,
		PrimaryColor:    request.PrimaryColor,
		SecondaryColor:  request.SecondaryColor,
		AccentColor:     request.AccentColor,
		BackgroundColor: request.BackgroundColor,
		SurfaceColor:    request.SurfaceColor,
		TextColor:       request.TextColor,
		BorderColor:     request.BorderColor,
		SuccessColor:    request.SuccessColor,
		ErrorColor:      request.ErrorColor,
		WarningColor:    request.WarningColor,
	})
}

func (s *service) GetBids(ctx context.Context, _ *model.BidRequest) (*bids.GetListBidsResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)

	return s.bids.GetListBids(filter.Outgoing(ctx), &bids.GetListBidsRequest{})
}

func (s *service) GetBid(ctx context.Context, request *model.BidRequest) (*bids.BidComplete, error) {
	return s.bids.GetBid(ctx, &bids.GetBidRequest{
		Id: request.Id,
	})
}

func (s *service) GetBidHandles(ctx context.Context, request *model.BidRequest) (*bids.GetBidClientHandlesResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.GetBidClientHandles(ctx, &bids.GetBidClientHandlesRequest{
		BidId:          request.Id,
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.PostBidProposal(ctx, &bids.PostBidProposalRequest{
		BidId:          request.Id,
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		Value:          request.Value,
		Files:          request.Files,
		Observation:    request.Observation,
	})
}

func (s *service) UpdateBidProposal(ctx context.Context, request *model.ProposalRequest) (*bids.ProposalComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.PutBidProposal(ctx, &bids.PutBidProposalRequest{
		BidId:          request.Id,
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		Value:          request.Value,
		Files:          request.Files,
		Observation:    request.Observation,
	})
}

func (s *service) PostPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostPaymentResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.billings.CreatePayment(ctx, &billings.PostPaymentRequest{
		UserId:                  user.Id,
		CompanyId:               companyId,
		OrganizationId:          organizationId,
		PlanId:                  request.PlanId,
		PaymentMethod:           request.PaymentMethod,
		ProviderPaymentMethodId: request.ProviderPaymentMethodId,
		Token:                   request.Token,
		Installments:            request.Installments,
		IssuerId:                request.IssuerId,
		Payer:                   request.Payer,
		IdempotencyKey:          request.IdempotencyKey,
		DeviceId:                request.DeviceId,
		IpAddress:               request.IpAddress,
		CallbackUrl:             request.CallbackUrl,
		CouponCode:              request.CouponCode,
		Metadata:                request.Metadata,
	})
}

func (s *service) CancelPayment(ctx context.Context, request *model.PaymentRequest) (*billings.PostCancelPaymentResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.billings.CancelPayment(ctx, &billings.PostCancelPaymentRequest{
		Id:             request.Id,
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostPlan(ctx context.Context, request *model.PlanRequest) (*billings.PlanComplete, error) {
	return s.billings.CreatePlan(ctx, &billings.PlanComplete{
		Name:                 request.Name,
		Slug:                 request.Slug,
		Description:          request.Description,
		PriceInCents:         request.PriceInCents,
		Currency:             request.Currency,
		BillingInterval:      request.BillingInterval,
		BillingIntervalCount: request.BillingIntervalCount,
		TrialDays:            request.TrialDays,
		Type:                 request.Type,
		Status:               request.Status,
		IsPublic:             request.IsPublic,
		Metadata:             request.Metadata,
		CallbackUrl:          request.CallbackUrl,
		Provider:             request.Provider,
	})
}

func (s *service) GetListPlans(ctx context.Context, _ *model.PlanRequest) (*billings.GetListPlansResponse, error) {
	return s.billings.GetListPlans(ctx, &billings.GetListPlansRequest{})
}

func (s *service) GetListUserPayments(ctx context.Context, _ *model.PaymentRequest) (*billings.GetListUserPaymentsResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.billings.GetListUserPayments(ctx, &billings.GetListUserPaymentsRequest{
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) GetListAllUsers(ctx context.Context, _ *model.ClientRequest) (*client.ListAllClientsResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)
	companyId, organizationId := ownership(ctx)

	return s.clients.ListAll(filter.Outgoing(ctx), &client.ListAllClientsRequest{
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostCreateBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*client.ClientComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.clients.CreateBusinessClient(ctx, &client.CreateBusinessClientRequest{
		ScopeId:                  r.ScopeId,
		Name:                     r.Name,
		Email:                    r.Email,
		Password:                 r.Password,
		CompanyId:                companyId,
		OrganizationId:           organizationId,
		Role:                     r.Role,
		SendEmailAccess:          r.SendEmailAccess,
		GenerateRandomicPassword: r.GenerateRandomicPassword,
	})
}

func (s *service) GetListAllScopes(ctx context.Context, _ *model.ScopeRequest) (*client.ListAllScopesResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)
	companyId, organizationId := ownership(ctx)

	return s.clients.ListAllScopes(filter.Outgoing(ctx), &client.ListAllScopesRequest{
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) DeleteBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*structpb.Struct, error) {
	companyId, organizationId := ownership(ctx)

	return s.clients.DeleteBusinessClient(ctx, &client.DeleteBusinessClientRequest{
		UserId:         r.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PatchBusinessClient(ctx context.Context, r *model.BusinessClientRequest) (*client.ClientComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.clients.PatchBusinessClient(ctx, BuildPatchBusinessClient(companyId, organizationId, r))
}

func (s *service) ListAutomations(ctx context.Context, r *model.AutomationsRequest) (*automations.ListAutomationsResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.automations.ListAutomations(ctx, &automations.ListAutomationsRequest{
		UserId: user.Id,
		Type:   r.Type,
		Status: r.Status,
	})
}

func (s *service) GetAndValidateCoupon(ctx context.Context, request *model.CouponRequest) (*coupons.CouponComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	return s.coupons.GetAndValidateCoupon(ctx, &coupons.GetAndValidateCouponRequest{
		Code:   request.Code,
		UserId: user.Id,
	})
}

func (s *service) PostCoupon(ctx context.Context, request *model.PostCouponRequest) (*coupons.CouponComplete, error) {
	return s.coupons.CreateCoupon(ctx, &coupons.PostCouponRequest{
		Code:               request.Code,
		DiscountPercentage: int32(request.DiscountPercentage),
		MaxIntCents:        int32(request.MaxIntCents),
		StartAt:            request.StartAt,
		EndAt:              request.EndAt,
		Active:             request.Active,
		MaxUses:            int32(request.MaxUses),
		Objective:          request.Objective,
		ForFirstBuy:        request.ForFirstBuy,
	})
}

func (s *service) PostCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	companyId, organizationId := ownership(ctx)

	url := request.Url

	if request.File != nil {
		if request.FileSize > maxCertificateBytes {
			level.Error(s.logger).Log("during", "post_certificate > size", "error", ErrCertificateTooLarge, "bytes", request.FileSize)
			return nil, ErrCertificateTooLarge
		}

		head := make([]byte, 512)
		read, err := io.ReadFull(request.File, head)

		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			level.Error(s.logger).Log("during", "post_certificate > read", "error", err)
			return nil, err
		}

		head = head[:read]

		if read == 0 {
			level.Error(s.logger).Log("during", "post_certificate > read", "error", ErrCertificateEmptyFile)
			return nil, ErrCertificateEmptyFile
		}

		detected := http.DetectContentType(head)

		if !slices.Contains(allowedCertificateTypes, detected) {
			level.Error(s.logger).Log("during", "post_certificate > content type", "error", ErrCertificateTypeNotAllowed, "detected", detected)
			return nil, ErrCertificateTypeNotAllowed
		}

		url, err = s.uploadFile(io.MultiReader(bytes.NewReader(head), request.File), request.FileName)

		if err != nil {
			level.Error(s.logger).Log("during", "post_certificate > upload", "error", err)
			return nil, err
		}

		if url == "" {
			level.Error(s.logger).Log("during", "post_certificate > upload", "error", ErrCertificateUploadFailed)
			return nil, ErrCertificateUploadFailed
		}
	}

	return s.certificates.CreateCertificate(ctx, &certificates.PostCertificateRequest{
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		DocumentCnpj:   request.DocumentCnpj,
		Type:           certificates.CertificateType(request.Type),
		Title:          request.Title,
		Url:            url,
		ExpiresAt:      request.ExpiresAt,
	})
}

func (s *service) GetCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.certificates.GetCertificate(ctx, &certificates.GetCertificateRequest{
		Id:             request.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PatchCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificateComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.certificates.UpdateCertificate(ctx, &certificates.UpdateCertificateRequest{
		Certificate: &certificates.CertificateComplete{
			Id:             request.Id,
			CompanyId:      companyId,
			OrganizationId: organizationId,
			DocumentCnpj:   request.DocumentCnpj,
			Type:           certificates.CertificateType(request.Type),
			Title:          request.Title,
			Url:            request.Url,
			ExpiresAt:      request.ExpiresAt,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: request.UpdateMask},
	})
}

func (s *service) DeleteCertificate(ctx context.Context, request *model.CertificateRequest) (*certificates.DeleteCertificateResponse, error) {
	companyId, organizationId := ownership(ctx)

	return s.certificates.DeleteCertificate(ctx, &certificates.DeleteCertificateRequest{
		Id:             request.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) ListCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.ListCertificatesResponse, error) {
	companyId, organizationId := ownership(ctx)

	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)

	return s.certificates.ListCertificates(filter.Outgoing(ctx), &certificates.ListCertificatesRequest{
		CompanyId:      companyId,
		OrganizationId: organizationId,
		DocumentCnpj:   request.DocumentCnpj,
		Type:           certificates.CertificateType(request.Type),
	})
}

func (s *service) PostExportCertificates(ctx context.Context, request *model.CertificateRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.certificates.PostExportCertificates(ctx, &certificates.PostExportCertificatesRequest{
		CompanyId:      companyId,
		OrganizationId: organizationId,
		DocumentCnpj:   request.DocumentCnpj,
		Type:           certificates.CertificateType(request.Type),
	})
}

func (s *service) GetCertificatesAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*certificates.CertificatesAsyncRunnerComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.certificates.GetAsyncRunner(ctx, &certificates.GetAsyncRunnerRequest{
		RunnerId:       request.RunnerId,
		RunnerType:     certificates.RunnerType_EXPORT_CERTIFICATES_RUNNER,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostExportReports(ctx context.Context, request *model.HoldingRequest) (*bids.BidsAsyncRunnerComplete, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	request.UserId = user.Id

	companyId, organizationId := ownership(ctx)

	return s.bids.PostExportReports(ctx, &bids.PostExportReportsRequest{
		UserId:           request.UserId,
		CompanyId:        companyId,
		OrganizationId:   organizationId,
		PublicationMonth: request.PublicationMonth,
		PublicationYear:  request.PublicationYear,
	})
}

func (s *service) GetReportsAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.bids.GetAsyncRunner(ctx, &bids.GetAsyncRunnerRequest{
		RunnerId:       request.RunnerId,
		RunnerType:     bids.RunnerType_EXPORT_REPORTS_RUNNER,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostExportSearch(ctx context.Context, _ *model.BidRequest) (*bids.BidsAsyncRunnerComplete, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)

	companyId, organizationId := ownership(ctx)

	return s.bids.PostExportSearch(filter.Outgoing(ctx), &bids.PostExportSearchRequest{
		UserId:         user.Id,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) GetSearchAsyncRunner(ctx context.Context, request *model.RunnerRequest) (*bids.BidsAsyncRunnerComplete, error) {
	companyId, organizationId := ownership(ctx)

	return s.bids.GetAsyncRunner(ctx, &bids.GetAsyncRunnerRequest{
		RunnerId:       request.RunnerId,
		RunnerType:     bids.RunnerType_EXPORT_SEARCH_RUNNER,
		CompanyId:      companyId,
		OrganizationId: organizationId,
	})
}

func (s *service) PostLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	return s.limits.PostLimit(ctx, &limits.PostLimitRequest{
		OrganizationId: request.OrganizationId,
		CompanyId:      request.CompanyId,
		PayerRole:      int32(request.PayerRole),
		Type:           request.Type,
		Quantity:       request.Quantity,
		Period:         request.Period,
		Unlimited:      request.Unlimited,
	})
}

func (s *service) PatchLimit(ctx context.Context, request *model.LimitRequest) (*limits.LimitsComplete, error) {
	return s.limits.UpdateLimit(ctx, &limits.UpdateLimitRequest{
		Limit: &limits.LimitsComplete{
			Id:        request.Id,
			Quantity:  request.Quantity,
			Period:    request.Period,
			Unlimited: request.Unlimited,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: request.UpdateMask},
	})
}

func (s *service) GetListLimits(ctx context.Context, request *model.LimitRequest) (*limits.GetListLimitsResponse, error) {
	filter, _ := decode.GetFromContext[look.Params](ctx, keys.LookContext)

	return s.limits.GetListLimits(filter.Outgoing(ctx), &limits.GetListLimitsRequest{
		OrganizationId: request.OrganizationId,
		CompanyId:      request.CompanyId,
		PayerRole:      int32(request.PayerRole),
	})
}

func (s *service) DeleteLimit(ctx context.Context, request *model.LimitRequest) (*limits.DeleteLimitResponse, error) {
	return s.limits.DeleteLimit(ctx, &limits.DeleteLimitRequest{
		Id: request.Id,
	})
}

func (s *service) GetListConsumes(ctx context.Context, request *model.LimitRequest) (*limits.GetListConsumesResponse, error) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	companyId, organizationId := ownership(ctx)

	return s.limits.GetListConsumes(ctx, &limits.GetListConsumesRequest{
		OrganizationId: organizationId,
		CompanyId:      companyId,
		PayerRole:      int32(user.PayerRole),
		Type:           request.Type,
	})
}

func (s *service) uploadFile(file io.Reader, fieldName string) (string, error) {
	var resp model.CdnResponse

	_, err := s.cdnApi.R().
		SetFileReader("file", fieldName, file).
		SetDebug(true).
		SetResult(&resp).
		Post("/api/upload")

	if err != nil {
		return "", err
	}

	if resp.URL == "" {
		return "", nil
	}

	return resp.URL, nil
}

// ownership reads the company the caller is acting as, already settled and checked by the authentication middleware, and the organization above it.
func ownership(ctx context.Context) (companyId string, organizationId string) {
	user, _ := decode.GetFromContext[*client.ClientComplete](ctx, keys.ClientContext)
	companyId, _ = ctx.Value(keys.CompanyContext).(string)

	if companyId == "" {
		companyId = user.GetId()
	}

	fmt.Println(user)

	return companyId, user.GetOrganizationId()
}

var documentClient = &http.Client{Timeout: 90 * time.Second}

const maxEditalBytes = 40 * 1024 * 1024

func download(ctx context.Context, fileURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := documentClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("document download failed: %s (%d)", fileURL, resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxEditalBytes+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > maxEditalBytes {
		return nil, fmt.Errorf("edital exceeds the %d MB limit: %s", maxEditalBytes>>20, fileURL)
	}

	return raw, nil
}

var editalCompanionRe = regexp.MustCompile(`(?i)anexo|adendo|errata|retifica|esclarec|impugna|resposta|minuta|planilha|termo\s+de\s+refer`)

func findEdital(files []*bids.BidFile) *bids.BidFile {
	for _, f := range files {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(f.GetType())), "edital") {
			return f
		}
	}
	for _, f := range files {
		name := strings.ToLower(f.GetName())
		if strings.Contains(name, "edital") && !editalCompanionRe.MatchString(name) {
			return f
		}
	}
	return nil
}
