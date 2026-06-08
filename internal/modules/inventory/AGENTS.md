# Inventory Module Guidelines

## Scope

This module owns host inventory, port groups, port slots, components, and repositories.

## Invariants

- A host IP is required and must remain unique.
- A port group must belong to an existing host.
- A port group must contain exactly 10 ports.
- Port group ranges on the same host must not overlap.
- Slot ports must stay inside their parent port group range and be unique within the group.
- Repository name and URL are required together.

## File Roles

- `schema.go`: table and index creation.
- `model.go`: Bun models.
- `dto.go`: request, response, and Huma input/output types.
- `handler.go`: Huma operation registration and request adaptation.
- `service.go`: business orchestration and transaction workflows.
- `repository.go`: Bun queries, persistence, and view assembly.
- `validation.go`: payload normalization, validation, and child model construction.
- `errors.go`: module-specific sentinel errors when needed.

## Testing

Run:

```bash
go test ./internal/modules/inventory
```

Add tests when changing port range validation, child replacement, or uniqueness behavior.
