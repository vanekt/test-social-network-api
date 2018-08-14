package error

import "github.com/gin-gonic/gin"

var HttpResponseErrorBadRequest = gin.H{"error": "Bad request"}
var HttpResponseErrorNotFound = gin.H{"error": "Not found"}
var HttpResponseErrorInternalServerError = gin.H{"error": "Internal Server Error"}
