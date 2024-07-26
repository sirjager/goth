package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func getPageAndLimitFromRequest(r *http.Request, defaultPage, defaultLimit *int) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	if (pageParam) != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			*defaultPage = p
		}
	}
	if (limitParam) != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			*defaultLimit = l
		}
	}
}

func (a *API) ParseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func (a *API) ParseAndValidate(r *http.Request, v interface{}) error {
	if err := a.ParseJSON(r, v); err != nil {
		return err
	}
	if err := a.validate.Struct(v); err != nil {
		return err
	}
	return nil
}
