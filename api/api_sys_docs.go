package api

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

// @Summary	Documentation
// @Tags		System
// @Produce	html
// @Success	200	{string}	string	"HTML content"
// @Router		/swagger [get]
func (a *API) SysDocs(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL:            "./docs/swagger.json",
		HideDownloadButton: true,
		DarkMode:           true,
		Theme:              scalar.ThemeSolarized,
		Layout:             scalar.LayoutModern,
		CustomOptions: scalar.CustomOptions{
			PageTitle: fmt.Sprintf("%s Api Documentations", a.config.ServiceName),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	if _, writeErr := w.Write([]byte(htmlContent)); writeErr != nil {
		http.Error(w, writeErr.Error(), 500)
		return
	}
}
