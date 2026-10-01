package endpoint

import (
	"audita-api-gateway/pkg/service"
	"audita-api-gateway/request"
	"context"

	"github.com/audita-bids/private-kit/middlewares"
	"github.com/audita-bids/private-kit/pkg/pb/protocols/pncp"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type EndpointSetup struct {
	GetAvailableLicenses     endpoint.Endpoint
	GetLicense               endpoint.Endpoint
	GetListFavoriteBid       endpoint.Endpoint
	PostFavoriteBid          endpoint.Endpoint
	DeleteFavoriteBid        endpoint.Endpoint
	PostAnalysis             endpoint.Endpoint
	GetAsyncRunner           endpoint.Endpoint
	PostCopilot              endpoint.Endpoint
	PostHoldingBid           endpoint.Endpoint
	GetListHoldingBid        endpoint.Endpoint
	PostWhitelabel           endpoint.Endpoint
	GetWhitelabel            endpoint.Endpoint
	UpdateWhitelabel         endpoint.Endpoint
	GetBids                  endpoint.Endpoint
	GetBid                   endpoint.Endpoint
	GetBidHandles            endpoint.Endpoint
	PostBidProposal          endpoint.Endpoint
	UpdateBidProposal        endpoint.Endpoint
	PostPayment              endpoint.Endpoint
	CancelPayment            endpoint.Endpoint
	PostPlan                 endpoint.Endpoint
	GetListPlans             endpoint.Endpoint
	GetListUserPayments      endpoint.Endpoint
	GetListAllUsers          endpoint.Endpoint
	PostCreateBusinessClient endpoint.Endpoint
	GetListAllScopes         endpoint.Endpoint
	DeleteBusinessClient     endpoint.Endpoint
	PatchBusinessClient      endpoint.Endpoint
	ListAutomations          endpoint.Endpoint
	GetAndValidateCoupon     endpoint.Endpoint
	PostCoupon               endpoint.Endpoint
	PostCertificate          endpoint.Endpoint
	GetCertificate           endpoint.Endpoint
	PatchCertificate         endpoint.Endpoint
	DeleteCertificate        endpoint.Endpoint
	ListCertificates         endpoint.Endpoint
	PostExportCertificates   endpoint.Endpoint
	GetCertificatesRunner    endpoint.Endpoint
	PostExportReports        endpoint.Endpoint
	GetReportsRunner         endpoint.Endpoint
	PostExportSearch         endpoint.Endpoint
	GetSearchRunner          endpoint.Endpoint
	PostLimit                endpoint.Endpoint
	PatchLimit               endpoint.Endpoint
	GetListLimits            endpoint.Endpoint
	DeleteLimit              endpoint.Endpoint
	GetListConsumes          endpoint.Endpoint
}

func NewEndpointSetup(s service.Service, logger log.Logger) *EndpointSetup {
	var GetAvailableLicenses endpoint.Endpoint
	var GetLicense endpoint.Endpoint
	var GetListFavoriteBid endpoint.Endpoint
	var PostFavoriteBid endpoint.Endpoint
	var PostAnalysis endpoint.Endpoint
	var GetAsyncRunner endpoint.Endpoint
	var PostCopilot endpoint.Endpoint
	var PostHoldingBid endpoint.Endpoint
	var GetListHoldingBid endpoint.Endpoint
	var PostWhitelabel endpoint.Endpoint
	var GetWhitelabel endpoint.Endpoint
	var UpdateWhitelabel endpoint.Endpoint
	var GetBids endpoint.Endpoint
	var GetBid endpoint.Endpoint
	var GetBidHandles endpoint.Endpoint
	var PostBidProposal endpoint.Endpoint
	var UpdateBidProposal endpoint.Endpoint
	var PostPayment endpoint.Endpoint
	var CancelPayment endpoint.Endpoint
	var PostPlan endpoint.Endpoint
	var GetListPlans endpoint.Endpoint
	var GetListUserPayments endpoint.Endpoint
	var GetListAllUsers endpoint.Endpoint
	var PostCreateBusinessClient endpoint.Endpoint
	var DeleteBusinessClient endpoint.Endpoint
	var GetListAllScopes endpoint.Endpoint
	var PatchBusinessClient endpoint.Endpoint
	var ListAutomations endpoint.Endpoint
	var GetAndValidateCoupon endpoint.Endpoint
	var PostCoupon endpoint.Endpoint
	var PostCertificate endpoint.Endpoint
	var GetCertificate endpoint.Endpoint
	var PatchCertificate endpoint.Endpoint
	var DeleteCertificate endpoint.Endpoint
	var ListCertificates endpoint.Endpoint
	var PostExportCertificates endpoint.Endpoint
	var GetCertificatesRunner endpoint.Endpoint
	var PostExportReports endpoint.Endpoint
	var GetReportsRunner endpoint.Endpoint
	var PostExportSearch endpoint.Endpoint
	var GetSearchRunner endpoint.Endpoint
	var DeleteFavoriteBid endpoint.Endpoint
	var PostLimit endpoint.Endpoint
	var PatchLimit endpoint.Endpoint
	var GetListLimits endpoint.Endpoint
	var DeleteLimit endpoint.Endpoint
	var GetListConsumes endpoint.Endpoint

	loggingMiddleware := middlewares.EndpointLoggingMiddleware(logger, "audita-api-gateway")
	metricsMiddleware := middlewares.MetricsMiddleware("audita-api-gateway")

	{
		GetAvailableLicenses = MakeGetAvailableLicensesEndpoint(s)
		GetAvailableLicenses = loggingMiddleware("GetAvailableLicenses")(GetAvailableLicenses)
		GetAvailableLicenses = metricsMiddleware("GetAvailableLicenses")(GetAvailableLicenses)

		GetLicense = MakeGetLicenseEndpoint(s)
		GetLicense = loggingMiddleware("GetLicense")(GetLicense)
		GetLicense = metricsMiddleware("GetLicense")(GetLicense)

		GetListFavoriteBid = MakeGetListFavoriteBidEndpoint(s)
		GetListFavoriteBid = loggingMiddleware("GetListFavoriteBid")(GetListFavoriteBid)
		GetListFavoriteBid = metricsMiddleware("GetListFavoriteBid")(GetListFavoriteBid)

		PostFavoriteBid = MakePostFavoriteBidEndpoint(s)
		PostFavoriteBid = loggingMiddleware("PostFavoriteBid")(PostFavoriteBid)
		PostFavoriteBid = metricsMiddleware("PostFavoriteBid")(PostFavoriteBid)

		PostAnalysis = MakePostAnalysisEndpoint(s)
		PostAnalysis = loggingMiddleware("PostAnalysis")(PostAnalysis)
		PostAnalysis = metricsMiddleware("PostAnalysis")(PostAnalysis)

		GetAsyncRunner = MakeGetAsyncRunnerEndpoint(s)
		GetAsyncRunner = loggingMiddleware("GetAsyncRunner")(GetAsyncRunner)
		GetAsyncRunner = metricsMiddleware("GetAsyncRunner")(GetAsyncRunner)

		PostCopilot = MakePostCopilotEndpoint(s)
		PostCopilot = loggingMiddleware("PostCopilot")(PostCopilot)
		PostCopilot = metricsMiddleware("PostCopilot")(PostCopilot)

		PostHoldingBid = MakePostHoldingBidEndpoint(s)
		PostHoldingBid = loggingMiddleware("PostHoldingBid")(PostHoldingBid)
		PostHoldingBid = metricsMiddleware("PostHoldingBid")(PostHoldingBid)

		GetListHoldingBid = MakeGetListHoldingBidEndpoint(s)
		GetListHoldingBid = loggingMiddleware("GetListHoldingBid")(GetListHoldingBid)
		GetListHoldingBid = metricsMiddleware("GetListHoldingBid")(GetListHoldingBid)

		PostWhitelabel = MakePostWhitelabelEndpoint(s)
		PostWhitelabel = loggingMiddleware("PostWhitelabel")(PostWhitelabel)
		PostWhitelabel = metricsMiddleware("PostWhitelabel")(PostWhitelabel)

		GetWhitelabel = MakeGetWhitelabelEndpoint(s)
		GetWhitelabel = loggingMiddleware("GetWhitelabel")(GetWhitelabel)
		GetWhitelabel = metricsMiddleware("GetWhitelabel")(GetWhitelabel)

		UpdateWhitelabel = MakeUpdateWhitelabelEndpoint(s)
		UpdateWhitelabel = loggingMiddleware("UpdateWhitelabel")(UpdateWhitelabel)
		UpdateWhitelabel = metricsMiddleware("UpdateWhitelabel")(UpdateWhitelabel)

		GetBids = MakeGetBidsEndpoint(s)
		GetBids = loggingMiddleware("GetBids")(GetBids)
		GetBids = metricsMiddleware("GetBids")(GetBids)

		GetBid = MakeGetBidEndpoint(s)
		GetBid = loggingMiddleware("GetBid")(GetBid)
		GetBid = metricsMiddleware("GetBid")(GetBid)

		GetBidHandles = MakeGetBidHandlesEndpoint(s)
		GetBidHandles = loggingMiddleware("GetBidHandles")(GetBidHandles)
		GetBidHandles = metricsMiddleware("GetBidHandles")(GetBidHandles)

		PostBidProposal = MakePostBidProposalEndpoint(s)
		PostBidProposal = loggingMiddleware("PostBidProposal")(PostBidProposal)
		PostBidProposal = metricsMiddleware("PostBidProposal")(PostBidProposal)

		UpdateBidProposal = MakeUpdateBidProposalEndpoint(s)
		UpdateBidProposal = loggingMiddleware("UpdateBidProposal")(UpdateBidProposal)
		UpdateBidProposal = metricsMiddleware("UpdateBidProposal")(UpdateBidProposal)

		PostPayment = MakePostPaymentEndpoint(s)
		PostPayment = loggingMiddleware("PostPayment")(PostPayment)
		PostPayment = metricsMiddleware("PostPayment")(PostPayment)

		CancelPayment = MakeCancelPaymentEndpoint(s)
		CancelPayment = loggingMiddleware("CancelPayment")(CancelPayment)
		CancelPayment = metricsMiddleware("CancelPayment")(CancelPayment)

		GetListPlans = MakeGetListPlansEndpoint(s)
		GetListPlans = loggingMiddleware("GetListPlans")(GetListPlans)
		GetListPlans = metricsMiddleware("GetListPlans")(GetListPlans)

		PostPlan = MakePostPlanEndpoint(s)
		PostPlan = loggingMiddleware("PostPlan")(PostPlan)
		PostPlan = metricsMiddleware("PostPlan")(PostPlan)

		GetListUserPayments = MakeGetListUserPaymentsEndpoint(s)
		GetListUserPayments = loggingMiddleware("GetListUserPayments")(GetListUserPayments)
		GetListUserPayments = metricsMiddleware("GetListUserPayments")(GetListUserPayments)

		GetListAllUsers = MakeGetListAllUsersEndpoint(s)
		GetListAllUsers = loggingMiddleware("GetListAllUsers")(GetListAllUsers)
		GetListAllUsers = metricsMiddleware("GetListAllUsers")(GetListAllUsers)

		PostCreateBusinessClient = MakePostCreateBusinessClientEndpoint(s)
		PostCreateBusinessClient = loggingMiddleware("PostCreateBusinessClient")(PostCreateBusinessClient)
		PostCreateBusinessClient = metricsMiddleware("PostCreateBusinessClient")(PostCreateBusinessClient)

		GetListAllScopes = MakeGetListAllScopesEndpoint(s)
		GetListAllScopes = loggingMiddleware("GetListAllScopes")(GetListAllScopes)
		GetListAllScopes = metricsMiddleware("GetListAllScopes")(GetListAllScopes)

		DeleteBusinessClient = MakeDeleteBusinessClientEndpoint(s)
		DeleteBusinessClient = loggingMiddleware("DeleteBusinessClient")(DeleteBusinessClient)
		DeleteBusinessClient = metricsMiddleware("DeleteBusinessClient")(DeleteBusinessClient)

		PatchBusinessClient = MakePatchBusinessClientEndpoint(s)
		PatchBusinessClient = loggingMiddleware("PatchBusinessClient")(PatchBusinessClient)
		PatchBusinessClient = metricsMiddleware("PatchBusinessClient")(PatchBusinessClient)

		ListAutomations = MakeListAutomationsEndpoint(s)
		ListAutomations = loggingMiddleware("ListAutomations")(ListAutomations)
		ListAutomations = metricsMiddleware("ListAutomations")(ListAutomations)

		GetAndValidateCoupon = MakeGetAndValidateCouponEndpoint(s)
		GetAndValidateCoupon = loggingMiddleware("GetAndValidateCoupon")(GetAndValidateCoupon)
		GetAndValidateCoupon = metricsMiddleware("GetAndValidateCoupon")(GetAndValidateCoupon)

		PostCoupon = MakePostCouponEndpoint(s)
		PostCoupon = loggingMiddleware("PostCoupon")(PostCoupon)
		PostCoupon = metricsMiddleware("PostCoupon")(PostCoupon)

		PostCertificate = MakePostCertificateEndpoint(s)
		PostCertificate = loggingMiddleware("PostCertificate")(PostCertificate)
		PostCertificate = metricsMiddleware("PostCertificate")(PostCertificate)

		GetCertificate = MakeGetCertificateEndpoint(s)
		GetCertificate = loggingMiddleware("GetCertificate")(GetCertificate)
		GetCertificate = metricsMiddleware("GetCertificate")(GetCertificate)

		PatchCertificate = MakePatchCertificateEndpoint(s)
		PatchCertificate = loggingMiddleware("PatchCertificate")(PatchCertificate)
		PatchCertificate = metricsMiddleware("PatchCertificate")(PatchCertificate)

		DeleteCertificate = MakeDeleteCertificateEndpoint(s)
		DeleteCertificate = loggingMiddleware("DeleteCertificate")(DeleteCertificate)
		DeleteCertificate = metricsMiddleware("DeleteCertificate")(DeleteCertificate)

		ListCertificates = MakeListCertificatesEndpoint(s)
		ListCertificates = loggingMiddleware("ListCertificates")(ListCertificates)
		ListCertificates = metricsMiddleware("ListCertificates")(ListCertificates)

		PostExportCertificates = MakePostExportCertificatesEndpoint(s)
		PostExportCertificates = loggingMiddleware("PostExportCertificates")(PostExportCertificates)
		PostExportCertificates = metricsMiddleware("PostExportCertificates")(PostExportCertificates)

		GetCertificatesRunner = MakeGetCertificatesRunnerEndpoint(s)
		GetCertificatesRunner = loggingMiddleware("GetCertificatesRunner")(GetCertificatesRunner)
		GetCertificatesRunner = metricsMiddleware("GetCertificatesRunner")(GetCertificatesRunner)

		PostExportReports = MakePostExportReportsEndpoint(s)
		PostExportReports = loggingMiddleware("PostExportReports")(PostExportReports)
		PostExportReports = metricsMiddleware("PostExportReports")(PostExportReports)

		GetReportsRunner = MakeGetReportsRunnerEndpoint(s)
		GetReportsRunner = loggingMiddleware("GetReportsRunner")(GetReportsRunner)
		GetReportsRunner = metricsMiddleware("GetReportsRunner")(GetReportsRunner)

		PostExportSearch = MakePostExportSearchEndpoint(s)
		PostExportSearch = loggingMiddleware("PostExportSearch")(PostExportSearch)
		PostExportSearch = metricsMiddleware("PostExportSearch")(PostExportSearch)

		GetSearchRunner = MakeGetSearchRunnerEndpoint(s)
		GetSearchRunner = loggingMiddleware("GetSearchRunner")(GetSearchRunner)
		GetSearchRunner = metricsMiddleware("GetSearchRunner")(GetSearchRunner)

		DeleteFavoriteBid = MakeDeleteFavoriteBidEndpoint(s)
		DeleteFavoriteBid = loggingMiddleware("DeleteFavoriteBid")(DeleteFavoriteBid)
		DeleteFavoriteBid = metricsMiddleware("DeleteFavoriteBid")(DeleteFavoriteBid)

		PostLimit = MakePostLimitEndpoint(s)
		PostLimit = loggingMiddleware("PostLimit")(PostLimit)
		PostLimit = metricsMiddleware("PostLimit")(PostLimit)

		PatchLimit = MakePatchLimitEndpoint(s)
		PatchLimit = loggingMiddleware("PatchLimit")(PatchLimit)
		PatchLimit = metricsMiddleware("PatchLimit")(PatchLimit)

		GetListLimits = MakeGetListLimitsEndpoint(s)
		GetListLimits = loggingMiddleware("GetListLimits")(GetListLimits)
		GetListLimits = metricsMiddleware("GetListLimits")(GetListLimits)

		DeleteLimit = MakeDeleteLimitEndpoint(s)
		DeleteLimit = loggingMiddleware("DeleteLimit")(DeleteLimit)
		DeleteLimit = metricsMiddleware("DeleteLimit")(DeleteLimit)

		GetListConsumes = MakeGetListConsumesEndpoint(s)
		GetListConsumes = loggingMiddleware("GetListConsumes")(GetListConsumes)
		GetListConsumes = metricsMiddleware("GetListConsumes")(GetListConsumes)
	}

	return &EndpointSetup{
		GetAvailableLicenses:     GetAvailableLicenses,
		GetLicense:               GetLicense,
		GetListFavoriteBid:       GetListFavoriteBid,
		PostFavoriteBid:          PostFavoriteBid,
		PostAnalysis:             PostAnalysis,
		GetAsyncRunner:           GetAsyncRunner,
		PostCopilot:              PostCopilot,
		PostHoldingBid:           PostHoldingBid,
		GetListHoldingBid:        GetListHoldingBid,
		PostWhitelabel:           PostWhitelabel,
		GetWhitelabel:            GetWhitelabel,
		UpdateWhitelabel:         UpdateWhitelabel,
		GetBids:                  GetBids,
		GetBid:                   GetBid,
		GetBidHandles:            GetBidHandles,
		PostBidProposal:          PostBidProposal,
		UpdateBidProposal:        UpdateBidProposal,
		PostPayment:              PostPayment,
		CancelPayment:            CancelPayment,
		PostPlan:                 PostPlan,
		GetListPlans:             GetListPlans,
		GetListUserPayments:      GetListUserPayments,
		GetListAllUsers:          GetListAllUsers,
		PostCreateBusinessClient: PostCreateBusinessClient,
		GetListAllScopes:         GetListAllScopes,
		DeleteBusinessClient:     DeleteBusinessClient,
		PatchBusinessClient:      PatchBusinessClient,
		ListAutomations:          ListAutomations,
		GetAndValidateCoupon:     GetAndValidateCoupon,
		PostCoupon:               PostCoupon,
		PostCertificate:          PostCertificate,
		GetCertificate:           GetCertificate,
		PatchCertificate:         PatchCertificate,
		DeleteCertificate:        DeleteCertificate,
		ListCertificates:         ListCertificates,
		PostExportCertificates:   PostExportCertificates,
		GetCertificatesRunner:    GetCertificatesRunner,
		PostExportReports:        PostExportReports,
		GetReportsRunner:         GetReportsRunner,
		PostExportSearch:         PostExportSearch,
		GetSearchRunner:          GetSearchRunner,
		DeleteFavoriteBid:        DeleteFavoriteBid,
		PostLimit:                PostLimit,
		PatchLimit:               PatchLimit,
		GetListLimits:            GetListLimits,
		DeleteLimit:              DeleteLimit,
		GetListConsumes:          GetListConsumes,
	}
}

func MakeGetAvailableLicensesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*pncp.PncpAvailableLicenseRequest)
		rpcRequest := &pncp.PncpAvailableLicenseRequest{
			CodigoMunicipioIbge:         r.CodigoMunicipioIbge,
			DataInicial:                 r.DataInicial,
			DataFinal:                   r.DataFinal,
			CodigoModalidadeContratacao: r.CodigoModalidadeContratacao,
			Pagina:                      r.Pagina,
			TamanhoPagina:               r.TamanhoPagina,
			Uf:                          r.Uf,
		}

		fc, err := s.GetAvailableLicenses(ctx, rpcRequest)

		if err != nil {
			return nil, err
		}
		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetLicenseEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*pncp.PncpFindLicenseRequest)
		fc, err := s.GetLicense(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetListFavoriteBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		r := request.(*model.FavoriteBidRequest)

		fc, err := s.GetListFavoriteBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
			Total: fc.Total,
		}, nil
	}
}

func MakePostFavoriteBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		r := request.(*model.FavoriteBidRequest)

		fc, err := s.PostFavoriteBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostAnalysisEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.AnalysisRequest)

		fc, err := s.PostAnalysis(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetAsyncRunnerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.RunnerRequest)

		fc, err := s.GetAsyncRunner(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostCopilotEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CopilotRequest)

		fc, err := s.PostCopilot(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostHoldingBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.HoldingRequest)

		fc, err := s.PostHoldingBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListHoldingBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.HoldingRequest)

		fc, err := s.GetListHoldingBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.GetWhitelabel(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.PostWhitelabel(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeUpdateWhitelabelEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.WhitelabelRequest)

		fc, err := s.UpdateWhitelabel(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetBidsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBids(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
		}, nil
	}
}

func MakeGetBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetBidHandlesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.GetBidHandles(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostBidProposalEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ProposalRequest)

		fc, err := s.PostBidProposal(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeUpdateBidProposalEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ProposalRequest)

		fc, err := s.UpdateBidProposal(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostPaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PaymentRequest)

		fc, err := s.PostPayment(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeCancelPaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PaymentRequest)

		fc, err := s.CancelPayment(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListPlansEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PlanRequest)

		fc, err := s.GetListPlans(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc.Plans,
			Total: fc.Total,
		}, nil
	}
}

func MakePostPlanEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PlanRequest)

		fc, err := s.PostPlan(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListUserPaymentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PaymentRequest)

		fc, err := s.GetListUserPayments(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListAllUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ClientRequest)

		fc, err := s.GetListAllUsers(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostCreateBusinessClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BusinessClientRequest)

		fc, err := s.PostCreateBusinessClient(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListAllScopesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.ScopeRequest)

		fc, err := s.GetListAllScopes(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeDeleteBusinessClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BusinessClientRequest)

		fc, err := s.DeleteBusinessClient(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePatchBusinessClientEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BusinessClientRequest)

		fc, err := s.PatchBusinessClient(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeListAutomationsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.AutomationsRequest)

		fc, err := s.ListAutomations(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostCertificateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.PostCertificate(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetCertificateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.GetCertificate(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePatchCertificateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.PatchCertificate(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeDeleteCertificateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.DeleteCertificate(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeListCertificatesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.ListCertificates(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
			Total: fc.Total,
		}, nil
	}
}

func MakePostExportCertificatesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CertificateRequest)

		fc, err := s.PostExportCertificates(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetCertificatesRunnerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.RunnerRequest)

		fc, err := s.GetCertificatesAsyncRunner(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostExportReportsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.HoldingRequest)

		fc, err := s.PostExportReports(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetReportsRunnerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.RunnerRequest)

		fc, err := s.GetReportsAsyncRunner(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostExportSearchEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.BidRequest)

		fc, err := s.PostExportSearch(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetSearchRunnerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.RunnerRequest)

		fc, err := s.GetSearchAsyncRunner(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostCouponEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.PostCouponRequest)

		fc, err := s.PostCoupon(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetAndValidateCouponEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.CouponRequest)

		fc, err := s.GetAndValidateCoupon(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeDeleteFavoriteBidEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.FavoriteBidRequest)

		fc, err := s.DeleteFavoriteBid(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePostLimitEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.LimitRequest)

		fc, err := s.PostLimit(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakePatchLimitEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.LimitRequest)

		fc, err := s.PatchLimit(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListLimitsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.LimitRequest)

		fc, err := s.GetListLimits(ctx, r)

		if err != nil {
			return nil, err
		}

		return &Resp{
			Items: fc,
			Total: fc.Total,
		}, nil
	}
}

func MakeDeleteLimitEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.LimitRequest)

		fc, err := s.DeleteLimit(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

func MakeGetListConsumesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*model.LimitRequest)

		fc, err := s.GetListConsumes(ctx, r)

		if err != nil {
			return nil, err
		}

		return fc, nil
	}
}

type Resp struct {
	Items  interface{} `json:"items,omitempty"`
	Total  int64       `json:"total,omitempty"`
	Cursor string      `json:"cursor,omitempty"`
}
