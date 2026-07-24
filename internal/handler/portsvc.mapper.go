package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToHostPayload(input HostPayloadDTO) service.HostPayload {
	return service.HostPayload(input)
}

func ToPortSlotPayload(input PortSlotPayloadDTO) service.PortSlotPayload {
	return service.PortSlotPayload(input)
}

func ToDependencyAssetPayload(input DependencyAssetPayloadDTO) service.DependencyAssetPayload {
	return service.DependencyAssetPayload(input)
}

func ToPortGroupPayload(input PortGroupPayloadDTO) service.PortGroupPayload {
	out := service.PortGroupPayload{
		HostID:           input.HostID,
		PortPrefix:       input.PortPrefix,
		EnvironmentName:  input.EnvironmentName,
		EnvironmentOwner: input.EnvironmentOwner,
		RuntimeMode:      input.RuntimeMode,
		RuntimeName:      input.RuntimeName,
		ServiceIP:        input.ServiceIP,
		Status:           input.Status,
		Tags:             input.Tags,
		Notes:            input.Notes,
		Slots:            make([]service.PortSlotPayload, 0, len(input.Slots)),
		AssetLinks:       make([]service.PortGroupAssetLinkPayload, 0, len(input.AssetLinks)),
		RepositoryLinks:  make([]service.PortGroupRepositoryLinkPayload, 0, len(input.RepositoryLinks)),
	}
	for _, slot := range input.Slots {
		out.Slots = append(out.Slots, service.PortSlotPayload(slot))
	}
	for _, link := range input.AssetLinks {
		out.AssetLinks = append(out.AssetLinks, service.PortGroupAssetLinkPayload(link))
	}
	for _, link := range input.RepositoryLinks {
		out.RepositoryLinks = append(out.RepositoryLinks, service.PortGroupRepositoryLinkPayload(link))
	}
	return out
}

func ToServiceGroupPayload(input ServiceGroupPayloadDTO) service.ServiceGroupPayload {
	out := service.ServiceGroupPayload{
		Name:        input.Name,
		Kind:        input.Kind,
		Status:      input.Status,
		Description: input.Description,
		Notes:       input.Notes,
		PortGroups:  make([]service.ServiceGroupPortGroupPayload, 0, len(input.PortGroups)),
	}
	for _, portGroup := range input.PortGroups {
		out.PortGroups = append(out.PortGroups, service.ServiceGroupPortGroupPayload(portGroup))
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

func ToDependencyAssetDTO(asset service.DependencyAsset) DependencyAssetDTO {
	return DependencyAssetDTO{
		ID:              asset.ID,
		Name:            asset.Name,
		AssetKind:       asset.AssetKind,
		AssetType:       asset.AssetType,
		Provider:        asset.Provider,
		URL:             asset.URL,
		FullName:        asset.FullName,
		ExternalID:      asset.ExternalID,
		Visibility:      asset.Visibility,
		Controllability: asset.Controllability,
		Status:          asset.Status,
		Description:     asset.Description,
		Metadata:        asset.Metadata,
		LastSyncedAt:    utilHTTPTime(asset.LastSyncedAt),
		Notes:           asset.Notes,
		CreatedAt:       utilHTTPTime(asset.CreatedAt),
		UpdatedAt:       utilHTTPTime(asset.UpdatedAt),
	}
}

func ToDependencyAssetDTOs(assets []service.DependencyAsset) []DependencyAssetDTO {
	out := make([]DependencyAssetDTO, 0, len(assets))
	for _, asset := range assets {
		out = append(out, ToDependencyAssetDTO(asset))
	}
	return out
}

func ToPortGroupAssetLinkDTO(link service.PortGroupAssetLink) PortGroupAssetLinkDTO {
	out := PortGroupAssetLinkDTO{
		ID:           link.ID,
		PortGroupID:  link.PortGroupID,
		PortSlotID:   link.PortSlotID,
		AssetID:      link.AssetID,
		RelationType: link.RelationType,
		Required:     link.Required,
		Notes:        link.Notes,
		CreatedAt:    utilHTTPTime(link.CreatedAt),
		UpdatedAt:    utilHTTPTime(link.UpdatedAt),
	}
	if link.Asset != nil {
		asset := ToDependencyAssetDTO(*link.Asset)
		out.Asset = &asset
	}
	return out
}

func ToPortGroupAssetLinkDTOs(links []service.PortGroupAssetLink) []PortGroupAssetLinkDTO {
	out := make([]PortGroupAssetLinkDTO, 0, len(links))
	for _, link := range links {
		out = append(out, ToPortGroupAssetLinkDTO(link))
	}
	return out
}

func ToPortGroupRepositoryLinkDTO(link service.PortGroupRepositoryLink) PortGroupRepositoryLinkDTO {
	out := PortGroupRepositoryLinkDTO{
		ID: link.ID, PortGroupID: link.PortGroupID, PortSlotID: link.PortSlotID,
		RepositoryID: link.RepositoryID, RelationType: link.RelationType, Required: link.Required,
		Notes: link.Notes, CreatedAt: utilHTTPTime(link.CreatedAt), UpdatedAt: utilHTTPTime(link.UpdatedAt),
	}
	if link.Repository != nil {
		repository := ToGithubRepositoryDTO(*link.Repository)
		out.Repository = &repository
	}
	return out
}

func ToPortGroupRepositoryLinkDTOs(links []service.PortGroupRepositoryLink) []PortGroupRepositoryLinkDTO {
	out := make([]PortGroupRepositoryLinkDTO, 0, len(links))
	for _, link := range links {
		out = append(out, ToPortGroupRepositoryLinkDTO(link))
	}
	return out
}

func ToPortGroupDTO(group service.PortGroupView) PortGroupDTO {
	out := PortGroupDTO{
		ID:               group.ID,
		HostID:           group.HostID,
		PortPrefix:       group.PortPrefix,
		EnvironmentName:  group.EnvironmentName,
		EnvironmentOwner: group.EnvironmentOwner,
		RuntimeMode:      group.RuntimeMode,
		RuntimeName:      group.RuntimeName,
		ServiceIP:        group.ServiceIP,
		Status:           group.Status,
		Tags:             group.Tags,
		Notes:            group.Notes,
		CreatedAt:        utilHTTPTime(group.CreatedAt),
		UpdatedAt:        utilHTTPTime(group.UpdatedAt),
		Slots:            ToPortSlotDTOs(group.Slots),
		AssetLinks:       ToPortGroupAssetLinkDTOs(group.AssetLinks),
		RepositoryLinks:  ToPortGroupRepositoryLinkDTOs(group.RepositoryLinks),
	}
	if group.Host != nil {
		host := ToHostDTO(*group.Host)
		out.Host = &host
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

func ToServiceGroupPortGroupDTO(group service.ServiceGroupPortGroup) ServiceGroupPortGroupDTO {
	out := ServiceGroupPortGroupDTO{
		ID:             group.ID,
		ServiceGroupID: group.ServiceGroupID,
		PortGroupID:    group.PortGroupID,
		Role:           group.Role,
		Notes:          group.Notes,
		CreatedAt:      utilHTTPTime(group.CreatedAt),
		UpdatedAt:      utilHTTPTime(group.UpdatedAt),
	}
	if group.PortGroup != nil {
		portGroup := ToPortGroupDTO(service.PortGroupView{PortGroup: *group.PortGroup})
		out.PortGroup = &portGroup
	}
	return out
}

func ToServiceGroupPortGroupDTOs(groups []service.ServiceGroupPortGroup) []ServiceGroupPortGroupDTO {
	out := make([]ServiceGroupPortGroupDTO, 0, len(groups))
	for _, group := range groups {
		out = append(out, ToServiceGroupPortGroupDTO(group))
	}
	return out
}

func ToServiceGroupDTO(group service.ServiceGroupView) ServiceGroupDTO {
	return ServiceGroupDTO{
		ID:          group.ID,
		Name:        group.Name,
		Kind:        group.Kind,
		Status:      group.Status,
		Description: group.Description,
		Notes:       group.Notes,
		CreatedAt:   utilHTTPTime(group.CreatedAt),
		UpdatedAt:   utilHTTPTime(group.UpdatedAt),
		PortGroups:  ToServiceGroupPortGroupDTOs(group.PortGroups),
	}
}

func ToServiceGroupDTOs(groups []service.ServiceGroupView) []ServiceGroupDTO {
	out := make([]ServiceGroupDTO, 0, len(groups))
	for _, group := range groups {
		out = append(out, ToServiceGroupDTO(group))
	}
	return out
}
