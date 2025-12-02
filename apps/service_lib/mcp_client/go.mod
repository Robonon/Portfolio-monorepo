module mcp_client

go 1.25.1

replace api => ../api

require (
	api v0.0.0-00010101000000-000000000000
	github.com/modelcontextprotocol/go-sdk v1.0.0
)

require (
	github.com/google/jsonschema-go v0.3.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)
