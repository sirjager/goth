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
        "/api": {
            "get": {
                "description": "Welcome message",
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
        "/api/admin/user/{identity}": {
            "patch": {
                "description": "Partially Update User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Update User",
                "parameters": [
                    {
                        "description": "Update User Params",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateUserParams"
                        }
                    },
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
            }
        },
        "/api/admin/users": {
            "get": {
                "description": "Fetch multiple users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Fetch Users",
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
        "/api/auth/refresh": {
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
                "summary": "Refresh Token",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "If true, returns User in body",
                        "name": "user",
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
        "/api/auth/reset": {
            "post": {
                "description": "Reset password with a verified email email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Reset Password",
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
        "/api/auth/signin": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "SignIn using credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignIn User",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "If true, returns User in body",
                        "name": "user",
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
        "/api/auth/signout/{provider}": {
            "get": {
                "description": "Signout session(s) or a provider",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignOut User",
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
        "/api/auth/signup": {
            "post": {
                "description": "Sign up a new user using email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "SignUp User",
                "parameters": [
                    {
                        "description": "sign up params : email and password",
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
        "/api/auth/user": {
            "get": {
                "description": "Get Authenticated User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Fetch",
                "responses": {
                    "200": {
                        "description": "UserResponse",
                        "schema": {
                            "$ref": "#/definitions/UserResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Initiate Authenticated User Deletion",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Delete",
                "parameters": [
                    {
                        "type": "string",
                        "description": "code if already have",
                        "name": "code",
                        "in": "query"
                    }
                ],
                "responses": {}
            },
            "patch": {
                "description": "Partially Update User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Update User",
                "parameters": [
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
        },
        "/api/auth/verify": {
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
                "summary": "Verify Email",
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
        "/api/auth/{provider}": {
            "get": {
                "description": "Authenticates a user with a specified oauth provider",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "OAuth Provider",
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
        "/api/docs": {
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
        "/api/health": {
            "get": {
                "description": "Api Health Check",
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
        }
    },
    "definitions": {
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
                "message": {
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
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "minLength": 3
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "UpdateUserParams": {
            "type": "object",
            "properties": {
                "currentPassword": {
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
                "newPassword": {
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
	Title:            "Goth",
	Description:      "Goth is a robust authentication API that supports both OAuth providers and traditional credentials-based authentication. It is designed to provide secure and flexible user authentication for various applications.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
