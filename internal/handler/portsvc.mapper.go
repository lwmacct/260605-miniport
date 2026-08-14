package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToHostCreatePayload(input HostCreateDTO) service.HostPayload {
	return service.HostPayload{Name: input.Name, IP: input.IP, Spec: input.Spec, Status: input.Status, Notes: input.Notes}
}

func ToHostUpdatePayload(input HostUpdateDTO) service.HostUpdateInput {
	return service.HostUpdateInput{ID: input.ID, Payload: service.HostPayload{Name: input.Name, IP: input.IP, Spec: input.Spec, Status: input.Status, Notes: input.Notes}}
}

func ToDependencyAssetCreatePayload(input DependencyAssetCreateDTO) service.DependencyAssetPayload {
	return service.DependencyAssetPayload{
		Name: input.Name, AssetKind: input.AssetKind, AssetType: input.AssetType, Provider: input.Provider,
		URL: input.URL, FullName: input.FullName, ExternalID: input.ExternalID, Visibility: input.Visibility,
		Controllability: input.Controllability, Status: input.Status, Description: input.Description,
		Metadata: input.Metadata, Notes: input.Notes,
	}
}

func ToDependencyAssetUpdatePayload(input DependencyAssetUpdateDTO) service.DependencyAssetUpdateInput {
	return service.DependencyAssetUpdateInput{ID: input.ID, Payload: service.DependencyAssetPayload{
		Name: input.Name, AssetKind: input.AssetKind, AssetType: input.AssetType, Provider: input.Provider,
		URL: input.URL, FullName: input.FullName, ExternalID: input.ExternalID, Visibility: input.Visibility,
		Controllability: input.Controllability, Status: input.Status, Description: input.Description,
		Metadata: input.Metadata, Notes: input.Notes,
	}}
}

func ToPortGroupCreatePayload(input PortGroupCreateDTO) service.PortGroupPayload {
	return service.PortGroupPayload{
		HostID: input.HostID, PortPrefix: input.PortPrefix, EnvironmentName: input.EnvironmentName,
		EnvironmentOwner: input.EnvironmentOwner, RuntimeMode: input.RuntimeMode, RuntimeName: input.RuntimeName,
		ServiceIP: input.ServiceIP, Status: input.Status, Tags: input.Tags, Notes: input.Notes,
		Slots: ToPortSlotPayloads(input.Slots), AssetLinks: ToAssetLinkPayloads(input.AssetLinks),
		RepositoryLinks: ToRepositoryLinkPayloads(input.RepositoryLinks),
	}
}

func ToPortGroupUpdatePayload(input PortGroupUpdateDTO) service.PortGroupUpdateInput {
	return service.PortGroupUpdateInput{ID: input.ID, Payload: service.PortGroupPayload{
		HostID: input.HostID, PortPrefix: input.PortPrefix, EnvironmentName: input.EnvironmentName,
		EnvironmentOwner: input.EnvironmentOwner, RuntimeMode: input.RuntimeMode, RuntimeName: input.RuntimeName,
		ServiceIP: input.ServiceIP, Status: input.Status, Tags: input.Tags, Notes: input.Notes,
		Slots: ToPortSlotPayloads(input.Slots), AssetLinks: ToAssetLinkPayloads(input.AssetLinks),
		RepositoryLinks: ToRepositoryLinkPayloads(input.RepositoryLinks),
	}}
}

func ToPortSlotPayloads(values []PortSlotInputDTO) []service.PortSlotPayload {
	out := make([]service.PortSlotPayload, 0, len(values))
	for _, value := range values {
		out = append(out, service.PortSlotPayload{
			ID: value.ID, Port: value.Port, Name: value.Name, Kind: value.Kind, Protocol: value.Protocol,
			ContainerName: value.ContainerName, Status: value.Status, Notes: value.Notes,
		})
	}
	return out
}

func ToAssetLinkPayloads(values []PortGroupAssetLinkInputDTO) []service.PortGroupAssetLinkPayload {
	out := make([]service.PortGroupAssetLinkPayload, 0, len(values))
	for _, value := range values {
		out = append(out, service.PortGroupAssetLinkPayload{
			ID: value.ID, PortSlotID: value.PortSlotID, AssetID: value.AssetID, RelationType: value.RelationType,
			Required: value.Required, Notes: value.Notes,
		})
	}
	return out
}

func ToRepositoryLinkPayloads(values []PortGroupRepositoryLinkInputDTO) []service.PortGroupRepositoryLinkPayload {
	out := make([]service.PortGroupRepositoryLinkPayload, 0, len(values))
	for _, value := range values {
		out = append(out, service.PortGroupRepositoryLinkPayload{
			ID: value.ID, PortSlotID: value.PortSlotID, RepositoryID: value.RepositoryID,
			RelationType: value.RelationType, Required: value.Required, Notes: value.Notes,
		})
	}
	return out
}

func ToServiceGroupCreatePayload(input ServiceGroupCreateDTO) service.ServiceGroupPayload {
	return service.ServiceGroupPayload{
		Name: input.Name, Kind: input.Kind, Status: input.Status, Description: input.Description,
		Notes: input.Notes, PortGroups: ToServiceGroupPortGroupPayloads(input.PortGroups),
	}
}

func ToServiceGroupUpdatePayload(input ServiceGroupUpdateDTO) service.ServiceGroupUpdateInput {
	return service.ServiceGroupUpdateInput{ID: input.ID, Payload: service.ServiceGroupPayload{
		Name: input.Name, Kind: input.Kind, Status: input.Status, Description: input.Description,
		Notes: input.Notes, PortGroups: ToServiceGroupPortGroupPayloads(input.PortGroups),
	}}
}

func ToServiceGroupPortGroupPayloads(values []ServiceGroupPortGroupInputDTO) []service.ServiceGroupPortGroupPayload {
	out := make([]service.ServiceGroupPortGroupPayload, 0, len(values))
	for _, value := range values {
		out = append(out, service.ServiceGroupPortGroupPayload{ID: value.ID, PortGroupID: value.PortGroupID, Role: value.Role, Notes: value.Notes})
	}
	return out
}

func ToHostDTO(host service.Host) HostDTO {
	return HostDTO{ID: host.ID, Name: host.Name, IP: host.IP, Spec: host.Spec, Status: host.Status, Notes: host.Notes, CreatedAt: utilHTTPTime(host.CreatedAt), UpdatedAt: utilHTTPTime(host.UpdatedAt)}
}

func ToHostDTOs(hosts []service.Host) []HostDTO {
	out := make([]HostDTO, 0, len(hosts))
	for _, host := range hosts {
		out = append(out, ToHostDTO(host))
	}
	return out
}

func ToDependencyAssetDTO(asset service.DependencyAsset) DependencyAssetDTO {
	return DependencyAssetDTO{
		ID: asset.ID, Name: asset.Name, AssetKind: asset.AssetKind, AssetType: asset.AssetType,
		Provider: asset.Provider, URL: asset.URL, FullName: asset.FullName, ExternalID: asset.ExternalID,
		Visibility: asset.Visibility, Controllability: asset.Controllability, Status: asset.Status,
		Description: asset.Description, Metadata: asset.Metadata, LastSyncedAt: utilHTTPTime(asset.LastSyncedAt),
		Notes: asset.Notes, CreatedAt: utilHTTPTime(asset.CreatedAt), UpdatedAt: utilHTTPTime(asset.UpdatedAt),
	}
}

func ToDependencyAssetDTOs(values []service.DependencyAsset) []DependencyAssetDTO {
	out := make([]DependencyAssetDTO, 0, len(values))
	for _, value := range values {
		out = append(out, ToDependencyAssetDTO(value))
	}
	return out
}

func ToPortSlotDTO(slot service.PortSlot) PortSlotDTO {
	return PortSlotDTO{ID: slot.ID, PortGroupID: slot.PortGroupID, Port: slot.Port, Name: slot.Name, Kind: slot.Kind, Protocol: slot.Protocol, ContainerName: slot.ContainerName, Status: slot.Status, Notes: slot.Notes, CreatedAt: utilHTTPTime(slot.CreatedAt), UpdatedAt: utilHTTPTime(slot.UpdatedAt)}
}

func ToPortGroupAssetLinkDTO(link service.PortGroupAssetLink) PortGroupAssetLinkDTO {
	out := PortGroupAssetLinkDTO{ID: link.ID, PortGroupID: link.PortGroupID, PortSlotID: link.PortSlotID, AssetID: link.AssetID, RelationType: link.RelationType, Required: link.Required, Notes: link.Notes, CreatedAt: utilHTTPTime(link.CreatedAt), UpdatedAt: utilHTTPTime(link.UpdatedAt)}
	if link.Asset != nil {
		asset := ToDependencyAssetDTO(*link.Asset)
		out.Asset = &asset
	}
	return out
}

func ToPortGroupRepositoryLinkDTO(link service.PortGroupRepositoryLink) PortGroupRepositoryLinkDTO {
	out := PortGroupRepositoryLinkDTO{ID: link.ID, PortGroupID: link.PortGroupID, PortSlotID: link.PortSlotID, RepositoryID: link.RepositoryID, RelationType: link.RelationType, Required: link.Required, Notes: link.Notes, CreatedAt: utilHTTPTime(link.CreatedAt), UpdatedAt: utilHTTPTime(link.UpdatedAt)}
	if link.Repository != nil {
		repository := ToGithubRepositoryDTO(*link.Repository)
		out.Repository = &repository
	}
	return out
}

func ToPortGroupDTO(group service.PortGroupView) PortGroupDTO {
	out := PortGroupDTO{
		ID: group.ID, HostID: group.HostID, PortPrefix: group.PortPrefix, EnvironmentName: group.EnvironmentName,
		EnvironmentOwner: group.EnvironmentOwner, RuntimeMode: group.RuntimeMode, RuntimeName: group.RuntimeName,
		ServiceIP: group.ServiceIP, Status: group.Status, Tags: group.Tags, Notes: group.Notes,
		CreatedAt: utilHTTPTime(group.CreatedAt), UpdatedAt: utilHTTPTime(group.UpdatedAt),
		Slots: make([]PortSlotDTO, 0, len(group.Slots)), AssetLinks: make([]PortGroupAssetLinkDTO, 0, len(group.AssetLinks)),
		RepositoryLinks: make([]PortGroupRepositoryLinkDTO, 0, len(group.RepositoryLinks)),
	}
	if group.Host != nil {
		host := ToHostDTO(*group.Host)
		out.Host = &host
	}
	for _, slot := range group.Slots {
		out.Slots = append(out.Slots, ToPortSlotDTO(slot))
	}
	for _, link := range group.AssetLinks {
		out.AssetLinks = append(out.AssetLinks, ToPortGroupAssetLinkDTO(link))
	}
	for _, link := range group.RepositoryLinks {
		out.RepositoryLinks = append(out.RepositoryLinks, ToPortGroupRepositoryLinkDTO(link))
	}
	return out
}

func ToPortGroupDTOs(values []service.PortGroupView) []PortGroupDTO {
	out := make([]PortGroupDTO, 0, len(values))
	for _, value := range values {
		out = append(out, ToPortGroupDTO(value))
	}
	return out
}

func ToServiceGroupPortGroupDTO(value service.ServiceGroupPortGroup) ServiceGroupPortGroupDTO {
	out := ServiceGroupPortGroupDTO{ID: value.ID, ServiceGroupID: value.ServiceGroupID, PortGroupID: value.PortGroupID, Role: value.Role, Notes: value.Notes, CreatedAt: utilHTTPTime(value.CreatedAt), UpdatedAt: utilHTTPTime(value.UpdatedAt)}
	if value.PortGroup != nil {
		portGroup := ToPortGroupDTO(service.PortGroupView{PortGroup: *value.PortGroup})
		out.PortGroup = &portGroup
	}
	return out
}

func ToServiceGroupDTO(value service.ServiceGroupView) ServiceGroupDTO {
	out := ServiceGroupDTO{ID: value.ID, Name: value.Name, Kind: value.Kind, Status: value.Status, Description: value.Description, Notes: value.Notes, CreatedAt: utilHTTPTime(value.CreatedAt), UpdatedAt: utilHTTPTime(value.UpdatedAt), PortGroups: make([]ServiceGroupPortGroupDTO, 0, len(value.PortGroups))}
	for _, portGroup := range value.PortGroups {
		out.PortGroups = append(out.PortGroups, ToServiceGroupPortGroupDTO(portGroup))
	}
	return out
}

func ToServiceGroupDTOs(values []service.ServiceGroupView) []ServiceGroupDTO {
	out := make([]ServiceGroupDTO, 0, len(values))
	for _, value := range values {
		out = append(out, ToServiceGroupDTO(value))
	}
	return out
}
