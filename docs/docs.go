// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Ankur Kumar",
            "url": "https://github.com/sirjager"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Welcome",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Welcome",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WelcomeResponse"
                        }
                    }
                }
            }
        },
        "/auth/delete": {
            "get": {
                "description": "Delete User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Delete",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code if already have",
                        "name": "code",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/auth/refresh": {
            "get": {
                "description": "Refreshes Access Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Refresh",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "If true, returns User in body",
                        "name": "user",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "If true, returns AccessToken and SessionID in body",
                        "name": "cookies",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "RefreshTokenResponse",
                        "schema": {
                            "$ref": "#/definitions/RefreshTokenResponse"
                        }
                    }
                }
            }
        },
        "/auth/reset": {
            "post": {
                "description": "Reset Password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Reset",
                "parameters": [
                    {
                        "description": "ResetPasswordParams",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ResetPasswordParams"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/signin": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Signin using credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Signin",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "If true, returns User in body",
                        "name": "user",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "If true, returns AccessToken, RefreshToken and SessionID in body",
                        "name": "cookies",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SignInResponse",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    }
                }
            }
        },
        "/auth/signout/{provider}": {
            "get": {
                "description": "Signout from a provider",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Signout",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provider Name",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Signup using email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Signup",
                "parameters": [
                    {
                        "description": "Signup request params",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SignUpRequestParams"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User object",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            }
        },
        "/auth/user": {
            "get": {
                "description": "Returns Authenticated User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User",
                "responses": {
                    "200": {
                        "description": "UserResponse",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "get": {
                "description": "Email Verification",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Verify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email to verify",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email verification code if already have any",
                        "name": "code",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/{provider}": {
            "get": {
                "description": "Authenticates a user with a specified oauth provider",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "OAuth",
                "parameters": [
                    {
                        "enum": [
                            "google",
                            "github"
                        ],
                        "type": "string",
                        "description": "OAuth provider name [google,github]",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User object",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Health Check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/HealthResponse"
                        }
                    }
                }
            }
        },
        "/swagger": {
            "get": {
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Documentation",
                "responses": {
                    "200": {
                        "description": "HTML content",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Fetch multiple users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resources"
                ],
                "summary": "Multiple Users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number: Default 1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Per Page: Default 100",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "UsersResponse",
                        "schema": {
                            "$ref": "#/definitions/UsersResponse"
                        }
                    }
                }
            }
        },
        "/users/{identity}": {
            "get": {
                "description": "Fetch specific user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resources"
                ],
                "summary": "Single User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identity can either be email or id",
                        "name": "identity",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "UserResponse",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Partially Update User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Resources"
                ],
                "summary": "Update User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identity can either be email or id",
                        "name": "identity",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update User Params",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateUserParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "UserResponse",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "EmailVerificationResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "HealthResponse": {
            "type": "object",
            "properties": {
                "server": {
                    "type": "string"
                },
                "service": {
                    "type": "string"
                },
                "started": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "uptime": {
                    "type": "string"
                }
            }
        },
        "RefreshTokenResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "sessionID": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/User"
                }
            }
        },
        "ResetPasswordParams": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "newPassword": {
                    "type": "string"
                }
            }
        },
        "SignInResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "sessionID": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/User"
                }
            }
        },
        "SignUpRequestParams": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "UpdateUserParams": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "pictureURL": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "blocked": {
                    "type": "boolean"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "pictureURL": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "verified": {
                    "type": "boolean"
                }
            }
        },
        "UserResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/User"
                }
            }
        },
        "UsersResponse": {
            "type": "object",
            "properties": {
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/User"
                    }
                }
            }
        },
        "WelcomeResponse": {
            "type": "object",
            "properties": {
                "docs": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "OAuthAPI",
	Description:      "OAuth API for 3rd party authentication",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
