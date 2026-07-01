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
	actor.OwnerSubject = user.Subject
	actor.OwnerName = user.Username
	actor.Admin = user.Admin
	return actor
}

func ToHostPayload(input HostPayloadDTO) service.HostPayload {
	return service.HostPayload(input)
}

func ToPortSlotPayload(input PortSlotPayloadDTO) service.PortSlotPayload {
	return service.PortSlotPayload(input)
}

func ToPortGroupPayload(input PortGroupPayloadDTO) service.PortGroupPayload {
	out := service.PortGroupPayload{
		OwnerSubject: input.OwnerSubject,
		HostID:       input.HostID,
		PortStart:    input.PortStart,
		PortEnd:      input.PortEnd,
		ProjectName:  input.ProjectName,
		ProjectOwner: input.ProjectOwner,
		RuntimeMode:  input.RuntimeMode,
		RuntimeName:  input.RuntimeName,
		ServiceIP:    input.ServiceIP,
		Status:       input.Status,
		Tags:         input.Tags,
		Notes:        input.Notes,
		Slots:        make([]service.PortSlotPayload, 0, len(input.Slots)),
		Repositories: make([]service.RepositoryPayload, 0, len(input.Repositories)),
		Dependencies: make([]service.DependencyPayload, 0, len(input.Dependencies)),
	}
	for _, slot := range input.Slots {
		out.Slots = append(out.Slots, service.PortSlotPayload(slot))
	}
	for _, repo := range input.Repositories {
		out.Repositories = append(out.Repositories, service.RepositoryPayload(repo))
	}
	for _, dep := range input.Dependencies {
		out.Dependencies = append(out.Dependencies, service.DependencyPayload(dep))
	}
	return out
}

func ToHostDTO(host service.Host) HostDTO {
	return HostDTO{
		ID:        host.ID,
		Name:      host.Name,
		IP:        host.IP,
		Spec:      host.Spec,
		Status:    host.Status,
		Notes:     host.Notes,
		CreatedAt: utilHTTPTime(host.CreatedAt),
		UpdatedAt: utilHTTPTime(host.UpdatedAt),
	}
}

func ToHostDTOs(hosts []service.Host) []HostDTO {
	out := make([]HostDTO, 0, len(hosts))
	for _, host := range hosts {
		out = append(out, ToHostDTO(host))
	}
	return out
}

func ToPortSlotDTO(slot service.PortSlot) PortSlotDTO {
	return PortSlotDTO{
		ID:            slot.ID,
		PortGroupID:   slot.PortGroupID,
		Port:          slot.Port,
		Name:          slot.Name,
		Kind:          slot.Kind,
		Protocol:      slot.Protocol,
		ContainerName: slot.ContainerName,
		Status:        slot.Status,
		Notes:         slot.Notes,
		CreatedAt:     utilHTTPTime(slot.CreatedAt),
		UpdatedAt:     utilHTTPTime(slot.UpdatedAt),
	}
}

func ToPortSlotDTOs(slots []service.PortSlot) []PortSlotDTO {
	out := make([]PortSlotDTO, 0, len(slots))
	for _, slot := range slots {
		out = append(out, ToPortSlotDTO(slot))
	}
	return out
}

func ToPortGroupDTO(group service.PortGroupView) PortGroupDTO {
	out := PortGroupDTO{
		ID:           group.ID,
		OwnerSubject: group.OwnerSubject,
		OwnerName:    group.OwnerName,
		HostID:       group.HostID,
		PortStart:    group.PortStart,
		PortEnd:      group.PortEnd,
		ProjectName:  group.ProjectName,
		ProjectOwner: group.ProjectOwner,
		RuntimeMode:  group.RuntimeMode,
		RuntimeName:  group.RuntimeName,
		ServiceIP:    group.ServiceIP,
		Status:       group.Status,
		Tags:         group.Tags,
		Notes:        group.Notes,
		CreatedAt:    utilHTTPTime(group.CreatedAt),
		UpdatedAt:    utilHTTPTime(group.UpdatedAt),
		Slots:        ToPortSlotDTOs(group.Slots),
		Repositories: make([]RepositoryDTO, 0, len(group.Repositories)),
		Dependencies: make([]DependencyDTO, 0, len(group.Dependencies)),
	}
	if group.Host != nil {
		host := ToHostDTO(*group.Host)
		out.Host = &host
	}
	for _, repo := range group.Repositories {
		out.Repositories = append(out.Repositories, ToRepositoryDTO(repo))
	}
	for _, dep := range group.Dependencies {
		out.Dependencies = append(out.Dependencies, ToDependencyDTO(dep))
	}
	return out
}

func ToPortGroupDTOs(groups []service.PortGroupView) []PortGroupDTO {
	out := make([]PortGroupDTO, 0, len(groups))
	for _, group := range groups {
		out = append(out, ToPortGroupDTO(group))
	}
	return out
}

func ToRepositoryDTO(repository service.RepositoryRef) RepositoryDTO {
	return RepositoryDTO{
		ID:           repository.ID,
		OwnerSubject: repository.OwnerSubject,
		PortGroupID:  repository.PortGroupID,
		Name:         repository.Name,
		URL:          repository.URL,
		Kind:         repository.Kind,
		Notes:        repository.Notes,
		CreatedAt:    utilHTTPTime(repository.CreatedAt),
		UpdatedAt:    utilHTTPTime(repository.UpdatedAt),
	}
}

func ToDependencyDTO(component service.Dependency) DependencyDTO {
	return DependencyDTO{
		ID:           component.ID,
		OwnerSubject: component.OwnerSubject,
		PortGroupID:  component.PortGroupID,
		Name:         component.Name,
		Type:         component.Type,
		URL:          component.URL,
		Version:      component.Version,
		Notes:        component.Notes,
		CreatedAt:    utilHTTPTime(component.CreatedAt),
		UpdatedAt:    utilHTTPTime(component.UpdatedAt),
	}
}
