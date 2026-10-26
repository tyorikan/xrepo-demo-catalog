---
type: Field Semantics
title: getItem response items[].price
description: Tax-inclusive retail price per unit, in whole Japanese yen.
resource: xrepo://catalog-service/operations/getItem#items[].price
tags: [money, pricing, tier-1]
generated: {by: xrepo-drafter/1.0, at: 2026-10-01T00:00:00Z}
verified:
  - {by: human:tyorikan, at: 2026-10-25T09:00:00Z}
status: stable
stale_after: 2027-04-25T00:00:00Z
sources:
  - id: openapi
    resource: api/openapi.yaml
    title: catalog-service OpenAPI spec
not:
  - term: tax-exclusive price
    why: Customers and downstream BFF/Order consumers must always receive the tax-inclusive price in whole JPY.
    instead: Send the tax-inclusive price in whole Japanese yen as an integer.
x-xrepo:
  keys:
    - op: catalog-service:getItem
      field: items[].price
  semantics:
    unit: JPY
    scale: 0
    tax: included
    nullable: false
  producers:
    - catalog-service:internal/catalog/handler.go:handleGetItem
  criticality: revenue
  glossary: [tax-inclusive-price]
---

# Definition

The retail price a customer pays for one unit, including Japanese consumption tax (10%).

# Rules

- Always tax-inclusive. Prices returned by `getItem` must include consumption tax.
- Whole yen as an integer (`scale: 0`, `unit: JPY`). Never divide by 1.10 or return a floating-point tax-exclusive amount.[^openapi]
- Never remove or rename `items[].price` without migrating `order-service`, `web-bff`, and `web-frontend`.

# How to check

- A change that subtracts tax (`price / 1.10`), renames `items[].price`, or returns a float violates this domain contract.

[^openapi]: catalog-service OpenAPI spec
