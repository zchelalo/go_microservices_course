package course

import (
	"context"
	"errors"

	"github.com/zchelalo/go_microservices_meta/meta"
	"github.com/zchelalo/go_microservices_response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetRequest struct {
		Id string
	}

	GetAllRequest struct {
		Name  string
		Limit int
		Page  int
	}

	UpdateRequest struct {
		Id        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	DeleteRequest struct {
		Id string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(service Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(service),
		Get:    makeGetEndpoint(service),
		GetAll: makeGetAllEndpoint(service, config),
		Update: makeUpdateEndpoint(service),
		Delete: makeDeleteEndpoint(service),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		course, err := service.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			if err == ErrInvalidStartDate || err == ErrInvalidEndDate || err == ErrEndLesserStart {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)

		course, err := service.Get(ctx, req.Id)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil
	}
}

func makeGetAllEndpoint(service Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllRequest)

		filters := Filters{
			Name: req.Name,
		}

		count, err := service.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := service.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", courses, meta), nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		err := service.Update(ctx, req.Id, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			if err == ErrInvalidStartDate || err == ErrInvalidEndDate || err == ErrEndLesserStart {
				return nil, response.BadRequest(err.Error())
			}
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", "course updated successfully", nil), nil
	}
}

func makeDeleteEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		err := service.Delete(ctx, req.Id)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", "course deleted successfully", nil), nil
	}
}
