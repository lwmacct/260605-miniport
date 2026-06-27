package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToInventoryActor(session *AuthSessionDTO) service.InventoryActor {
	actor := service.InventoryActor{}
	if session == nil || session.User == nil {
		return actor
	}
	actor.UserID = session.User.ID
	actor.Username = session.User.Username
	actor.Admin = session.User.Admin
	return actor
}

func ToPortGroupPayload(input AllocationPayloadDTO) service.PortGroupPayload {
	out := service.PortGroupPayload{
		UserID:        input.UserID,
		PortStart:     input.PortStart,
		PortEnd:       input.PortEnd,
		Name:          input.Name,
		DindIP:        input.DindIP,
		DindContainer: input.DindContainer,
		Status:        input.Status,
		Owner:         input.Owner,
		Tags:          input.Tags,
		Notes:         input.Notes,
		Slots:         make([]service.PortSlotPayload, 0, len(input.Slots)),
		Projects:      make([]service.ProjectPayload, 0, len(input.Projects)),
		Components:    make([]service.ComponentPayload, 0, len(input.Components)),
		Repositories:  make([]service.RepositoryPayload, 0, len(input.Repositories)),
	}
	for _, slot := range input.Slots {
		out.Slots = append(out.Slots, service.PortSlotPayload(slot))
	}
	for _, project := range input.Projects {
		out.Projects = append(out.Projects, service.ProjectPayload(project))
	}
	for _, component := range input.Components {
		out.Components = append(out.Components, service.ComponentPayload(component))
	}
	for _, repository := range input.Repositories {
		out.Repositories = append(out.Repositories, service.RepositoryPayload(repository))
	}
	return out
}

func ToBatchUpdateInput(actor service.InventoryActor, input AllocationBatchUpdateDTO) service.PortGroupBatchUpdateInput {
	return service.PortGroupBatchUpdateInput{
		Actor:  actor,
		IDs:    input.IDs,
		Owner:  input.Owner,
		Status: input.Status,
		Tags:   input.Tags,
	}
}

func ToBatchDeleteInput(actor service.InventoryActor, input AllocationBatchDeleteDTO) service.PortGroupBatchDeleteInput {
	return service.PortGroupBatchDeleteInput{
		Actor: actor,
		IDs:   input.IDs,
	}
}

func ToPortSlotDTO(slot service.InventoryPortSlot) PortSlotDTO {
	return PortSlotDTO{
		ID:           slot.ID,
		AllocationID: slot.PortGroupID,
		Port:         slot.Port,
		Name:         slot.Name,
		Protocol:     slot.Protocol,
		Purpose:      slot.Purpose,
		Status:       slot.Status,
		Notes:        slot.Notes,
		CreatedAt:    utilHTTPTime(slot.CreatedAt),
		UpdatedAt:    utilHTTPTime(slot.UpdatedAt),
	}
}

func ToProjectDTO(project service.InventoryProject) ProjectDTO {
	return ProjectDTO{
		ID:           project.ID,
		AllocationID: project.PortGroupID,
		Name:         project.Name,
		Description:  project.Description,
		Notes:        project.Notes,
		CreatedAt:    utilHTTPTime(project.CreatedAt),
		UpdatedAt:    utilHTTPTime(project.UpdatedAt),
	}
}

func ToComponentDTO(component service.InventoryComponent) ComponentDTO {
	return ComponentDTO{
		ID:           component.ID,
		AllocationID: component.PortGroupID,
		Name:         component.Name,
		Type:         component.Type,
		URL:          component.URL,
		Version:      component.Version,
		Notes:        component.Notes,
		CreatedAt:    utilHTTPTime(component.CreatedAt),
		UpdatedAt:    utilHTTPTime(component.UpdatedAt),
	}
}

func ToRepositoryDTO(repository service.InventoryRepositoryRef) RepositoryDTO {
	return RepositoryDTO{
		ID:           repository.ID,
		AllocationID: repository.PortGroupID,
		ProjectID:    repository.ProjectID,
		Name:         repository.Name,
		URL:          repository.URL,
		Kind:         repository.Kind,
		Notes:        repository.Notes,
		CreatedAt:    utilHTTPTime(repository.CreatedAt),
		UpdatedAt:    utilHTTPTime(repository.UpdatedAt),
	}
}

func ToAllocationDTO(group service.PortGroupView) AllocationDTO {
	out := AllocationDTO{
		ID:            group.ID,
		UserID:        group.UserID,
		Username:      group.Username,
		PortStart:     group.PortStart,
		PortEnd:       group.PortEnd,
		Name:          group.Name,
		DindIP:        group.DindIP,
		DindContainer: group.DindContainer,
		Status:        group.Status,
		Owner:         group.Owner,
		Tags:          group.Tags,
		Notes:         group.Notes,
		CreatedAt:     utilHTTPTime(group.CreatedAt),
		UpdatedAt:     utilHTTPTime(group.UpdatedAt),
		Slots:         make([]PortSlotDTO, 0, len(group.Slots)),
		Projects:      make([]ProjectDTO, 0, len(group.Projects)),
		Components:    make([]ComponentDTO, 0, len(group.Components)),
		Repositories:  make([]RepositoryDTO, 0, len(group.Repositories)),
	}
	for _, slot := range group.Slots {
		out.Slots = append(out.Slots, ToPortSlotDTO(slot))
	}
	for _, project := range group.Projects {
		out.Projects = append(out.Projects, ToProjectDTO(project))
	}
	for _, component := range group.Components {
		out.Components = append(out.Components, ToComponentDTO(component))
	}
	for _, repository := range group.Repositories {
		out.Repositories = append(out.Repositories, ToRepositoryDTO(repository))
	}
	return out
}

func ToAllocationDTOs(groups []service.PortGroupView) []AllocationDTO {
	out := make([]AllocationDTO, 0, len(groups))
	for _, group := range groups {
		out = append(out, ToAllocationDTO(group))
	}
	return out
}
