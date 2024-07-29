package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/zchelalo/go_microservices_course/internal/course"
	"github.com/zchelalo/go_microservices_response/response"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {
	router := http.NewServeMux()

	opts := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("POST /courses", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	))
	router.Handle("GET /courses", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourse,
		encodeResponse,
		opts...,
	))
	router.Handle("GET /courses/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	))
	router.Handle("PATCH /courses/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	))
	router.Handle("DELETE /courses/{id}", httpTransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	))

	return router
}

func decodeCreateCourse(_ context.Context, request *http.Request) (interface{}, error) {
	var req course.CreateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetCourse(_ context.Context, request *http.Request) (interface{}, error) {
	id := request.PathValue("id")
	req := course.GetRequest{
		Id: id,
	}

	return req, nil
}

func decodeGetAllCourse(_ context.Context, request *http.Request) (interface{}, error) {
	queries := request.URL.Query()

	limit, _ := strconv.Atoi(queries.Get("limit"))
	page, _ := strconv.Atoi(queries.Get("page"))

	req := course.GetAllRequest{
		Name:  queries.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func decodeUpdateCourse(_ context.Context, request *http.Request) (interface{}, error) {
	var req course.UpdateRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	id := request.PathValue("id")
	req.Id = id

	return req, nil
}

func decodeDeleteCourse(_ context.Context, request *http.Request) (interface{}, error) {
	id := request.PathValue("id")
	req := course.DeleteRequest{
		Id: id,
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
