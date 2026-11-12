package docs

import (
	"os"
	"strings"
	"github.com/swaggo/swag"
)

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "paths": {
        "/health": {
            "get": {
                "description": "Check if the service is alive",
                "produces": ["text/plain"],
                "tags": ["System"],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func init() {
	host := getEnv("SWAGGER_HOST", "localhost:80")
	version := getEnv("API_VERSION", "1.0")
	title := getEnv("API_TITLE", "API Gateway")
	description := getEnv("API_DESCRIPTION", "API Gateway for Movie Recommendation Platform")
	basePath := getEnv("API_BASE_PATH", "/")
	schemesStr := getEnv("SWAGGER_SCHEMES", "http")

	schemes := strings.Split(schemesStr, ",")
	for i := range schemes {
		schemes[i] = strings.TrimSpace(schemes[i])
	}

	SwaggerInfo = &swag.Spec{
		Version:          version,
		Host:             host,
		BasePath:         basePath,
		Schemes:          schemes,
		Title:            title,
		Description:      description,
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
		LeftDelim:        "{{",
		RightDelim:       "}}",
	}

	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

var SwaggerInfo *swag.Spec
