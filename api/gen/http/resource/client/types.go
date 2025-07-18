// Code generated by goa v3.21.1, DO NOT EDIT.
//
// resource HTTP client types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	resource "github.com/tektoncd/hub/api/gen/resource"
	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goa "goa.design/goa/v3/pkg"
)

// ListResponseBody is the type of the "resource" service "List" endpoint HTTP
// response body.
type ListResponseBody struct {
	Data ResourceDataCollectionResponseBody `form:"data,omitempty" json:"data,omitempty" xml:"data,omitempty"`
}

// ResourceDataCollectionResponseBody is used to define fields on response body
// types.
type ResourceDataCollectionResponseBody []*ResourceDataResponseBody

// ResourceDataResponseBody is used to define fields on response body types.
type ResourceDataResponseBody struct {
	// ID is the unique id of the resource
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of resource
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// Type of catalog to which resource belongs
	Catalog *CatalogResponseBody `form:"catalog,omitempty" json:"catalog,omitempty" xml:"catalog,omitempty"`
	// Categories related to resource
	Categories []*CategoryResponseBody `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Kind of resource
	Kind *string `form:"kind,omitempty" json:"kind,omitempty" xml:"kind,omitempty"`
	// Url path of the resource in hub
	HubURLPath *string `form:"hubURLPath,omitempty" json:"hubURLPath,omitempty" xml:"hubURLPath,omitempty"`
	// Path of the api to get the raw yaml of resource from hub apiserver
	HubRawURLPath *string `form:"hubRawURLPath,omitempty" json:"hubRawURLPath,omitempty" xml:"hubRawURLPath,omitempty"`
	// Latest version of resource
	LatestVersion *ResourceVersionDataResponseBody `form:"latestVersion,omitempty" json:"latestVersion,omitempty" xml:"latestVersion,omitempty"`
	// Tags related to resource
	Tags []*TagResponseBody `form:"tags,omitempty" json:"tags,omitempty" xml:"tags,omitempty"`
	// Platforms related to resource
	Platforms []*PlatformResponseBody `form:"platforms,omitempty" json:"platforms,omitempty" xml:"platforms,omitempty"`
	// Rating of resource
	Rating *float64 `form:"rating,omitempty" json:"rating,omitempty" xml:"rating,omitempty"`
	// List of all versions of a resource
	Versions []*ResourceVersionDataResponseBody `form:"versions,omitempty" json:"versions,omitempty" xml:"versions,omitempty"`
}

// CatalogResponseBody is used to define fields on response body types.
type CatalogResponseBody struct {
	// ID is the unique id of the catalog
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of catalog
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// Type of catalog
	Type *string `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
	// URL of catalog
	URL *string `form:"url,omitempty" json:"url,omitempty" xml:"url,omitempty"`
	// Provider of catalog
	Provider *string `form:"provider,omitempty" json:"provider,omitempty" xml:"provider,omitempty"`
}

// CategoryResponseBody is used to define fields on response body types.
type CategoryResponseBody struct {
	// ID is the unique id of the category
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of category
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// ResourceVersionDataResponseBody is used to define fields on response body
// types.
type ResourceVersionDataResponseBody struct {
	// ID is the unique id of resource's version
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Version of resource
	Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	// Display name of version
	DisplayName *string `form:"displayName,omitempty" json:"displayName,omitempty" xml:"displayName,omitempty"`
	// Deprecation status of a version
	Deprecated *bool `form:"deprecated,omitempty" json:"deprecated,omitempty" xml:"deprecated,omitempty"`
	// Description of version
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion *string `form:"minPipelinesVersion,omitempty" json:"minPipelinesVersion,omitempty" xml:"minPipelinesVersion,omitempty"`
	// Raw URL of resource's yaml file of the version
	RawURL *string `form:"rawURL,omitempty" json:"rawURL,omitempty" xml:"rawURL,omitempty"`
	// Web URL of resource's yaml file of the version
	WebURL *string `form:"webURL,omitempty" json:"webURL,omitempty" xml:"webURL,omitempty"`
	// Path of the api to get the raw yaml of resource from hub apiserver
	HubRawURLPath *string `form:"hubRawURLPath,omitempty" json:"hubRawURLPath,omitempty" xml:"hubRawURLPath,omitempty"`
	// Timestamp when version was last updated
	UpdatedAt *string `form:"updatedAt,omitempty" json:"updatedAt,omitempty" xml:"updatedAt,omitempty"`
	// Platforms related to resource version
	Platforms []*PlatformResponseBody `form:"platforms,omitempty" json:"platforms,omitempty" xml:"platforms,omitempty"`
	// Url path of the resource in hub
	HubURLPath *string `form:"hubURLPath,omitempty" json:"hubURLPath,omitempty" xml:"hubURLPath,omitempty"`
	// Resource to which the version belongs
	Resource *ResourceDataResponseBody `form:"resource,omitempty" json:"resource,omitempty" xml:"resource,omitempty"`
}

// PlatformResponseBody is used to define fields on response body types.
type PlatformResponseBody struct {
	// ID is the unique id of platform
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of platform
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// TagResponseBody is used to define fields on response body types.
type TagResponseBody struct {
	// ID is the unique id of tag
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of tag
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// NewQueryResultFound builds a "resource" service "Query" endpoint result from
// a HTTP "Found" response.
func NewQueryResultFound(location string) *resource.QueryResult {
	v := &resource.QueryResult{}
	v.Location = location

	return v
}

// NewListResourcesMovedPermanently builds a "resource" service "List" endpoint
// result from a HTTP "MovedPermanently" response.
func NewListResourcesMovedPermanently(body *ListResponseBody) *resourceviews.ResourcesView {
	v := &resourceviews.ResourcesView{}
	v.Data = make([]*resourceviews.ResourceDataView, len(body.Data))
	for i, val := range body.Data {
		v.Data[i] = unmarshalResourceDataResponseBodyToResourceviewsResourceDataView(val)
	}

	return v
}

// NewVersionsByIDResultFound builds a "resource" service "VersionsByID"
// endpoint result from a HTTP "Found" response.
func NewVersionsByIDResultFound(location string) *resource.VersionsByIDResult {
	v := &resource.VersionsByIDResult{}
	v.Location = location

	return v
}

// NewByCatalogKindNameVersionResultFound builds a "resource" service
// "ByCatalogKindNameVersion" endpoint result from a HTTP "Found" response.
func NewByCatalogKindNameVersionResultFound(location string) *resource.ByCatalogKindNameVersionResult {
	v := &resource.ByCatalogKindNameVersionResult{}
	v.Location = location

	return v
}

// NewByVersionIDResultFound builds a "resource" service "ByVersionId" endpoint
// result from a HTTP "Found" response.
func NewByVersionIDResultFound(location string) *resource.ByVersionIDResult {
	v := &resource.ByVersionIDResult{}
	v.Location = location

	return v
}

// NewByCatalogKindNameResultFound builds a "resource" service
// "ByCatalogKindName" endpoint result from a HTTP "Found" response.
func NewByCatalogKindNameResultFound(location string) *resource.ByCatalogKindNameResult {
	v := &resource.ByCatalogKindNameResult{}
	v.Location = location

	return v
}

// NewByIDResultFound builds a "resource" service "ById" endpoint result from a
// HTTP "Found" response.
func NewByIDResultFound(location string) *resource.ByIDResult {
	v := &resource.ByIDResult{}
	v.Location = location

	return v
}

// ValidateResourceDataCollectionResponseBody runs the validations defined on
// ResourceDataCollectionResponseBody
func ValidateResourceDataCollectionResponseBody(body ResourceDataCollectionResponseBody) (err error) {
	for _, e := range body {
		if e != nil {
			if err2 := ValidateResourceDataResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateResourceDataResponseBody runs the validations defined on
// ResourceDataResponseBody
func ValidateResourceDataResponseBody(body *ResourceDataResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.Catalog == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("catalog", "body"))
	}
	if body.Categories == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("categories", "body"))
	}
	if body.Kind == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("kind", "body"))
	}
	if body.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "body"))
	}
	if body.LatestVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("latestVersion", "body"))
	}
	if body.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "body"))
	}
	if body.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "body"))
	}
	if body.Rating == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rating", "body"))
	}
	if body.Versions == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("versions", "body"))
	}
	if body.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "body"))
	}
	if body.Catalog != nil {
		if err2 := ValidateCatalogResponseBody(body.Catalog); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range body.Categories {
		if e != nil {
			if err2 := ValidateCategoryResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if body.LatestVersion != nil {
		if err2 := ValidateResourceVersionDataResponseBody(body.LatestVersion); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range body.Tags {
		if e != nil {
			if err2 := ValidateTagResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.Platforms {
		if e != nil {
			if err2 := ValidatePlatformResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.Versions {
		if e != nil {
			if err2 := ValidateResourceVersionDataResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateCatalogResponseBody runs the validations defined on
// CatalogResponseBody
func ValidateCatalogResponseBody(body *CatalogResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "body"))
	}
	if body.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "body"))
	}
	if body.Provider == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("provider", "body"))
	}
	if body.Type != nil {
		if !(*body.Type == "official" || *body.Type == "community") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.type", *body.Type, []any{"official", "community"}))
		}
	}
	return
}

// ValidateCategoryResponseBody runs the validations defined on
// CategoryResponseBody
func ValidateCategoryResponseBody(body *CategoryResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}

// ValidateResourceVersionDataResponseBody runs the validations defined on
// ResourceVersionDataResponseBody
func ValidateResourceVersionDataResponseBody(body *ResourceVersionDataResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "body"))
	}
	if body.DisplayName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("displayName", "body"))
	}
	if body.Description == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("description", "body"))
	}
	if body.MinPipelinesVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minPipelinesVersion", "body"))
	}
	if body.RawURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rawURL", "body"))
	}
	if body.WebURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("webURL", "body"))
	}
	if body.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "body"))
	}
	if body.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "body"))
	}
	if body.Resource == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("resource", "body"))
	}
	if body.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "body"))
	}
	if body.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "body"))
	}
	if body.RawURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.rawURL", *body.RawURL, goa.FormatURI))
	}
	if body.WebURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.webURL", *body.WebURL, goa.FormatURI))
	}
	if body.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.updatedAt", *body.UpdatedAt, goa.FormatDateTime))
	}
	for _, e := range body.Platforms {
		if e != nil {
			if err2 := ValidatePlatformResponseBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if body.Resource != nil {
		if err2 := ValidateResourceDataResponseBody(body.Resource); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidatePlatformResponseBody runs the validations defined on
// PlatformResponseBody
func ValidatePlatformResponseBody(body *PlatformResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}

// ValidateTagResponseBody runs the validations defined on TagResponseBody
func ValidateTagResponseBody(body *TagResponseBody) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}
