## Compound Interest BNPL Domain

Domain-driven design sketch for an overdue/compound-interest BNPL system in Go.

### Package layout
- `domain/payment`: Payment aggregate, lifecycle state (`Status*`), overdue info, domain errors.
- `domain/money`: Decimal-based Money value object with currency scale handling and BPS helpers.
- `domain/user`: User aggregate stub with scoped ID and validation.
- `domain/shared`: Cross-cutting ID helper (ULID).

### Key design points
- Payment encapsulates transitions (`Pay`, `MarkOverdue`) to guard invariants (no double-pay, no overdue after pay).
- Money uses `shopspring/decimal` and currency-specific scale (KRW:0, USD:2) to preserve precision; BPS helpers support interest calculations.
- IDs are generated via ULID and wrapped per aggregate to avoid zero-value leaks.

### Running tests
```bash
go test ./...
```

### Next steps
- Introduce a clock/rate interface and implement compound interest accrual on Payment.
- Add persistence/adapters while keeping domain free of transport/types.
- Extend Money with more currencies and rounding policies as needed.
