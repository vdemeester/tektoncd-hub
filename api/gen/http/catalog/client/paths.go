// Code generated by goa v3.21.1, DO NOT EDIT.
//
// HTTP request path constructors for the catalog service.
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"fmt"
)

// RefreshCatalogPath returns the URL path to the catalog service Refresh HTTP endpoint.
func RefreshCatalogPath(catalogName string) string {
	return fmt.Sprintf("/catalog/%v/refresh", catalogName)
}

// RefreshAllCatalogPath returns the URL path to the catalog service RefreshAll HTTP endpoint.
func RefreshAllCatalogPath() string {
	return "/catalog/refresh"
}

// CatalogErrorCatalogPath returns the URL path to the catalog service CatalogError HTTP endpoint.
func CatalogErrorCatalogPath(catalogName string) string {
	return fmt.Sprintf("/catalog/%v/error", catalogName)
}
