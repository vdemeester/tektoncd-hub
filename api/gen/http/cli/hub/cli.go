// Code generated by goa v3.2.0, DO NOT EDIT.
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

	categoryc "github.com/tektoncd/hub/api/gen/http/category/client"
	resourcec "github.com/tektoncd/hub/api/gen/http/resource/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `category list
resource (query|list|versions-by-id|by-type-name-version|by-version-id)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` category list` + "\n" +
		os.Args[0] + ` resource query --name "Occaecati officia inventore adipisci." --type "pipeline" --limit 364697698949290656` + "\n" +
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
) (goa.Endpoint, interface{}, error) {
	var (
		categoryFlags = flag.NewFlagSet("category", flag.ContinueOnError)

		categoryListFlags = flag.NewFlagSet("list", flag.ExitOnError)

		resourceFlags = flag.NewFlagSet("resource", flag.ContinueOnError)

		resourceQueryFlags     = flag.NewFlagSet("query", flag.ExitOnError)
		resourceQueryNameFlag  = resourceQueryFlags.String("name", "", "")
		resourceQueryTypeFlag  = resourceQueryFlags.String("type", "", "")
		resourceQueryLimitFlag = resourceQueryFlags.String("limit", "100", "")

		resourceListFlags     = flag.NewFlagSet("list", flag.ExitOnError)
		resourceListLimitFlag = resourceListFlags.String("limit", "100", "")

		resourceVersionsByIDFlags  = flag.NewFlagSet("versions-by-id", flag.ExitOnError)
		resourceVersionsByIDIDFlag = resourceVersionsByIDFlags.String("id", "REQUIRED", "ID of a resource")

		resourceByTypeNameVersionFlags       = flag.NewFlagSet("by-type-name-version", flag.ExitOnError)
		resourceByTypeNameVersionTypeFlag    = resourceByTypeNameVersionFlags.String("type", "REQUIRED", "type of resource")
		resourceByTypeNameVersionNameFlag    = resourceByTypeNameVersionFlags.String("name", "REQUIRED", "name of resource")
		resourceByTypeNameVersionVersionFlag = resourceByTypeNameVersionFlags.String("version", "REQUIRED", "version of resource")

		resourceByVersionIDFlags         = flag.NewFlagSet("by-version-id", flag.ExitOnError)
		resourceByVersionIDVersionIDFlag = resourceByVersionIDFlags.String("version-id", "REQUIRED", "Version ID of a resource's version")
	)
	categoryFlags.Usage = categoryUsage
	categoryListFlags.Usage = categoryListUsage

	resourceFlags.Usage = resourceUsage
	resourceQueryFlags.Usage = resourceQueryUsage
	resourceListFlags.Usage = resourceListUsage
	resourceVersionsByIDFlags.Usage = resourceVersionsByIDUsage
	resourceByTypeNameVersionFlags.Usage = resourceByTypeNameVersionUsage
	resourceByVersionIDFlags.Usage = resourceByVersionIDUsage

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
		case "category":
			svcf = categoryFlags
		case "resource":
			svcf = resourceFlags
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
		case "category":
			switch epn {
			case "list":
				epf = categoryListFlags

			}

		case "resource":
			switch epn {
			case "query":
				epf = resourceQueryFlags

			case "list":
				epf = resourceListFlags

			case "versions-by-id":
				epf = resourceVersionsByIDFlags

			case "by-type-name-version":
				epf = resourceByTypeNameVersionFlags

			case "by-version-id":
				epf = resourceByVersionIDFlags

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
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "category":
			c := categoryc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "list":
				endpoint = c.List()
				data = nil
			}
		case "resource":
			c := resourcec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "query":
				endpoint = c.Query()
				data, err = resourcec.BuildQueryPayload(*resourceQueryNameFlag, *resourceQueryTypeFlag, *resourceQueryLimitFlag)
			case "list":
				endpoint = c.List()
				data, err = resourcec.BuildListPayload(*resourceListLimitFlag)
			case "versions-by-id":
				endpoint = c.VersionsByID()
				data, err = resourcec.BuildVersionsByIDPayload(*resourceVersionsByIDIDFlag)
			case "by-type-name-version":
				endpoint = c.ByTypeNameVersion()
				data, err = resourcec.BuildByTypeNameVersionPayload(*resourceByTypeNameVersionTypeFlag, *resourceByTypeNameVersionNameFlag, *resourceByTypeNameVersionVersionFlag)
			case "by-version-id":
				endpoint = c.ByVersionID()
				data, err = resourcec.BuildByVersionIDPayload(*resourceByVersionIDVersionIDFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// categoryUsage displays the usage of the category command and its subcommands.
func categoryUsage() {
	fmt.Fprintf(os.Stderr, `The category service provides details about category
Usage:
    %s [globalflags] category COMMAND [flags]

COMMAND:
    list: List all categories along with their tags sorted by name

Additional help:
    %s category COMMAND --help
`, os.Args[0], os.Args[0])
}
func categoryListUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] category list

List all categories along with their tags sorted by name

Example:
    `+os.Args[0]+` category list
`, os.Args[0])
}

// resourceUsage displays the usage of the resource command and its subcommands.
func resourceUsage() {
	fmt.Fprintf(os.Stderr, `The resource service provides details about all type of resources
Usage:
    %s [globalflags] resource COMMAND [flags]

COMMAND:
    query: Find resources by a combination of name, type
    list: List all resources sorted by rating and name
    versions-by-id: Find all versions of a resource by its id
    by-type-name-version: Find resource using name, type and version of resource
    by-version-id: Find a resource using its version's id

Additional help:
    %s resource COMMAND --help
`, os.Args[0], os.Args[0])
}
func resourceQueryUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource query -name STRING -type STRING -limit UINT

Find resources by a combination of name, type
    -name STRING: 
    -type STRING: 
    -limit UINT: 

Example:
    `+os.Args[0]+` resource query --name "Occaecati officia inventore adipisci." --type "pipeline" --limit 364697698949290656
`, os.Args[0])
}

func resourceListUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource list -limit UINT

List all resources sorted by rating and name
    -limit UINT: 

Example:
    `+os.Args[0]+` resource list --limit 11281538076509796713
`, os.Args[0])
}

func resourceVersionsByIDUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource versions-by-id -id UINT

Find all versions of a resource by its id
    -id UINT: ID of a resource

Example:
    `+os.Args[0]+` resource versions-by-id --id 10108643860615476272
`, os.Args[0])
}

func resourceByTypeNameVersionUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-type-name-version -type STRING -name STRING -version STRING

Find resource using name, type and version of resource
    -type STRING: type of resource
    -name STRING: name of resource
    -version STRING: version of resource

Example:
    `+os.Args[0]+` resource by-type-name-version --type "task" --name "Omnis quas deserunt nostrum assumenda." --version "Occaecati voluptas assumenda."
`, os.Args[0])
}

func resourceByVersionIDUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] resource by-version-id -version-id UINT

Find a resource using its version's id
    -version-id UINT: Version ID of a resource's version

Example:
    `+os.Args[0]+` resource by-version-id --version-id 11745309465362025378
`, os.Args[0])
}