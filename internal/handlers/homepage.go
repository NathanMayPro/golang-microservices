package handlers

import (
	"net/http"

	"golang-microservices/internal/customlog" // import loggers package
	"golang-microservices/internal/datamodel" // for HttpResponse struct
)





func HandlerHomepage(w http.ResponseWriter, r *http.Request) (datamodel.Response, error) {
	// log activity
	customlog.NewLog("endpoint:homepage hit", r)
	return datamodel.Response{
		Status:  200,
		Message: "success",
	}, nil
}