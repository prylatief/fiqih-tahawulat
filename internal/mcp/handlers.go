package mcp

import (
	"context"
	"fmt"

	"github.com/brolatief/fiqih-tahawulat/internal/repository"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) RegisterHandlers(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool("list_kitab",
		mcp.WithDescription("Browse kitab catalog with pagination"),
		mcp.WithString("judul", mcp.Description("Filter by title")),
		mcp.WithNumber("page", mcp.Description("Page number")),
		mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
	), s.handleListKitab)

	mcpServer.AddTool(mcp.NewTool("get_kitab",
		mcp.WithDescription("Get metadata for a single kitab by ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Kitab ID")),
	), s.handleGetKitab)

	mcpServer.AddTool(mcp.NewTool("get_konten",
		mcp.WithDescription("Read paginated content sections"),
		mcp.WithNumber("kitab_id", mcp.Required(), mcp.Description("Kitab ID")),
		mcp.WithNumber("page", mcp.Description("Page number")),
		mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
		mcp.WithString("lang", mcp.Description("Translation language (e.g. 'id' for Indonesian)")),
	), s.handleGetKonten)

	mcpServer.AddTool(mcp.NewTool("search_kitab",
		mcp.WithDescription("Full-text search across all library content"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Arabic search query")),
		mcp.WithString("kategori", mcp.Description("Filter by category")),
		mcp.WithString("kitab", mcp.Description("Filter by kitab title")),
		mcp.WithNumber("page", mcp.Description("Page number")),
		mcp.WithNumber("limit", mcp.Description("Items per page (max 100)")),
	), s.handleSearchKitab)

	mcpServer.AddTool(mcp.NewTool("list_kategori",
		mcp.WithDescription("List all available categories"),
	), s.handleListKategori)
}

func (s *Server) handleListKitab(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	judul, _ := req.Params.Arguments["judul"].(string)
	page := 1
	limit := 20

	if p, ok := req.Params.Arguments["page"].(float64); ok {
		page = int(p)
	}
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	filter := repository.KitabFilter{
		Judul: judul,
		Page:  page,
		Limit: limit,
	}

	data, total, err := s.svc.ListKitab(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list kitab failed: %w", err)
	}

	result := map[string]interface{}{
		"data":  data,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	return mcp.NewToolResultText(jsonResult(result)), nil
}

func (s *Server) handleGetKitab(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := int(req.Params.Arguments["id"].(float64))

	kitab, err := s.svc.GetKitab(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get kitab failed: %w", err)
	}

	return mcp.NewToolResultText(jsonResult(kitab)), nil
}

func (s *Server) handleGetKonten(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	kitabID := int(req.Params.Arguments["kitab_id"].(float64))
	page := 1
	limit := 10
	lang := ""

	if p, ok := req.Params.Arguments["page"].(float64); ok {
		page = int(p)
	}
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}
	if l, ok := req.Params.Arguments["lang"].(string); ok {
		lang = l
	}

	data, total, err := s.svc.GetKonten(ctx, kitabID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("get konten failed: %w", err)
	}

	result := map[string]interface{}{
		"data":     data,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"kitab_id": kitabID,
		"lang":     lang,
	}
	return mcp.NewToolResultText(jsonResult(result)), nil
}

func (s *Server) handleSearchKitab(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, _ := req.Params.Arguments["query"].(string)
	kategori, _ := req.Params.Arguments["kategori"].(string)
	kitab, _ := req.Params.Arguments["kitab"].(string)
	page := 1
	limit := 10

	if p, ok := req.Params.Arguments["page"].(float64); ok {
		page = int(p)
	}
	if l, ok := req.Params.Arguments["limit"].(float64); ok {
		limit = int(l)
	}

	filter := repository.SearchFilter{
		Query:    query,
		Kategori: kategori,
		Kitab:    kitab,
		Page:     page,
		Limit:    limit,
	}

	data, total, err := s.svc.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	result := map[string]interface{}{
		"data":     data,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"query":    query,
		"kategori": kategori,
		"kitab":    kitab,
	}
	return mcp.NewToolResultText(jsonResult(result)), nil
}

func (s *Server) handleListKategori(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	data, err := s.svc.ListKategori(ctx)
	if err != nil {
		return nil, fmt.Errorf("list kategori failed: %w", err)
	}

	result := map[string]interface{}{
		"data":  data,
		"total": len(data),
	}
	return mcp.NewToolResultText(jsonResult(result)), nil
}
