The idea of field collections is borrowed from pimcore.
A field collection is in theory similar to a content-type or class that can be encapsulated in a Content or Item
without it being particularly a reference type or self-contained item/content.

An example of a field collection can be:
For a content-type Page, there can exist a field contentBlocks which can be a varied list of other field collection
types e.g.
- ImageSliderBlock
- ImageAndTextBlock
- TabGroup
e.t.c.

Field Collections will allow composing more complex types.

Steps To Achieve Field Collections:
- We need to be able to generate plain-types.
  + A plain-type by definition is a type that can be added as a nested field to another type.
  + A plain-type cannot nest itself.
  + Generate command should by default generate plain-types. An extra argument "content-type" must be specified
    to make a plain-type a content-type.
  + plain-types cannot be self persisted. They must be used as fields in content-types.
- When generating a content-type, one can reference a plain-type to be used in any field as a nested field
- Generate command can generate a FieldCollection type. A FieldCollection type must have as properties:
  + Title
  + Accepted Types: Can be a map of plain-type constructors