---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clickhouseops_mergetree Resource - clickhouseops"
subcategory: ""
description: |-
  Clickhouse MergeTree Table
---

# clickhouseops_mergetree (Resource)

Clickhouse MergeTree Table



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `columns` (Attributes List) Clickhouse Table Column List (see [below for nested schema](#nestedatt--columns))
- `database_name` (String) Clickhouse Database Name
- `name` (String) Clickhouse Table Name
- `order_by` (List of String) Clickhous columns list for order by

### Optional

- `cluster_name` (String) Clickhouse Cluster Name
- `is_replicated` (Boolean) Clickhouse replicated ReplacingMergeTree
- `partition_by` (String) Clickhouse Cluster Name
- `primary_key` (String) Clickhouse Cluster Name
- `settings` (Attributes List) MergeTree optional settings (see [below for nested schema](#nestedatt--settings))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--columns"></a>
### Nested Schema for `columns`

Required:

- `name` (String) Clickhouse Table Column name
- `type` (String) Clickhouse Table Column type


<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Required:

- `name` (String) Clickhouse table setting name
- `value` (String) Clickhouse table setting value
