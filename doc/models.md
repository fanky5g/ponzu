- BoltDB is a key-value store (requires no models)
- Postgres is a table store (requiring models)
- We need to pass models when declaring a new database
- Models will be stored in user application dir and user of ponzu has to pass them in application init
- System models for (Config, Uploads, Users, Analytics) have to be generated if it does not exist.
- We mark system models file as generated (non-editable)
- Currently, we do not have system items generated to user application dir. It seems a contradiction if we generate
  content models to user application dir but not system content