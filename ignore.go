package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const doubleAst = "**"

type Pattern struct {
	patterns []string
}

// TODO: improve pattern matching algorithm to work in O(n)
// TODO: add support for subdirectory ignore file support

func ReadIgnoreFile(filename string) *Pattern {
	data, err := os.ReadFile(filename)
	p := new(Pattern)
	if err == nil {
		p.patterns = strings.Split(string(data), "\n")
		return p
	}
	return &Pattern{}
}

func (p *Pattern) Match(fpath string) bool {
	for _, pattern := range p.patterns {
		if matchSinglePattern(pattern, fpath) {
			return true
		}
	}
	return false
}

// https://github.com/zabawaba99/go-gitignore

// Match matches patterns in the same manner that gitignore does.
// Reference https://git-scm.com/docs/gitignore.
func matchSinglePattern(pattern, value string) bool {
	// A blank line matches no files, so it can serve as a separator for readability.
	if pattern == "" {
		return false
	}

	// A line starting with # serves as a comment. Put a backslash ("\") in front of the first hash for patterns that begin with a hash.
	if strings.HasPrefix(pattern, "#") {
		return false
	}

	// Trailing spaces are ignored unless they are quoted with backslash ("\"). (To be implemented)
	pattern = strings.TrimSuffix(pattern, " ")

	// An optional prefix "!" which negates the pattern; any matching file
	// excluded by a previous pattern will become included again. It is not
	// possible to re-include a file if a parent directory of that file is excluded.
	// Git doesn’t list excluded directories for performance reasons, so any patterns
	// on contained files have no effect, no matter where they are defined.
	// Put a backslash ("\") in front of the first "!" for patterns that begin
	// with a literal "!", for example, "\!important!.txt".
	negate := strings.HasPrefix(pattern, "!")
	if negate {
		pattern = strings.TrimPrefix(pattern, "!")
	}

	// If the pattern ends with a slash, it is removed for the purpose of the
	// following description, but it would only find a match with a directory.
	// In other words, foo/ will match a directory foo and paths underneath it,
	// but will not match a regular file or a symbolic link foo (this is consistent
	// with the way how pathspec works in general in Git).
	pattern = strings.TrimSuffix(pattern, string(os.PathSeparator))

	// Two consecutive asterisks ("**") in patterns matched
	// against full pathname may have special meaning:
	if strings.Contains(pattern, doubleAst) {
		result := evalDoubleAst(pattern, value)
		if negate {
			result = !result
		}
		return result
	}

	// If the pattern does not contain a slash /, Git treats it as a shell glob
	// pattern and checks for a match against the pathname relative to the location
	// of the .gitignore file (relative to the toplevel of the work tree if not from
	// a .gitignore file).
	if !strings.Contains(pattern, string(os.PathSeparator)) {
		m, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("ERROR %s\n", err)
			return false
		}

		var found bool
		for _, v := range m {
			if v == value {
				found = true
				break
			}
		}

		// Fixes bug that it can access sub directories and files even if parent folder
		// is ignored (By DelusionaL)
		if !found && strings.HasPrefix(value, pattern+string(os.PathSeparator)) {
			return true
		}

		if negate {
			return !found
		}
		return found
	}

	// Otherwise, Git treats the pattern as a shell glob suitable for consumption by
	// fnmatch(3) with the FNM_PATHNAME flag: wildcards in the pattern will not match
	// a / in the pathname. For example, "Documentation/*.html" matches
	// "Documentation/git.html" but not "Documentation/ppc/ppc.html" or
	// "tools/perf/Documentation/perf.html".

	// A leading slash matches the beginning of the pathname. For example, "/*.c" matches "cat-file.c" but not "mozilla-sha1/sha1.c".

	matched, err := path.Match(pattern, value)
	if err != nil {
		return false
	}

	// Fixes bug that it can access sub directories and files even if parent folder
	// is ignored (By DelusionaL)
	if !matched && strings.HasPrefix(value, pattern+string(os.PathSeparator)) {
		return true
	}

	if negate {
		return !matched
	}
	return matched

}

func evalDoubleAst(pattern, value string) bool {
	// A leading "**" followed by a slash means match in all directories.
	// For example, "**/foo" matches file or directory "foo" anywhere,
	// the same as pattern "foo". "**/foo/bar" matches file or directory
	// "bar" anywhere that is directly under directory "foo".
	if strings.HasPrefix(pattern, doubleAst) {
		pattern = strings.TrimPrefix(pattern, doubleAst)
		return strings.HasSuffix(value, pattern)
	}

	// A trailing "/**" matches everything inside. For example, "abc/**"
	// matches all files inside directory "abc", relative to the location
	// of the .gitignore file, with infinite depth.
	if strings.HasSuffix(pattern, doubleAst) {
		pattern = strings.TrimSuffix(pattern, doubleAst)
		return strings.HasPrefix(value, pattern)
	}

	// A slash followed by two consecutive asterisks then a slash matches
	// zero or more directories. For example, "a/**/b" matches "a/b",
	// /"a/x/b", "a/x/y/b" and so on.
	parts := strings.Split(pattern, doubleAst)
	for i, part := range parts {
		switch i {
		case 0:
			if !strings.HasPrefix(value, part) {
				return false
			}
		case len(parts) - 1: // last part
			part = strings.TrimPrefix(part, string(os.PathSeparator))
			return strings.HasSuffix(value, part)
		default:
			if !strings.Contains(value, part) {
				return false
			}
		}

		// trim evaluated text
		index := strings.Index(value, part) + len(part)
		value = value[index:]
	}

	// Other consecutive asterisks are considered invalid.
	return false
}
