package handler

import "github.com/lwmacct/260605-miniport/internal/service"

func ToAdminUserDTO(user service.AdminUser) AdminUserDTO {
	var disabledAt *string
	if user.DisabledAt != nil {
		value := utilHTTPTime(*user.DisabledAt)
		disabledAt = &value
	}
	return AdminUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		Admin:       user.Admin,
		DisabledAt:  disabledAt,
	}
}
