# S022 Runtime Evidence Model

## Family Outcome

Each fixed family descriptor has an A/AAAA label, expected IP family 4/6, and numeric DoH type 1/28. A successful query produces valid family-specific addresses. A valid empty response supplies no evidence; a rejected query supplies an operational or response error retaining available code/message.

## Resolver Outcome

Both families settle independently. Any usable family yields a unique lexically sorted address array. Without usable evidence, an error names the hostname and resolver and summarizes both families in fixed A then AAAA order. System and DoH results cannot compensate for each other.

## Validation

System responses must be arrays of valid requested-family IP strings. DoH requires an object with numeric successful Status and an omitted or array Answer collection. Only requested-type records with string data validating as the requested family enter the inventory. Aliases may coexist with valid target addresses but cannot establish success alone.

No new serialized product fields, persistent schema, or source content is created.
