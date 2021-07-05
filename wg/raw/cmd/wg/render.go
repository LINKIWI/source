package main

import (
	"fmt"
	"strconv"
	"strings"

	"wg/internal/cli"
	"wg/pkg/webgrep"
)

var (
	// renderers is a map of known renderer aliases to their corresponding implementations.
	renderers = map[string]resultsRenderer{
		"table":   &tableResultsRenderer{},
		"stacked": &stackedResultsRenderer{},
	}
)

// resultsRenderer provides an interface for logic that renders structured search results into a
// string sequence printed to the terminal.
type resultsRenderer interface {
	// RenderCodeSearchResults renders code search results.
	RenderCodeSearchResults(code []webgrep.CodeSearchResult) (string, error)

	// RenderFileSearchResults renders file search results.
	RenderFileSearchResults(files []webgrep.FileSearchResult) (string, error)
}

// tableResultsRenderer is an implementation of the default renderer style, which prints results
// as a table with columns for the version, repository, file path, and result contents (code line
// or file name). It is intended to be accessible to humans while remaining machine-parsable by
// standard Unix tools.
type tableResultsRenderer struct{}

// RenderCodeSearchResults renders a table of version, repository, file path, line number, and line
// contents.
func (d *tableResultsRenderer) RenderCodeSearchResults(code []webgrep.CodeSearchResult) (string, error) {
	table := cli.NewTable()

	for _, result := range code {
		for _, line := range result.Lines {
			source := line.Line
			if len(line.Bounds) == 2 {
				// Apply highlighting to the matching portion of the source
				// line, if applicable
				source = strings.Join([]string{
					line.Line[:line.Bounds[0]],
					cli.Highlight(line.Line[line.Bounds[0]:line.Bounds[1]]),
					line.Line[line.Bounds[1]:],
				}, "")
			}

			row := []string{
				result.Version,
				result.Repo,
				result.Path,
				strconv.Itoa(line.Number),
				fmt.Sprintf("|%s", source),
			}

			if err := table.Add(row); err != nil {
				return "", err
			}
		}
	}

	return table.String(), nil
}

// RenderFileSearchResults renders a table of version, repository, and file path.
func (d *tableResultsRenderer) RenderFileSearchResults(files []webgrep.FileSearchResult) (string, error) {
	table := cli.NewTable()

	for _, result := range files {
		path := strings.Join([]string{
			result.Path[:result.Bounds[0]],
			cli.Highlight(result.Path[result.Bounds[0]:result.Bounds[1]]),
			result.Path[result.Bounds[1]:],
		}, "")

		row := []string{
			result.Version,
			result.Repo,
			path,
		}

		if err := table.Add(row); err != nil {
			return "", err
		}
	}

	return table.String(), nil
}

// stackedResultsRenderer is an implementation of a renderer style that organizes results as stacked
// blocks grouped by repository and/or matching file path name. It resembles the results layout of
// the webgrep UI and is generally intended for consumption by humans due to its machine-unfriendly
// format.
type stackedResultsRenderer struct{}

// RenderCodeSearchResults renders blocks of code, each with a header indicating the source
// repository and file path.
func (s *stackedResultsRenderer) RenderCodeSearchResults(code []webgrep.CodeSearchResult) (string, error) {
	var lines []string

	for idx, result := range code {
		block := cli.NewTable()
		header := fmt.Sprintf("%s > %s", result.Repo, result.Path)

		lines = append(lines, cli.Bold(header))
		lines = append(lines, cli.Bold(strings.Repeat("-", len(header))))

		for _, line := range result.Lines {
			source := line.Line
			if len(line.Bounds) == 2 {
				source = strings.Join([]string{
					line.Line[:line.Bounds[0]],
					cli.Highlight(line.Line[line.Bounds[0]:line.Bounds[1]]),
					line.Line[line.Bounds[1]:],
				}, "")
			}

			block.Add([]string{
				strconv.Itoa(line.Number),
				fmt.Sprintf("|%s", source),
			})
		}

		lines = append(lines, block.String())

		if idx < len(code)-1 {
			lines = append(lines, "")
		}
	}

	return strings.Join(lines, "\n"), nil
}

// RenderFileSearchResults renders blocks of file paths, each with a header indicating the source
// repository by which it is grouped.
func (s *stackedResultsRenderer) RenderFileSearchResults(files []webgrep.FileSearchResult) (string, error) {
	var lines []string
	var idx int

	repoFiles := make(map[string][]webgrep.FileSearchResult)

	for _, result := range files {
		repoFiles[result.Repo] = append(repoFiles[result.Repo], result)
	}

	for repo, result := range repoFiles {
		lines = append(lines, cli.Bold(repo))
		lines = append(lines, cli.Bold(strings.Repeat("-", len(repo))))

		for _, file := range result {
			path := strings.Join([]string{
				file.Path[:file.Bounds[0]],
				cli.Highlight(file.Path[file.Bounds[0]:file.Bounds[1]]),
				file.Path[file.Bounds[1]:],
			}, "")

			lines = append(lines, path)
		}

		if idx < len(repoFiles)-1 {
			lines = append(lines, "")
		}

		idx++
	}

	return strings.Join(lines, "\n"), nil
}
