package constants

const RouteTagIdentifier = "route_tag"

type RouteTag string

var (
	AdminRoute RouteTag = "admin"
	APIRoute   RouteTag = "api"
)
