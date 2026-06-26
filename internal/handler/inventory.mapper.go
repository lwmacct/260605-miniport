package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToHostDTO(host service.InventoryHost) HostDTO {
	return HostDTO{
		ID:          host.ID,
		IP:          host.IP,
		Name:        host.Name,
		Network:     host.Network,
		Environment: host.Environment,
		Notes:       host.Notes,
		CreatedAt:   utilHTTPTime(host.CreatedAt),
		UpdatedAt:   utilHTTPTime(host.UpdatedAt),
	}
}

func ToHostDTOs(hosts []service.InventoryHost) []HostDTO {
	out := make([]HostDTO, 0, len(hosts))
	for _, host := range hosts {
		out = append(out, ToHostDTO(host))
	}
	return out
}

func ToPortSlotDTO(slot service.InventoryPortSlot) PortSlotDTO {
	return PortSlotDTO{
		ID:          slot.ID,
		PortGroupID: slot.PortGroupID,
		Port:        slot.Port,
		Name:        slot.Name,
		Protocol:    slot.Protocol,
		Purpose:     slot.Purpose,
		Status:      slot.Status,
		Notes:       slot.Notes,
		CreatedAt:   utilHTTPTime(slot.CreatedAt),
		UpdatedAt:   utilHTTPTime(slot.UpdatedAt),
	}
}

func ToComponentDTO(component service.InventoryComponent) ComponentDTO {
	return ComponentDTO{
		ID:          component.ID,
		PortGroupID: component.PortGroupID,
		Name:        component.Name,
		Type:        component.Type,
		URL:         component.URL,
		Version:     component.Version,
		Notes:       component.Notes,
		CreatedAt:   utilHTTPTime(component.CreatedAt),
		UpdatedAt:   utilHTTPTime(component.UpdatedAt),
	}
}

func ToRepositoryDTO(repository service.InventoryRepositoryRef) RepositoryDTO {
	return RepositoryDTO{
		ID:          repository.ID,
		PortGroupID: repository.PortGroupID,
		Name:        repository.Name,
		URL:         repository.URL,
		Kind:        repository.Kind,
		Notes:       repository.Notes,
		CreatedAt:   utilHTTPTime(repository.CreatedAt),
		UpdatedAt:   utilHTTPTime(repository.UpdatedAt),
	}
}

func ToPortGroupDTO(group service.PortGroupView) PortGroupDTO {
	out := PortGroupDTO{
		ID:            group.ID,
		HostID:        group.HostID,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		ServiceName:   group.ServiceName,
		ContainerName: group.ContainerName,
		DindHost:      group.DindHost,
		Status:        group.Status,
		Owner:         group.Owner,
		Tags:          group.Tags,
		Notes:         group.Notes,
		CreatedAt:     utilHTTPTime(group.CreatedAt),
		UpdatedAt:     utilHTTPTime(group.UpdatedAt),
		Slots:         make([]PortSlotDTO, 0, len(group.Slots)),
		Components:    make([]ComponentDTO, 0, len(group.Components)),
		Repositories:  make([]RepositoryDTO, 0, len(group.Repositories)),
	}
	if group.Host != nil {
		host := ToHostDTO(*group.Host)
		out.Host = &host
	}
	for _, slot := range group.Slots {
		out.Slots = append(out.Slots, ToPortSlotDTO(slot))
	}
	for _, component := range group.Components {
		out.Components = append(out.Components, ToComponentDTO(component))
	}
	for _, repository := range group.Repositories {
		out.Repositories = append(out.Repositories, ToRepositoryDTO(repository))
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
