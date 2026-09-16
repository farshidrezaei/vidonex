// Package api provides HTTP request handlers and OpenAPI specifications for the Vidonex server.
package api

import (
	"net/http"
)

// OpenAPISpecDefinition returns the OpenAPI 3.0 specification for Vidonex REST & WebSocket endpoints.
const OpenAPISpecDefinition = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Vidonex Video Composition Engine API",
    "version": "1.6.0",
    "description": "RESTful and WebSocket API for declarative video composition, intermediate representation DAG compilation, timeline management, and real-time FFmpeg telemetry.",
    "contact": {
      "name": "Vidonex Engineering",
      "url": "https://github.com/farshidrezaei/vidonex"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "servers": [
    {
      "url": "/",
      "description": "Local Vidonex Studio Server"
    }
  ],
  "paths": {
    "/api/projects": {
      "get": {
        "summary": "List all video composition projects",
        "operationId": "listProjects",
        "responses": {
          "200": {
            "description": "Array of projects",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": { "$ref": "#/components/schemas/Project" }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create a new project",
        "operationId": "createProject",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/CreateProjectRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Project created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/Project" }
              }
            }
          }
        }
      }
    },
    "/api/projects/{id}": {
      "get": {
        "summary": "Retrieve project details and timeline specification",
        "operationId": "getProject",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Project found" },
          "404": { "description": "Project not found" }
        }
      },
      "put": {
        "summary": "Update project timeline specification",
        "operationId": "updateProject",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Project updated" },
          "404": { "description": "Project not found" }
        }
      },
      "delete": {
        "summary": "Delete a project",
        "operationId": "deleteProject",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Project deleted" },
          "404": { "description": "Project not found" }
        }
      }
    },
    "/api/media": {
      "get": {
        "summary": "List uploaded media assets",
        "operationId": "listMedia",
        "responses": {
          "200": { "description": "List of media assets" }
        }
      }
    },
    "/api/media/upload": {
      "post": {
        "summary": "Upload a media file (video, audio, image, subtitle)",
        "operationId": "uploadMedia",
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "properties": {
                  "file": { "type": "string", "format": "binary" }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Media uploaded and probed" }
        }
      }
    },
    "/api/spec/validate": {
      "post": {
        "summary": "Validate a declarative VideoSpec timeline AST",
        "operationId": "validateSpec",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "type": "object" }
            }
          }
        },
        "responses": {
          "200": { "description": "Specification validation result" }
        }
      }
    },
    "/api/spec/graph": {
      "post": {
        "summary": "Generate a Mermaid.js or Graphviz DOT diagram from timeline AST",
        "operationId": "generateGraph",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "type": "object" }
            }
          }
        },
        "responses": {
          "200": { "description": "Filtergraph diagram code" }
        }
      }
    },
    "/api/render/start": {
      "post": {
        "summary": "Queue and start a background FFmpeg render job",
        "operationId": "startRender",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "project_id": { "type": "string" },
                  "output_path": { "type": "string" },
                  "format": { "type": "string", "enum": ["mp4", "mkv", "mov", "webm", "gif"] },
                  "gpu_accel": { "type": "string", "enum": ["auto", "none", "nvenc", "videotoolbox", "qsv", "vaapi"] }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Render job queued" }
        }
      }
    },
    "/api/render/{jobId}": {
      "get": {
        "summary": "Get render job progress status and telemetry",
        "operationId": "getRenderJob",
        "parameters": [
          { "name": "jobId", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Job status and progress" },
          "404": { "description": "Job not found" }
        }
      }
    },
    "/ws/telemetry": {
      "get": {
        "summary": "WebSocket streaming real-time render telemetry",
        "operationId": "streamTelemetry",
        "responses": {
          "101": { "description": "Switching Protocols to WebSocket" }
        }
      }
    },
    "/api/version/check": {
      "get": {
        "summary": "Check for engine updates and retrieve host system diagnostics",
        "operationId": "checkVersion",
        "parameters": [
          { "name": "force", "in": "query", "required": false, "schema": { "type": "boolean" }, "description": "Bypass 10-minute cache and force GitHub API query" }
        ],
        "responses": {
          "200": { "description": "Version check and system hardware diagnostics" }
        }
      }
    },
    "/api/version/upgrade": {
      "post": {
        "summary": "Trigger in-app self-update with SSE progress streaming",
        "operationId": "upgradeEngine",
        "parameters": [
          { "name": "force", "in": "query", "required": false, "schema": { "type": "boolean" }, "description": "Force upgrade even if already on latest version" }
        ],
        "responses": {
          "200": { "description": "Server-Sent Events (SSE) stream of update progress" }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Project": {
        "type": "object",
        "properties": {
          "id": { "type": "string" },
          "name": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "CreateProjectRequest": {
        "type": "object",
        "required": ["name"],
        "properties": {
          "name": { "type": "string" },
          "description": { "type": "string" }
        }
      }
    }
  }
}`

// ScalarHTML is the HTML payload embedding the standalone Scalar API documentation sandbox.
const ScalarHTML = `<!doctype html>
<html>
  <head>
    <title>Vidonex API Reference & Interactive Sandbox</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="icon" type="image/png" href="/favicon-32x32.png" />
    <style>
      body {
        margin: 0;
        background-color: #0b0f19;
      }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/openapi.json"
      data-configuration='{"theme":"purple","layout":"modern","darkMode":true}'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

// HandleOpenAPISpec serves the OpenAPI 3.0 JSON specification.
func (handlers *Handlers) HandleOpenAPISpec(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(OpenAPISpecDefinition))
}

// HandleScalarDocs serves the interactive Scalar documentation HTML page.
func (handlers *Handlers) HandleScalarDocs(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(ScalarHTML))
}
