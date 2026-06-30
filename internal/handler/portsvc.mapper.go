package handler

import (
	"github.com/lwmacct/260630-go-hsr-shared/pkg/identity"

	"github.com/lwmacct/260605-miniport/internal/service"
)

func ToPortsvcActor(user *identity.Principal) service.PortsvcActor {
	actor := service.PortsvcActor{}
	if user == nil {
		return actor
	}
	actor.UserID = user.ID
	actor.Username = user.Username
	actor.Admin = user.Admin
	return actor
}

func ToServicePayload(input ServicePayloadDTO) service.ServicePayload {
	out := service.ServicePayload{
		UserID:           input.UserID,
		PortAllocationID: input.PortAllocationID,
		Name:             input.Name,
		ProjectName:      input.ProjectName,
		DindIP:           input.DindIP,
		DindContainer:    input.DindContainer,
		Status:           input.Status,
		Owner:            input.Owner,
		Tags:             input.Tags,
		Notes:            input.Notes,
		Repositories:     make([]service.RepositoryPayload, 0, len(input.Repositories)),
		Dependencies:     make([]service.DependencyPayload, 0, len(input.Dependencies)),
	}
	for _, repo := range input.Repositories {
		out.Repositories = append(out.Repositories, service.RepositoryPayload(repo))
	}
	for _, dep := range input.Dependencies {
		out.Dependencies = append(out.Dependencies, service.DependencyPayload(dep))
	}
	return out
}

func ToPortAllocationPayload(input PortAllocationPayloadDTO) service.PortAllocationPayload {
	return service.PortAllocationPayload{
		UserID:    input.UserID,
		PortStart: input.PortStart,
		PortEnd:   input.PortEnd,
		Status:    input.Status,
		Notes:     input.Notes,
	}
}

func ToServiceBatchDeleteInput(actor service.PortsvcActor, ids []int64) service.ServiceBatchDeleteInput {
	return service.ServiceBatchDeleteInput{Actor: actor, IDs: ids}
}

func ToPortAllocationDTO(group service.PortAllocation) PortAllocationDTO {
	return PortAllocationDTO{
		ID:        group.ID,
		UserID:    group.UserID,
		Username:  group.Username,
		PortStart: group.PortStart,
		PortEnd:   group.PortEnd,
		Status:    group.Status,
		Notes:     group.Notes,
		CreatedAt: utilHTTPTime(group.CreatedAt),
		UpdatedAt: utilHTTPTime(group.UpdatedAt),
	}
}

func ToPortAllocationDTOs(groups []service.PortAllocation) []PortAllocationDTO {
	out := make([]PortAllocationDTO, 0, len(groups))
	for _, group := range groups {
		out = append(out, ToPortAllocationDTO(group))
	}
	return out
}

func ToRepositoryDTO(repository service.RepositoryRef) RepositoryDTO {
	return RepositoryDTO{
		ID:        repository.ID,
		UserID:    repository.UserID,
		Name:      repository.Name,
		URL:       repository.URL,
		Kind:      repository.Kind,
		Notes:     repository.Notes,
		CreatedAt: utilHTTPTime(repository.CreatedAt),
		UpdatedAt: utilHTTPTime(repository.UpdatedAt),
	}
}

func ToDependencyDTO(component service.Dependency) DependencyDTO {
	return DependencyDTO{
		ID:        component.ID,
		UserID:    component.UserID,
		Name:      component.Name,
		Type:      component.Type,
		URL:       component.URL,
		Version:   component.Version,
		Notes:     component.Notes,
		CreatedAt: utilHTTPTime(component.CreatedAt),
		UpdatedAt: utilHTTPTime(component.UpdatedAt),
	}
}

func ToServiceDTO(item service.ServiceView) ServiceDTO {
	out := ServiceDTO{
		ID:               item.ID,
		UserID:           item.UserID,
		Username:         item.Username,
		PortAllocationID: item.PortAllocationID,
		Name:             item.Name,
		ProjectName:      item.ProjectName,
		DindIP:           item.DindIP,
		DindContainer:    item.DindContainer,
		Status:           item.Status,
		Owner:            item.Owner,
		Tags:             item.Tags,
		Notes:            item.Notes,
		CreatedAt:        utilHTTPTime(item.CreatedAt),
		UpdatedAt:        utilHTTPTime(item.UpdatedAt),
		Repositories:     make([]RepositoryDTO, 0, len(item.Repositories)),
		Dependencies:     make([]DependencyDTO, 0, len(item.Dependencies)),
	}
	if item.PortAllocation != nil {
		allocation := ToPortAllocationDTO(*item.PortAllocation)
		out.PortAllocation = &allocation
	}
	for _, repo := range item.Repositories {
		out.Repositories = append(out.Repositories, ToRepositoryDTO(repo))
	}
	for _, dep := range item.Dependencies {
		out.Dependencies = append(out.Dependencies, ToDependencyDTO(dep))
	}
	return out
}

func ToServiceDTOs(items []service.ServiceView) []ServiceDTO {
	out := make([]ServiceDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ToServiceDTO(item))
	}
	return out
}
