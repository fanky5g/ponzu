### Big Picture
We will stick with parent ponzu [reference type implementation](https://docs.ponzu-cms.org/References/Overview/) 
with a few modifications.

1. Database drivers that support relations i.e RDMS are encouraged to implement references in their own way.
As an indicator to reference type fields, these fields will have a struct tag `reference:"ReferenceTypeName"` that 
database drivers could use to implement reference types. Foreign keys and data loading implementation are left to
drivers.
2. A field Refererences will be added to Ponzu Content Types which will be a plain map containing loaded data:
`
[TypeName]: {
    [ID]: {
        "field": "value",
    },
}
`
A convenience method will be added to Ponzu Item type to fetch a referenced type by ID.
Database drivers must load references and populate them in this map.
Note: Reference fields themselves will not hold their values. They will be as implemented in ponzu. ie. Plain string 
URLs pointing to their location. 
3. Ponzu Repository structure will be modified to allow callers to specify if they want to load data with or without
references.
4. HTTP/2 Server push for sending reference types will be ommitted.
Clients will have to check if reference field is populated before initiating a request to fetch type.
