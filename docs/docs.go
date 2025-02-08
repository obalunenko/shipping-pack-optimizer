// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Oleg Balunenko",
            "email": "oleg.balunenko@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/license/mit"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/pack": {
            "post": {
                "description": "Calculates the number of packs needed to ship to a customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pack"
                ],
                "summary": "Get the number of packs needed to ship to a customer",
                "operationId": "shipping-pack-optimizer-pack\tpost",
                "parameters": [
                    {
                        "description": "Request data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.PackRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with packs data",
                        "schema": {
                            "$ref": "#/definitions/service.PackResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request data",
                        "schema": {
                            "$ref": "#/definitions/service.badRequestError"
                        }
                    },
                    "405": {
                        "description": "Method not allowed",
                        "schema": {
                            "$ref": "#/definitions/service.methodNotAllowedError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/service.internalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "service.Pack": {
            "type": "object",
            "properties": {
                "box": {
                    "type": "integer",
                    "format": "uint",
                    "example": 50
                },
                "quantity": {
                    "type": "integer",
                    "format": "uint",
                    "example": 3
                }
            }
        },
        "service.PackRequest": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "integer",
                    "format": "uint",
                    "example": 543
                }
            }
        },
        "service.PackResponse": {
            "type": "object",
            "properties": {
                "packs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.Pack"
                    }
                }
            }
        },
        "service.badRequestError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "Bad request"
                }
            }
        },
        "service.internalServerError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 500
                },
                "message": {
                    "type": "string",
                    "example": "Internal server error"
                }
            }
        },
        "service.methodNotAllowedError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 405
                },
                "message": {
                    "type": "string",
                    "example": "Method not allowed"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{"http"},
	Title:            "Shipping Pack Optimizer API",
	Description:      "This is a simple API for packing orders",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
