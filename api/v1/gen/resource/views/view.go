// Code generated by goa v3.14.1, DO NOT EDIT.
//
// resource views
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/v1/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// Resources is the viewed result type that is projected based on a view.
type Resources struct {
	// Type to project
	Projected *ResourcesView
	// View to render
	View string
}

// ResourceVersions is the viewed result type that is projected based on a view.
type ResourceVersions struct {
	// Type to project
	Projected *ResourceVersionsView
	// View to render
	View string
}

// ResourceVersion is the viewed result type that is projected based on a view.
type ResourceVersion struct {
	// Type to project
	Projected *ResourceVersionView
	// View to render
	View string
}

// ResourceVersionReadme is the viewed result type that is projected based on a
// view.
type ResourceVersionReadme struct {
	// Type to project
	Projected *ResourceVersionReadmeView
	// View to render
	View string
}

// ResourceVersionYaml is the viewed result type that is projected based on a
// view.
type ResourceVersionYaml struct {
	// Type to project
	Projected *ResourceVersionYamlView
	// View to render
	View string
}

// Resource is the viewed result type that is projected based on a view.
type Resource struct {
	// Type to project
	Projected *ResourceView
	// View to render
	View string
}

// ResourcesView is a type that runs validations on a projected type.
type ResourcesView struct {
	Data ResourceDataCollectionView
}

// ResourceDataCollectionView is a type that runs validations on a projected
// type.
type ResourceDataCollectionView []*ResourceDataView

// ResourceDataView is a type that runs validations on a projected type.
type ResourceDataView struct {
	// ID is the unique id of the resource
	ID *uint
	// Name of resource
	Name *string
	// Type of catalog to which resource belongs
	Catalog *CatalogView
	// Categories related to resource
	Categories []*CategoryView
	// Kind of resource
	Kind *string
	// Url path of the resource in hub
	HubURLPath *string
	// Path of the api to get the raw yaml of resource from hub apiserver
	HubRawURLPath *string
	// Latest version of resource
	LatestVersion *ResourceVersionDataView
	// Tags related to resource
	Tags []*TagView
	// Platforms related to resource
	Platforms []*PlatformView
	// Rating of resource
	Rating *float64
	// List of all versions of a resource
	Versions []*ResourceVersionDataView
}

// CatalogView is a type that runs validations on a projected type.
type CatalogView struct {
	// ID is the unique id of the catalog
	ID *uint
	// Name of catalog
	Name *string
	// Type of catalog
	Type *string
	// URL of catalog
	URL *string
	// Provider of catalog
	Provider *string
}

// CategoryView is a type that runs validations on a projected type.
type CategoryView struct {
	// ID is the unique id of the category
	ID *uint
	// Name of category
	Name *string
}

// ResourceVersionDataView is a type that runs validations on a projected type.
type ResourceVersionDataView struct {
	// ID is the unique id of resource's version
	ID *uint
	// Version of resource
	Version *string
	// Display name of version
	DisplayName *string
	// Deprecation status of a version
	Deprecated *bool
	// Description of version
	Description *string
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion *string
	// Raw URL of resource's yaml file of the version
	RawURL *string
	// Web URL of resource's yaml file of the version
	WebURL *string
	// Path of the api to get the raw yaml of resource from hub apiserver
	HubRawURLPath *string
	// Timestamp when version was last updated
	UpdatedAt *string
	// Platforms related to resource version
	Platforms []*PlatformView
	// Url path of the resource in hub
	HubURLPath *string
	// Resource to which the version belongs
	Resource *ResourceDataView
}

// PlatformView is a type that runs validations on a projected type.
type PlatformView struct {
	// ID is the unique id of platform
	ID *uint
	// Name of platform
	Name *string
}

// TagView is a type that runs validations on a projected type.
type TagView struct {
	// ID is the unique id of tag
	ID *uint
	// Name of tag
	Name *string
}

// ResourceVersionsView is a type that runs validations on a projected type.
type ResourceVersionsView struct {
	Data *VersionsView
}

// VersionsView is a type that runs validations on a projected type.
type VersionsView struct {
	// Latest Version of resource
	Latest *ResourceVersionDataView
	// List of all versions of resource
	Versions []*ResourceVersionDataView
}

// ResourceVersionView is a type that runs validations on a projected type.
type ResourceVersionView struct {
	Data *ResourceVersionDataView
}

// ResourceVersionReadmeView is a type that runs validations on a projected
// type.
type ResourceVersionReadmeView struct {
	Data *ResourceContentView
}

// ResourceContentView is a type that runs validations on a projected type.
type ResourceContentView struct {
	// Readme
	Readme *string
	// Yaml
	Yaml *string
}

// ResourceVersionYamlView is a type that runs validations on a projected type.
type ResourceVersionYamlView struct {
	Data *ResourceContentView
}

// ResourceView is a type that runs validations on a projected type.
type ResourceView struct {
	Data *ResourceDataView
}

var (
	// ResourcesMap is a map indexing the attribute names of Resources by view name.
	ResourcesMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceVersionsMap is a map indexing the attribute names of
	// ResourceVersions by view name.
	ResourceVersionsMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceVersionMap is a map indexing the attribute names of ResourceVersion
	// by view name.
	ResourceVersionMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceVersionReadmeMap is a map indexing the attribute names of
	// ResourceVersionReadme by view name.
	ResourceVersionReadmeMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceVersionYamlMap is a map indexing the attribute names of
	// ResourceVersionYaml by view name.
	ResourceVersionYamlMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceMap is a map indexing the attribute names of Resource by view name.
	ResourceMap = map[string][]string{
		"default": {
			"data",
		},
	}
	// ResourceDataCollectionMap is a map indexing the attribute names of
	// ResourceDataCollection by view name.
	ResourceDataCollectionMap = map[string][]string{
		"info": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"tags",
			"platforms",
			"rating",
		},
		"withoutVersion": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"hubRawURLPath",
			"latestVersion",
			"tags",
			"platforms",
			"rating",
		},
		"default": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"hubRawURLPath",
			"latestVersion",
			"tags",
			"platforms",
			"rating",
			"versions",
		},
	}
	// ResourceDataMap is a map indexing the attribute names of ResourceData by
	// view name.
	ResourceDataMap = map[string][]string{
		"info": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"tags",
			"platforms",
			"rating",
		},
		"withoutVersion": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"hubRawURLPath",
			"latestVersion",
			"tags",
			"platforms",
			"rating",
		},
		"default": {
			"id",
			"name",
			"catalog",
			"categories",
			"kind",
			"hubURLPath",
			"hubRawURLPath",
			"latestVersion",
			"tags",
			"platforms",
			"rating",
			"versions",
		},
	}
	// CatalogMap is a map indexing the attribute names of Catalog by view name.
	CatalogMap = map[string][]string{
		"min": {
			"id",
			"name",
			"type",
		},
		"default": {
			"id",
			"name",
			"type",
			"url",
			"provider",
		},
	}
	// ResourceVersionDataMap is a map indexing the attribute names of
	// ResourceVersionData by view name.
	ResourceVersionDataMap = map[string][]string{
		"tiny": {
			"id",
			"version",
		},
		"min": {
			"id",
			"version",
			"rawURL",
			"webURL",
			"hubRawURLPath",
			"hubURLPath",
			"platforms",
		},
		"withoutResource": {
			"id",
			"version",
			"displayName",
			"deprecated",
			"description",
			"minPipelinesVersion",
			"rawURL",
			"webURL",
			"hubRawURLPath",
			"hubURLPath",
			"updatedAt",
			"platforms",
		},
		"default": {
			"id",
			"version",
			"displayName",
			"deprecated",
			"description",
			"minPipelinesVersion",
			"rawURL",
			"webURL",
			"hubURLPath",
			"hubRawURLPath",
			"updatedAt",
			"resource",
			"platforms",
		},
	}
	// VersionsMap is a map indexing the attribute names of Versions by view name.
	VersionsMap = map[string][]string{
		"default": {
			"latest",
			"versions",
		},
	}
	// ResourceContentMap is a map indexing the attribute names of ResourceContent
	// by view name.
	ResourceContentMap = map[string][]string{
		"readme": {
			"readme",
		},
		"yaml": {
			"yaml",
		},
		"default": {
			"readme",
			"yaml",
		},
	}
)

// ValidateResources runs the validations defined on the viewed result type
// Resources.
func ValidateResources(result *Resources) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourcesView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResourceVersions runs the validations defined on the viewed result
// type ResourceVersions.
func ValidateResourceVersions(result *ResourceVersions) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourceVersionsView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResourceVersion runs the validations defined on the viewed result
// type ResourceVersion.
func ValidateResourceVersion(result *ResourceVersion) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourceVersionView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResourceVersionReadme runs the validations defined on the viewed
// result type ResourceVersionReadme.
func ValidateResourceVersionReadme(result *ResourceVersionReadme) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourceVersionReadmeView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResourceVersionYaml runs the validations defined on the viewed
// result type ResourceVersionYaml.
func ValidateResourceVersionYaml(result *ResourceVersionYaml) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourceVersionYamlView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResource runs the validations defined on the viewed result type
// Resource.
func ValidateResource(result *Resource) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateResourceView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []any{"default"})
	}
	return
}

// ValidateResourcesView runs the validations defined on ResourcesView using
// the "default" view.
func ValidateResourcesView(result *ResourcesView) (err error) {

	if result.Data != nil {
		if err2 := ValidateResourceDataCollectionViewWithoutVersion(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataCollectionViewInfo runs the validations defined on
// ResourceDataCollectionView using the "info" view.
func ValidateResourceDataCollectionViewInfo(result ResourceDataCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateResourceDataViewInfo(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataCollectionViewWithoutVersion runs the validations
// defined on ResourceDataCollectionView using the "withoutVersion" view.
func ValidateResourceDataCollectionViewWithoutVersion(result ResourceDataCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateResourceDataViewWithoutVersion(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataCollectionView runs the validations defined on
// ResourceDataCollectionView using the "default" view.
func ValidateResourceDataCollectionView(result ResourceDataCollectionView) (err error) {
	for _, item := range result {
		if err2 := ValidateResourceDataView(item); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataViewInfo runs the validations defined on
// ResourceDataView using the "info" view.
func ValidateResourceDataViewInfo(result *ResourceDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Categories == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("categories", "result"))
	}
	if result.Kind == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("kind", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.Rating == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rating", "result"))
	}
	for _, e := range result.Categories {
		if e != nil {
			if err2 := ValidateCategoryView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Tags {
		if e != nil {
			if err2 := ValidateTagView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Catalog != nil {
		if err2 := ValidateCatalogViewMin(result.Catalog); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataViewWithoutVersion runs the validations defined on
// ResourceDataView using the "withoutVersion" view.
func ValidateResourceDataViewWithoutVersion(result *ResourceDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Categories == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("categories", "result"))
	}
	if result.Kind == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("kind", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.Rating == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rating", "result"))
	}
	if result.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "result"))
	}
	for _, e := range result.Categories {
		if e != nil {
			if err2 := ValidateCategoryView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Tags {
		if e != nil {
			if err2 := ValidateTagView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Catalog != nil {
		if err2 := ValidateCatalogViewMin(result.Catalog); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.LatestVersion != nil {
		if err2 := ValidateResourceVersionDataViewWithoutResource(result.LatestVersion); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceDataView runs the validations defined on ResourceDataView
// using the "default" view.
func ValidateResourceDataView(result *ResourceDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Categories == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("categories", "result"))
	}
	if result.Kind == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("kind", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.Rating == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rating", "result"))
	}
	if result.Versions == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("versions", "result"))
	}
	if result.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "result"))
	}
	for _, e := range result.Categories {
		if e != nil {
			if err2 := ValidateCategoryView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Tags {
		if e != nil {
			if err2 := ValidateTagView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range result.Versions {
		if e != nil {
			if err2 := ValidateResourceVersionDataView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Catalog != nil {
		if err2 := ValidateCatalogViewMin(result.Catalog); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if result.LatestVersion != nil {
		if err2 := ValidateResourceVersionDataViewWithoutResource(result.LatestVersion); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateCatalogViewMin runs the validations defined on CatalogView using the
// "min" view.
func ValidateCatalogViewMin(result *CatalogView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "result"))
	}
	if result.Type != nil {
		if !(*result.Type == "official" || *result.Type == "community") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.type", *result.Type, []any{"official", "community"}))
		}
	}
	return
}

// ValidateCatalogView runs the validations defined on CatalogView using the
// "default" view.
func ValidateCatalogView(result *CatalogView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "result"))
	}
	if result.URL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("url", "result"))
	}
	if result.Provider == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("provider", "result"))
	}
	if result.Type != nil {
		if !(*result.Type == "official" || *result.Type == "community") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("result.type", *result.Type, []any{"official", "community"}))
		}
	}
	return
}

// ValidateCategoryView runs the validations defined on CategoryView.
func ValidateCategoryView(result *CategoryView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	return
}

// ValidateResourceVersionDataViewTiny runs the validations defined on
// ResourceVersionDataView using the "tiny" view.
func ValidateResourceVersionDataViewTiny(result *ResourceVersionDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "result"))
	}
	return
}

// ValidateResourceVersionDataViewMin runs the validations defined on
// ResourceVersionDataView using the "min" view.
func ValidateResourceVersionDataViewMin(result *ResourceVersionDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "result"))
	}
	if result.RawURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rawURL", "result"))
	}
	if result.WebURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("webURL", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "result"))
	}
	if result.RawURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.rawURL", *result.RawURL, goa.FormatURI))
	}
	if result.WebURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.webURL", *result.WebURL, goa.FormatURI))
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateResourceVersionDataViewWithoutResource runs the validations defined
// on ResourceVersionDataView using the "withoutResource" view.
func ValidateResourceVersionDataViewWithoutResource(result *ResourceVersionDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "result"))
	}
	if result.DisplayName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("displayName", "result"))
	}
	if result.Description == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("description", "result"))
	}
	if result.MinPipelinesVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minPipelinesVersion", "result"))
	}
	if result.RawURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rawURL", "result"))
	}
	if result.WebURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("webURL", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "result"))
	}
	if result.RawURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.rawURL", *result.RawURL, goa.FormatURI))
	}
	if result.WebURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.webURL", *result.WebURL, goa.FormatURI))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateResourceVersionDataView runs the validations defined on
// ResourceVersionDataView using the "default" view.
func ValidateResourceVersionDataView(result *ResourceVersionDataView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "result"))
	}
	if result.DisplayName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("displayName", "result"))
	}
	if result.Description == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("description", "result"))
	}
	if result.MinPipelinesVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minPipelinesVersion", "result"))
	}
	if result.RawURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rawURL", "result"))
	}
	if result.WebURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("webURL", "result"))
	}
	if result.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "result"))
	}
	if result.Platforms == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("platforms", "result"))
	}
	if result.HubURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubURLPath", "result"))
	}
	if result.HubRawURLPath == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("hubRawURLPath", "result"))
	}
	if result.RawURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.rawURL", *result.RawURL, goa.FormatURI))
	}
	if result.WebURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.webURL", *result.WebURL, goa.FormatURI))
	}
	if result.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("result.updatedAt", *result.UpdatedAt, goa.FormatDateTime))
	}
	for _, e := range result.Platforms {
		if e != nil {
			if err2 := ValidatePlatformView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Resource != nil {
		if err2 := ValidateResourceDataViewInfo(result.Resource); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidatePlatformView runs the validations defined on PlatformView.
func ValidatePlatformView(result *PlatformView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	return
}

// ValidateTagView runs the validations defined on TagView.
func ValidateTagView(result *TagView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	return
}

// ValidateResourceVersionsView runs the validations defined on
// ResourceVersionsView using the "default" view.
func ValidateResourceVersionsView(result *ResourceVersionsView) (err error) {

	if result.Data != nil {
		if err2 := ValidateVersionsView(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateVersionsView runs the validations defined on VersionsView using the
// "default" view.
func ValidateVersionsView(result *VersionsView) (err error) {
	if result.Versions == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("versions", "result"))
	}
	for _, e := range result.Versions {
		if e != nil {
			if err2 := ValidateResourceVersionDataView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	if result.Latest != nil {
		if err2 := ValidateResourceVersionDataViewMin(result.Latest); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceVersionView runs the validations defined on
// ResourceVersionView using the "default" view.
func ValidateResourceVersionView(result *ResourceVersionView) (err error) {

	if result.Data != nil {
		if err2 := ValidateResourceVersionDataView(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceVersionReadmeView runs the validations defined on
// ResourceVersionReadmeView using the "default" view.
func ValidateResourceVersionReadmeView(result *ResourceVersionReadmeView) (err error) {

	if result.Data != nil {
		if err2 := ValidateResourceContentViewReadme(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceContentViewReadme runs the validations defined on
// ResourceContentView using the "readme" view.
func ValidateResourceContentViewReadme(result *ResourceContentView) (err error) {

	return
}

// ValidateResourceContentViewYaml runs the validations defined on
// ResourceContentView using the "yaml" view.
func ValidateResourceContentViewYaml(result *ResourceContentView) (err error) {

	return
}

// ValidateResourceContentView runs the validations defined on
// ResourceContentView using the "default" view.
func ValidateResourceContentView(result *ResourceContentView) (err error) {

	return
}

// ValidateResourceVersionYamlView runs the validations defined on
// ResourceVersionYamlView using the "default" view.
func ValidateResourceVersionYamlView(result *ResourceVersionYamlView) (err error) {

	if result.Data != nil {
		if err2 := ValidateResourceContentViewYaml(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateResourceView runs the validations defined on ResourceView using the
// "default" view.
func ValidateResourceView(result *ResourceView) (err error) {

	if result.Data != nil {
		if err2 := ValidateResourceDataView(result.Data); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}
