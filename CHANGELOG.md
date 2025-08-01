## v0.14.0
- Reimplement `InputRepeater` using patterns borrowed from `NestedRepeater`

## v0.13.1
- Fix Nested Repeater injecting child in the wrong place with ':scope' selector without polyfill

## v0.13.0
- Reference loader, load uncached objects


## v0.12.2
- Fix: `GCS` write object content-type


## v0.12.1
- Fix: Upload service GetUploads usage of slice with capacity instead of array


## v0.12.0
- Support new config field "WorkflowStateChangeHandler" of type "workflow.StateChangeTrigger" called when entity's workflow state changes


## v0.11.1
- Remove delete button in content editor view during content creation


## v0.11.0
- Add helper text to nested repeater when empty


## v0.10.1
- Support Nesting NestedRepeaters (Only tested to one level deep).


## v0.10.0
- DataBase config. Support ponzu.props from ${USER_HOME}/.config/ponzu
- Use more dynamic connection string support. Set postgres ssl_mode to disable by default.


## v0.9.1
- Remove .config support in cwd
- Add loading ponzu.props from ${USER_HOME}/.config/ponzu


## v0.9.0
- Support loading ponzu config file from .config directory from working directory


## v0.8.4
- Fix: remove cache control usage in CORS middleware


## v0.8.3
- Fix: missing save button in configuration page
- Fix: setup not saving app name
- Fix: empty public path handling
- Fix: GCS storage handling of opening files


### v0.8.2
- Fix blank public path behavior.


### v.0.8.1
- Localize ponzu-driver-gcs
- Fix broken gcs file URLs


### v0.8.0
- Support nested repeaters as direct children in Field Collections.


### v0.7.3
- Bundle assets and templates with binary


### v0.7.2
- BUGFIX: fix field collections usage of static method receiver variable


### v0.7.1
- BUGFIX: repeated nested type


### v0.7.0
- Support generating select with initial options


### v0.6.0
- Localize ponzu-driver-postgres
- Implement initial concept of workflow transitions


### v0.5.0
- Remove dependency on material-css and replace with [m2-material](https://m2.material.io/)


### v0.4.0
- Update ponzu-driver-gcs to v1.1.0


### v0.3.0
- Select between `local` or `gcs` upload storage driver


### v0.2.1
- Fix storage driver, add searchClient


### v0.2.0
- Update ponzu-driver-postgres to v1.0.0


### v0.1.0
- Update ponzu-driver-bleve to v1.0.3
- Redefine SearchInterface
- Remove SearchClientInterface


### v0.0.0
- Initial refactoring & concept of drivers