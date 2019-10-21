package router

import (

	// Native Go Libs
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	http "net/http"
	models "oniondrop/models"
	middlewares "oniondrop/router/middlewares"
	"reflect"
	"runtime"
	"strings"

	// 3rd Party Libs
	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

type (
	// Handler : Custom type to work with CustomHandle wrapper
	Handler func(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error)
)

type Greeter struct {
	Message string
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(20 << 30)
		// logic part of log in
		fmt.Println("username:", r.Form["user_name"])
		fmt.Println("password:", r.Form["password1"])
		fmt.Println("email:", r.Form["email"])
		fmt.Println("File Upload Endpoint Hit")
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, handler, err := r.FormFile("fileupload")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("tmp", "name-*.zip")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		fmt.Fprintf(w, "Successfully Uploaded File\n")

		models.ContainerRun()
	}
}

// CustomHandle : Custom Handlers Wrapper for API
func CustomHandle(env *models.Env, handlers ...Handler) http.Handler {

	statusCode := ""

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseDetails := &customhttpresponse.ResponseDetails{}

		// Retrieve AuthMiddleware method name for response details
		action := strings.Split(runtime.FuncForPC(reflect.ValueOf(middlewares.AuthMiddleware).Pointer()).Name(), ".")[1]

		// Get UserID through authentication middleware
		userID, err := middlewares.AuthMiddleware(env, w, r)

		if err != nil {
			responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, customhttpresponse.CodeInvalidToken)
			customhttpresponse.WriteResponse(nil, responseDetails, w)
			return
		}

		// Pass UserID to request context
		ctx := context.WithValue(r.Context(), middlewares.ContextUserKey, userID)

		if err != nil {

			if statusCode == customhttpresponse.CodeValidationFailed {
				responseDetails = customhttpresponse.NewResponseDetailsWithFields(strings.Split(err.Error(), "|"), env.Config.Service, runtime.FuncForPC(reflect.ValueOf(middlewares.AuthMiddleware).Pointer()).Name(), statusCode)
			} else {
				// FIXME: Remove Debugging Mode before production
				responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, statusCode)
			}

			customhttpresponse.WriteResponse(nil, responseDetails, w)

			// We can then log error somewhere here

			return
		}

		// If auth check is successful, trigger handlers
		for _, h := range handlers {

			// Retrieve handler method name for response details
			action = strings.Split(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), ".")[1]
			statusCode, err = h(env, w, r.WithContext(ctx))

			if err != nil {

				if statusCode == customhttpresponse.CodeValidationFailed {
					responseDetails = customhttpresponse.NewResponseDetailsWithFields(strings.Split(err.Error(), "|"), env.Config.Service, runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), statusCode)
				} else {
					// FIXME: Remove Debugging Mode before production
					responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, statusCode)
				}

				customhttpresponse.WriteResponse(nil, responseDetails, w)

				// We can then log error somewhere here

				return
			}
		}

		// We can then log success somewhere here
	})
}
