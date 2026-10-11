package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //Request method
		origin := c.Request.Header.Get("Origin") //Request headers
		var headerKeys []string                  // Declare request header keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // Allow access from all origins
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //All cross-origin request methods supported by the server; this avoids repeated preflight requests by the browser
			//  Header type
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              Allow cross-origin settings                                                                                                      Other fields can be returned
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // Key CORS setting that allows the browser to parse the response
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // Cache request information, in seconds
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  Whether cross-origin requests need cookies; defaults to true
			c.Set("content-type", "application/json")                                                                                                                                                              // Set the response format to JSON
		}
		//Allow all OPTIONS methods
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// Handle the request
		c.Next() //  Handle the request
		_ = headerStr
	}
}
