package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rawdaGastan/urls_checker/pkg/swagger/server/models"
	"github.com/rawdaGastan/urls_checker/pkg/swagger/server/restapi"

	"github.com/rawdaGastan/urls_checker/pkg/swagger/server/restapi/operations"

	"github.com/rawdaGastan/urls_checker/pkg"
)

func main() {

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewUrlsCheckerAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	server.Port = 8080

	api.GetCheckWebsiteHandler = operations.GetCheckWebsiteHandlerFunc(CheckWebsite)

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

//CheckWebsite route returns status of the urls of a website
func CheckWebsite(params operations.GetCheckWebsiteParams) middleware.Responder {

	linksStatus, err := pkg.Check(params.Website)
	if err != nil {
		fmt.Printf("checking links of %v failed with error: %v\n", params.Website, err)
	}

	log.Fatalln(reflect.TypeOf(linksStatus), []*models.URLStatus{})
	//linksStatus = linksStatus.([]*models.URLStatus)
	return operations.NewGetCheckWebsiteOK().WithPayload(linksStatus)
}
