# Collection Log Source Scraper

This is a tool written in Go that (responsibly) scrapes the [OSRS wiki](https://oldschool.runescape.wiki) to populate the data required for the Collection Log Source website. The wiki contains all the required info, just not structured in the way we want it structured, hence this tool.

The intention is for this to run once a week, to update the data for the website with each OSRS update. Effort has been taken to try to be a good citizen; a descriptive user agent is set, and there is a 1 second delay between each request to the Wiki.

Note it is pretty brittle; I've not really put any effort into safe error handling so it is pretty prone to panicking if, for example, the Wiki slightly changes its formatting.

This tool is not associated with or endorsed by the OSRS Wiki, Jagex, or Old School RuneScape.

I claim no ownership of the data that this tool scrapes.

## Usage

| flag            | action                                       |
|-----------------|----------------------------------------------|
| `--mock`        | use the stored mock data                     |
| `--pretty`      | pretty print the result                      |
| `--debug`       | additional debug logging to stderr           |
| `--no-log`      | do not write loglines to stderr              |
| `-o={filename}` | output the results to `{filename}`           |
| `--tee`         | if `-o` has been used, also output to stdout |
| `-t={num}`      | truncate the result set to `{num}`           |

### Examples

use the mock data, pretty print the output. **Note** that the mock data assumes every item is an abyssal whip, but is useful for testing the tool.
```sh
go run . -- --mock --pretty
```

use real data, output results to a file. **Note** that takes a long time due to 1s waits between requests so we don't hammer the wiki
```sh
go run . > ./results.json
```

or
```sh
go run . -- -o=./results.json
```

If you use Jetbrains Goland/IntelliJ, there are 3 run configurations delivered with this project, 

- one that uses mock output ([`Scraper (mock, pretty)`](./.idea/runConfigurations/Scraper__mock__pretty_.xml))
- one that uses real data but truncates to 5 results ([`Scraper only 5`](./.idea/runConfigurations/Scraper_only_5.xml))
- one that uses real data and does not pretty print ([`Scraper`](./.idea/runConfigurations/Scraper.xml))

All of them output the results to the (git-ignored) file `results.json`.

## Testing

Currently, the only test validates that the output matches the schema.

```sh
go test ./...
```

## Output format

The output format follows [this JSON Schema](./schema.json).