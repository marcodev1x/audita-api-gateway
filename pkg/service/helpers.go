package service

import (
	model "audita-api-gateway/request"

	"github.com/audita-bids/private-kit/pkg/pb/protocols/client"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func BuildPatchBusinessClient(companyId string, organizationId string, r *model.BusinessClientRequest) *client.PatchBusinessClientRequest {
	merged := &client.ClientComplete{
		Id:      r.Id,
		ScopeId: r.ScopeId,
		Status:  r.Status,
		Name:    r.Name,
		Email:   r.Email,
	}

	return &client.PatchBusinessClientRequest{
		Client:         merged,
		CompanyId:      companyId,
		OrganizationId: organizationId,
		UpdateMask:     &fieldmaskpb.FieldMask{Paths: r.UpdateMask},
	}
}
