// Code generated by goa v3.21.1, DO NOT EDIT.
//
// hub HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	adminc "github.com/tektoncd/hub/api/gen/http/admin/client"
	catalogc "github.com/tektoncd/hub/api/gen/http/catalog/client"
	categoryc "github.com/tektoncd/hub/api/gen/http/category/client"
	ratingc "github.com/tektoncd/hub/api/gen/http/rating/client"
	resourcec "github.com/tektoncd/hub/api/gen/http/resource/client"
	statusc "github.com/tektoncd/hub/api/gen/http/status/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `admin (update-agent|refresh-config)
catalog (refresh|refresh-all|catalog-error)
category list
rating (get|update)
resource (query|list|versions-by-id|by-catalog-kind-name-version|by-version-id|by-catalog-kind-name|by-id)
status status
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` admin update-agent --body '{
      "name": "abc",
      "scopes": [
         "catalog-refresh",
         "agent:create"
      ]
   }' --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"` + "\n" +
		os.Args[0] + ` catalog refresh --catalog-name "tekton" --token "Ad dicta."` + "\n" +
		os.Args[0] + ` category list` + "\n" +
		os.Args[0] + ` rating get --id 1154472414850186965 --token "Dolorem voluptatem cum."` + "\n" +
		os.Args[0] + ` resource query --name "buildah" --catalogs '[
      "tekton",
      "openshift"
   ]' --categories '[
      "build",
      "tools"
   ]' --kinds '[
      "task",
      "pipelines"
   ]' --tags '[
      "image",
      "build"
   ]' --platforms '[
      "linux/s390x",
      "linux/amd64"
   ]' --limit 100 --match "contains"` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, any, error) {
	var (
		adminFlags = flag.NewFlagSet("admin", flag.ContinueOnError)

		adminUpdateAgentFlags     = flag.NewFlagSet("update-agent", flag.ExitOnError)
		adminUpdateAgentBodyFlag  = adminUpdateAgentFlags.String("body", "REQUIRED", "")
		adminUpdateAgentTokenFlag = adminUpdateAgentFlags.String("token", "REQUIRED", "")

		adminRefreshConfigFlags     = flag.NewFlagSet("refresh-config", flag.ExitOnError)
		adminRefreshConfigTokenFlag = adminRefreshConfigFlags.String("token", "REQUIRED", "")

		catalogFlags = flag.NewFlagSet("catalog", flag.ContinueOnError)

		catalogRefreshFlags           = flag.NewFlagSet("refresh", flag.ExitOnError)
		catalogRefreshCatalogNameFlag = catalogRefreshFlags.String("catalog-name", "REQUIRED", "Name of catalog")
		catalogRefreshTokenFlag       = catalogRefreshFlags.String("token", "REQUIRED", "")

		catalogRefreshAllFlags     = flag.NewFlagSet("refresh-all", flag.ExitOnError)
		catalogRefreshAllTokenFlag = catalogRefreshAllFlags.String("token", "REQUIRED", "")

		catalogCatalogErrorFlags           = flag.NewFlagSet("catalog-error", flag.ExitOnError)
		catalogCatalogErrorCatalogNameFlag = catalogCatalogErrorFlags.String("catalog-name", "REQUIRED", "Name of catalog")
		catalogCatalogErrorTokenFlag       = catalogCatalogErrorFlags.String("token", "REQUIRED", "")

		categoryFlags = flag.NewFlagSet("category", flag.ContinueOnError)

		categoryListFlags = flag.NewFlagSet("list", flag.ExitOnError)

		ratingFlags = flag.NewFlagSet("rating", flag.ContinueOnError)

		ratingGetFlags     = flag.NewFlagSet("get", flag.ExitOnError)
		ratingGetIDFlag    = ratingGetFlags.String("id", "REQUIRED", "ID of a resource")
		ratingGetTokenFlag = ratingGetFlags.String("token", "REQUIRED", "")

		ratingUpdateFlags     = flag.NewFlagSet("update", flag.ExitOnError)
		ratingUpdateBodyFlag  = ratingUpdateFlags.String("body", "REQUIRED", "")
		ratingUpdateIDFlag    = ratingUpdateFlags.String("id", "REQUIRED", "ID of a resource")
		ratingUpdateTokenFlag = ratingUpdateFlags.String("token", "REQUIRED", "")

		resourceFlags = flag.NewFlagSet("resource", flag.ContinueOnError)

		resourceQueryFlags          = flag.NewFlagSet("query", flag.ExitOnError)
		resourceQueryNameFlag       = resourceQueryFlags.String("name", "", "")
		resourceQueryCatalogsFlag   = resourceQueryFlags.String("catalogs", "", "")
		resourceQueryCategoriesFlag = resourceQueryFlags.String("categories", "", "")
		resourceQueryKindsFlag      = resourceQueryFlags.String("kinds", "", "")
		resourceQueryTagsFlag       = resourceQueryFlags.String("tags", "", "")
		resourceQueryPlatformsFlag  = resourceQueryFlags.String("platforms", "", "")
		resourceQueryLimitFlag      = resourceQueryFlags.String("limit", "1000", "")
		resourceQueryMatchFlag      = resourceQueryFlags.String("match", "contains", "")

		resourceListFlags = flag.NewFlagSet("list", flag.ExitOnError)

		resourceVersionsByIDFlags  = flag.NewFlagSet("versions-by-id", flag.ExitOnError)
		resourceVersionsByIDIDFlag = resourceVersionsByIDFlags.String("id", "REQUIRED", "ID of a resource")

		resourceByCatalogKindNameVersionFlags       = flag.NewFlagSet("by-catalog-kind-name-version", flag.ExitOnError)
		resourceByCatalogKindNameVersionCatalogFlag = resourceByCatalogKindNameVersionFlags.String("catalog", "REQUIRED", "name of catalog")
		resourceByCatalogKindNameVersionKindFlag    = resourceByCatalogKindNameVersionFlags.String("kind", "REQUIRED", "kind of resource")
		resourceByCatalogKindNameVersionNameFlag    = resourceByCatalogKindNameVersionFlags.String("name", "REQUIRED", "name of resource")
		resourceByCatalogKindNameVersionVersionFlag = resourceByCatalogKindNameVersionFlags.String("version", "REQUIRED", "version of resource")

		resourceByVersionIDFlags         = flag.NewFlagSet("by-version-id", flag.ExitOnError)
		resourceByVersionIDVersionIDFlag = resourceByVersionIDFlags.String("version-id", "REQUIRED", "Version ID of a resource's version")

		resourceByCatalogKindNameFlags                = flag.NewFlagSet("by-catalog-kind-name", flag.ExitOnError)
		resourceByCatalogKindNameCatalogFlag          = resourceByCatalogKindNameFlags.String("catalog", "REQUIRED", "name of catalog")
		resourceByCatalogKindNameKindFlag             = resourceByCatalogKindNameFlags.String("kind", "REQUIRED", "kind of resource")
		resourceByCatalogKindNameNameFlag             = resourceByCatalogKindNameFlags.String("name", "REQUIRED", "Name of resource")
		resourceByCatalogKindNamePipelinesversionFlag = resourceByCatalogKindNameFlags.String("pipelinesversion", "", "")

		resourceByIDFlags  = flag.NewFlagSet("by-id", flag.ExitOnError)
		resourceByIDIDFlag = resourceByIDFlags.String("id", "REQUIRED", "ID of a resource")

		statusFlags = flag.NewFlagSet("status", flag.ContinueOnError)

		statusStatusFlags = flag.NewFlagSet("status", flag.ExitOnError)
	)
	adminFlags.Usage = adminUsage
	adminUpdateAgentFlags.Usage = adminUpdateAgentUsage
	adminRefreshConfigFlags.Usage = adminRefreshConfigUsage

	catalogFlags.Usage = catalogUsage
	catalogRefreshFlags.Usage = catalogRefreshUsage
	catalogRefreshAllFlags.Usage = catalogRefreshAllUsage
	catalogCatalogErrorFlags.Usage = catalogCatalogErrorUsage

	categoryFlags.Usage = categoryUsage
	categoryListFlags.Usage = categoryListUsage

	ratingFlags.Usage = ratingUsage
	ratingGetFlags.Usage = ratingGetUsage
	ratingUpdateFlags.Usage = ratingUpdateUsage

	resourceFlags.Usage = resourceUsage
	resourceQueryFlags.Usage = resourceQueryUsage
	resourceListFlags.Usage = resourceListUsage
	resourceVersionsByIDFlags.Usage = resourceVersionsByIDUsage
	resourceByCatalogKindNameVersionFlags.Usage = resourceByCatalogKindNameVersionUsage
	resourceByVersionIDFlags.Usage = resourceByVersionIDUsage
	resourceByCatalogKindNameFlags.Usage = resourceByCatalogKindNameUsage
	resourceByIDFlags.Usage = resourceByIDUsage

	statusFlags.Usage = statusUsage
	statusStatusFlags.Usage = statusStatusUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "admin":
			svcf = adminFlags
		case "catalog":
			svcf = catalogFlags
		case "category":
			svcf = categoryFlags
		case "rating":
			svcf = ratingFlags
		case "resource":
			svcf = resourceFlags
		case "status":
			svcf = statusFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "admin":
			switch epn {
			case "update-agent":
				epf = adminUpdateAgentFlags

			case "refresh-config":
				epf = adminRefreshConfigFlags

			}

		case "catalog":
			switch epn {
			case "refresh":
				epf = catalogRefreshFlags

			case "refresh-all":
				epf = catalogRefreshAllFlags

			case "catalog-error":
				epf = catalogCatalogErrorFlags

			}

		case "category":
			switch epn {
			case "list":
				epf = categoryListFlags

			}

		case "rating":
			switch epn {
			case "get":
				epf = ratingGetFlags

			case "update":
				epf = ratingUpdateFlags

			}

		case "resource":
			switch epn {
			case "query":
				epf = resourceQueryFlags

			case "list":
				epf = resourceListFlags

			case "versions-by-id":
				epf = resourceVersionsByIDFlags

			case "by-catalog-kind-name-version":
				epf = resourceByCatalogKindNameVersionFlags

			case "by-version-id":
				epf = resourceByVersionIDFlags

			case "by-catalog-kind-name":
				epf = resourceByCatalogKindNameFlags

			case "by-id":
				epf = resourceByIDFlags

			}

		case "status":
			switch epn {
			case "status":
				epf = statusStatusFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "admin":
			c := adminc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "update-agent":
				endpoint = c.UpdateAgent()
				data, err = adminc.BuildUpdateAgentPayload(*adminUpdateAgentBodyFlag, *adminUpdateAgentTokenFlag)
			case "refresh-config":
				endpoint = c.RefreshConfig()
				data, err = adminc.BuildRefreshConfigPayload(*adminRefreshConfigTokenFlag)
			}
		case "catalog":
			c := catalogc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "refresh":
				endpoint = c.Refresh()
				data, err = catalogc.BuildRefreshPayload(*catalogRefreshCatalogNameFlag, *catalogRefreshTokenFlag)
			case "refresh-all":
				endpoint = c.RefreshAll()
				data, err = catalogc.BuildRefreshAllPayload(*catalogRefreshAllTokenFlag)
			case "catalog-error":
				endpoint = c.CatalogError()
				data, err = catalogc.BuildCatalogErrorPayload(*catalogCatalogErrorCatalogNameFlag, *catalogCatalogErrorTokenFlag)
			}
		case "category":
			c := categoryc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "list":
				endpoint = c.List()
			}
		case "rating":
			c := ratingc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get":
				endpoint = c.Get()
				data, err = ratingc.BuildGetPayload(*ratingGetIDFlag, *ratingGetTokenFlag)
			case "update":
				endpoint = c.Update()
				data, err = ratingc.BuildUpdatePayload(*ratingUpdateBodyFlag, *ratingUpdateIDFlag, *ratingUpdateTokenFlag)
			}
		case "resource":
			c := resourcec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "query":
				endpoint = c.Query()
				data, err = resourcec.BuildQueryPayload(*resourceQueryNameFlag, *resourceQueryCatalogsFlag, *resourceQueryCategoriesFlag, *resourceQueryKindsFlag, *resourceQueryTagsFlag, *resourceQueryPlatformsFlag, *resourceQueryLimitFlag, *resourceQueryMatchFlag)
			case "list":
				endpoint = c.List()
			case "versions-by-id":
				endpoint = c.VersionsByID()
				data, err = resourcec.BuildVersionsByIDPayload(*resourceVersionsByIDIDFlag)
			case "by-catalog-kind-name-version":
				endpoint = c.ByCatalogKindNameVersion()
				data, err = resourcec.BuildByCatalogKindNameVersionPayload(*resourceByCatalogKindNameVersionCatalogFlag, *resourceByCatalogKindNameVersionKindFlag, *resourceByCatalogKindNameVersionNameFlag, *resourceByCatalogKindNameVersionVersionFlag)
			case "by-version-id":
				endpoint = c.ByVersionID()
				data, err = resourcec.BuildByVersionIDPayload(*resourceByVersionIDVersionIDFlag)
			case "by-catalog-kind-name":
				endpoint = c.ByCatalogKindName()
				data, err = resourcec.BuildByCatalogKindNamePayload(*resourceByCatalogKindNameCatalogFlag, *resourceByCatalogKindNameKindFlag, *resourceByCatalogKindNameNameFlag, *resourceByCatalogKindNamePipelinesversionFlag)
			case "by-id":
				endpoint = c.ByID()
				data, err = resourcec.BuildByIDPayload(*resourceByIDIDFlag)
			}
		case "status":
			c := statusc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "status":
				endpoint = c.Status()
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// adminUsage displays the usage of the admin command and its subcommands.
func adminUsage() {
	fmt.Fprintf(os.Stderr, `Admin service
Usage:
    %[1]s [globalflags] admin COMMAND [flags]

COMMAND:
    update-agent: Create or Update an agent user with required scopes
    refresh-config: Refresh the changes in config file

Additional help:
    %[1]s admin COMMAND --help
`, os.Args[0])
}
func adminUpdateAgentUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] admin update-agent -body JSON -token STRING

Create or Update an agent user with required scopes
    -body JSON: 
    -token STRING: 

Example:
    %[1]s admin update-agent --body '{
      "name": "abc",
      "scopes": [
         "catalog-refresh",
         "agent:create"
      ]
   }' --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"
`, os.Args[0])
}

func adminRefreshConfigUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] admin refresh-config -token STRING

Refresh the changes in config file
    -token STRING: 

Example:
    %[1]s admin refresh-config --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODAzMDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJhdGluZzpyZWFkIiwicmF0aW5nOndyaXRlIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.6pDmziSKkoSqI1f0rc4-AqVdcfY0Q8wA-tSLzdTCLgM"
`, os.Args[0])
}

// catalogUsage displays the usage of the catalog command and its subcommands.
func catalogUsage() {
	fmt.Fprintf(os.Stderr, `The Catalog Service exposes endpoints to interact with catalogs
Usage:
    %[1]s [globalflags] catalog COMMAND [flags]

COMMAND:
    refresh: Refresh a Catalog by it's name
    refresh-all: Refresh all catalogs
    catalog-error: List all errors occurred refreshing a catalog

Additional help:
    %[1]s catalog COMMAND --help
`, os.Args[0])
}
func catalogRefreshUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] catalog refresh -catalog-name STRING -token STRING

Refresh a Catalog by it's name
    -catalog-name STRING: Name of catalog
    -token STRING: 

Example:
    %[1]s catalog refresh --catalog-name "tekton" --token "Ad dicta."
`, os.Args[0])
}

func catalogRefreshAllUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] catalog refresh-all -token STRING

Refresh all catalogs
    -token STRING: 

Example:
    %[1]s catalog refresh-all --token "Tempore laboriosam placeat eveniet perspiciatis ut ut."
`, os.Args[0])
}

func catalogCatalogErrorUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] catalog catalog-error -catalog-name STRING -token STRING

List all errors occurred refreshing a catalog
    -catalog-name STRING: Name of catalog
    -token STRING: 

Example:
    %[1]s catalog catalog-error --catalog-name "tekton" --token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODM2MDAsImlhdCI6MTU3Nzg4MDAwMCwiaWQiOjExLCJpc3MiOiJUZWt0b24gSHViIiwic2NvcGVzIjpbInJlZnJlc2g6dG9rZW4iXSwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ.4RdUk5ttHdDiymurlZ_f7Uy5Pas3Lq9w04BjKQKRiCE"
`, os.Args[0])
}

// categoryUsage displays the usage of the category command and its subcommands.
func categoryUsage() {
	fmt.Fprintf(os.Stderr, `The category service provides details about category
Usage:
    %[1]s [globalflags] category COMMAND [flags]

COMMAND:
    list: List all categories along with their tags sorted by name

Additional help:
    %[1]s category COMMAND --help
`, os.Args[0])
}
func categoryListUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] category list

List all categories along with their tags sorted by name

Example:
    %[1]s category list
`, os.Args[0])
}

// ratingUsage displays the usage of the rating command and its subcommands.
func ratingUsage() {
	fmt.Fprintf(os.Stderr, `The rating service exposes endpoints to read and write user's rating for resources
Usage:
    %[1]s [globalflags] rating COMMAND [flags]

COMMAND:
    get: Find user's rating for a resource
    update: Update user's rating for a resource

Additional help:
    %[1]s rating COMMAND --help
`, os.Args[0])
}
func ratingGetUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] rating get -id UINT -token STRING

Find user's rating for a resource
    -id UINT: ID of a resource
    -token STRING: 

Example:
    %[1]s rating get --id 1154472414850186965 --token "Dolorem voluptatem cum."
`, os.Args[0])
}

func ratingUpdateUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] rating update -body JSON -id UINT -token STRING

Update user's rating for a resource
    -body JSON: 
    -id UINT: ID of a resource
    -token STRING: 

Example:
    %[1]s rating update --body '{
      "rating": 2
   }' --id 15596292303976604081 --token "Numquam et ea itaque nam rerum aut."
`, os.Args[0])
}

// resourceUsage displays the usage of the resource command and its subcommands.
func resourceUsage() {
	fmt.Fprintf(os.Stderr, `The resource service provides details about all kind of resources
Usage:
    %[1]s [globalflags] resource COMMAND [flags]

COMMAND:
    query: Find resources by a combination of name, kind, catalog, categories, platforms and tags
    list: List all resources sorted by rating and name
    versions-by-id: Find all versions of a resource by its id
    by-catalog-kind-name-version: Find resource using name of catalog & name, kind and version of resource
    by-version-id: Find a resource using its version's id
    by-catalog-kind-name: Find resources using name of catalog, resource name and kind of resource
    by-id: Find a resource using it's id

Additional help:
    %[1]s resource COMMAND --help
`, os.Args[0])
}
func resourceQueryUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource query -name STRING -catalogs JSON -categories JSON -kinds JSON -tags JSON -platforms JSON -limit UINT -match STRING

Find resources by a combination of name, kind, catalog, categories, platforms and tags
    -name STRING: 
    -catalogs JSON: 
    -categories JSON: 
    -kinds JSON: 
    -tags JSON: 
    -platforms JSON: 
    -limit UINT: 
    -match STRING: 

Example:
    %[1]s resource query --name "buildah" --catalogs '[
      "tekton",
      "openshift"
   ]' --categories '[
      "build",
      "tools"
   ]' --kinds '[
      "task",
      "pipelines"
   ]' --tags '[
      "image",
      "build"
   ]' --platforms '[
      "linux/s390x",
      "linux/amd64"
   ]' --limit 100 --match "contains"
`, os.Args[0])
}

func resourceListUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource list

List all resources sorted by rating and name

Example:
    %[1]s resource list
`, os.Args[0])
}

func resourceVersionsByIDUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource versions-by-id -id UINT

Find all versions of a resource by its id
    -id UINT: ID of a resource

Example:
    %[1]s resource versions-by-id --id 1
`, os.Args[0])
}

func resourceByCatalogKindNameVersionUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource by-catalog-kind-name-version -catalog STRING -kind STRING -name STRING -version STRING

Find resource using name of catalog & name, kind and version of resource
    -catalog STRING: name of catalog
    -kind STRING: kind of resource
    -name STRING: name of resource
    -version STRING: version of resource

Example:
    %[1]s resource by-catalog-kind-name-version --catalog "tektoncd" --kind "task" --name "buildah" --version "0.1"
`, os.Args[0])
}

func resourceByVersionIDUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource by-version-id -version-id UINT

Find a resource using its version's id
    -version-id UINT: Version ID of a resource's version

Example:
    %[1]s resource by-version-id --version-id 1
`, os.Args[0])
}

func resourceByCatalogKindNameUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource by-catalog-kind-name -catalog STRING -kind STRING -name STRING -pipelinesversion STRING

Find resources using name of catalog, resource name and kind of resource
    -catalog STRING: name of catalog
    -kind STRING: kind of resource
    -name STRING: Name of resource
    -pipelinesversion STRING: 

Example:
    %[1]s resource by-catalog-kind-name --catalog "tektoncd" --kind "task" --name "buildah" --pipelinesversion "0.21.0"
`, os.Args[0])
}

func resourceByIDUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] resource by-id -id UINT

Find a resource using it's id
    -id UINT: ID of a resource

Example:
    %[1]s resource by-id --id 1
`, os.Args[0])
}

// statusUsage displays the usage of the status command and its subcommands.
func statusUsage() {
	fmt.Fprintf(os.Stderr, `Describes the status of each service
Usage:
    %[1]s [globalflags] status COMMAND [flags]

COMMAND:
    status: Return status of the services

Additional help:
    %[1]s status COMMAND --help
`, os.Args[0])
}
func statusStatusUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] status status

Return status of the services

Example:
    %[1]s status status
`, os.Args[0])
}
