package references

import (
	"errors"
	"net/http"

	"github.com/fanky5g/ponzu/exceptions"
	"github.com/fanky5g/ponzu/internal/http/response"
	log "github.com/sirupsen/logrus"
)

var ErrInternalServerError = errors.New("internal server error")

func NewListReferencesHandler(service *Service, mapper *Mapper) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		listReferencesInput, err := mapper.MapReqToListReferencesInputResource(req)
		if err != nil {
			response.Respond(res, req, response.NewJSONResponse(http.StatusBadRequest, nil, err))
		}

		matchedReferences, count, err := service.ListReferences(listReferencesInput.Type, listReferencesInput.Search)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("failed to list references")

			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusInternalServerError, nil, ErrInternalServerError),
			)
			return
		}

		outputResource, err := mapper.MapReferenceMatchesToReferencesResourceOutput(matchedReferences, count)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("failed to map references to response")

			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusInternalServerError, nil, ErrInternalServerError),
			)
			return
		}

		response.Respond(res, req, response.NewJSONResponse(http.StatusOK, outputResource, nil))
	}
}

func NewGetReferenceHandler(service *Service, mapper *Mapper) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		inputResource, err := mapper.MapRequestToGetReferenceInputResource(req)
		if err != nil {
			exceptions.Log(err)
			response.Respond(res, req, response.NewJSONResponse(http.StatusBadRequest, nil, err))
		}

		entity, err := service.GetReference(inputResource.Type, inputResource.ID)
		if err != nil {
			exceptions.Log(err)
			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusInternalServerError, nil, ErrInternalServerError),
			)
			return
		}

		outputResource, err := mapper.MapEntityToReference(entity)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("failed to map entity to response")

			response.Respond(
				res,
				req,
				response.NewJSONResponse(http.StatusInternalServerError, nil, ErrInternalServerError),
			)
			return
		}

		response.Respond(res, req, response.NewJSONResponse(http.StatusOK, outputResource, nil))
	}
}
