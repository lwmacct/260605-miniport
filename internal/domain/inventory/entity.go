package inventory

import (
	"time"

	"github.com/uptrace/bun"
)

type Host struct {
	bun.BaseModel `bun:"table:hosts"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	IP          string    `bun:"ip,notnull,unique" json:"ip"`
	Name        string    `bun:"name" json:"name"`
	Network     string    `bun:"network" json:"network"`
	Environment string    `bun:"environment" json:"environment"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type PortGroup struct {
	bun.BaseModel `bun:"table:port_groups"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	HostID        int64     `bun:"host_id,notnull" json:"hostId"`
	Host          *Host     `bun:"rel:belongs-to,join:host_id=id" json:"host,omitempty"`
	PortStart     int       `bun:"port_start,notnull" json:"portStart"`
	PortEnd       int       `bun:"port_end,notnull" json:"portEnd"`
	ServiceName   string    `bun:"service_name,notnull" json:"serviceName"`
	ContainerName string    `bun:"container_name" json:"containerName"`
	DindHost      string    `bun:"dind_host" json:"dindHost"`
	Status        string    `bun:"status,notnull" json:"status"`
	Owner         string    `bun:"owner" json:"owner"`
	Tags          string    `bun:"tags" json:"tags"`
	Notes         string    `bun:"notes" json:"notes"`
	CreatedAt     time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt     time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type PortSlot struct {
	bun.BaseModel `bun:"table:port_slots"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
	Port        int       `bun:"port,notnull" json:"port"`
	Name        string    `bun:"name" json:"name"`
	Protocol    string    `bun:"protocol,notnull" json:"protocol"`
	Purpose     string    `bun:"purpose" json:"purpose"`
	Status      string    `bun:"status,notnull" json:"status"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type Component struct {
	bun.BaseModel `bun:"table:components"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
	Name        string    `bun:"name,notnull" json:"name"`
	Type        string    `bun:"type,notnull" json:"type"`
	URL         string    `bun:"url" json:"url"`
	Version     string    `bun:"version" json:"version"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}

type Repository struct {
	bun.BaseModel `bun:"table:repositories"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	PortGroupID int64     `bun:"port_group_id,notnull" json:"portGroupId"`
	Name        string    `bun:"name,notnull" json:"name"`
	URL         string    `bun:"url,notnull" json:"url"`
	Kind        string    `bun:"kind,notnull" json:"kind"`
	Notes       string    `bun:"notes" json:"notes"`
	CreatedAt   time.Time `bun:"created_at,notnull" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull" json:"updatedAt"`
}
