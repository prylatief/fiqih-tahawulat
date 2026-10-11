package mcp

import (
	"encoding/json"

	"github.com/brolatief/fiqih-tahawulat/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
)

type Server struct {
	svc service.Service
}

func NewServer(svc service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Tools() []mcp.Tool {
	return []mcp.Tool{
		mcp.NewTool("list_kitab",
			mcp.WithDescription("Browse kitab catalog with pagination"),
			mcp.WithString("judul", mcp.Description("Filter by title")),
			mcp.WithNumber("page", mcp.Description("Page number")),
			mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
		),
		mcp.NewTool("get_kitab",
			mcp.WithDescription("Get metadata for a single kitab by ID"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Kitab ID")),
		),
		mcp.NewTool("get_konten",
			mcp.WithDescription("Read paginated content sections"),
			mcp.WithNumber("kitab_id", mcp.Required(), mcp.Description("Kitab ID")),
			mcp.WithNumber("page", mcp.Description("Page number")),
			mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
			mcp.WithString("lang", mcp.Description("Translation language (e.g. 'id' for Indonesian)")),
		),
		mcp.NewTool("search_kitab",
			mcp.WithDescription("Full-text search across all library content"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Arabic search query")),
			mcp.WithString("kategori", mcp.Description("Filter by category")),
			mcp.WithString("kitab", mcp.Description("Filter by kitab title")),
			mcp.WithNumber("page", mcp.Description("Page number")),
			mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
		),
		mcp.NewTool("list_kategori",
			mcp.WithDescription("List all available categories"),
		),
	}
}

func jsonResult(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
