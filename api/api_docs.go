package api

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/sirjager/gopkg/httpx"
)

// @Summary	Documentation
// @Tags		System
// @Produce	html
// @Success	200	{string}	string	"HTML content"
// @Router		/api/docs [get]
func (s *Server) swaggerDocs(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL:            fmt.Sprintf("http://localhost:%d/api/docs/swagger.json", s.Config().Port),
		HideDownloadButton: true,
		DarkMode:           true,
		ShowSidebar:        true,
		Theme:              scalar.ThemeBluePlanet,
		Layout:             scalar.LayoutModern,
		CustomOptions: scalar.CustomOptions{
			PageTitle: fmt.Sprintf("%s Api Documentations", s.Config().ServiceName),
		},
	})
	if err != nil {
		httpx.Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	if _, writeErr := w.Write([]byte(htmlContent)); writeErr != nil {
		http.Error(w, writeErr.Error(), 500)
		return
	}
}
