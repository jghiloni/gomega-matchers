package matchers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"

	slices2 "github.com/jghiloni/go-commonutils/v2/slices"
	"github.com/onsi/gomega/types"
)

func BeAZipFile() types.GomegaMatcher {
	return &zipfileMatcher{}
}

func BeAZipFileWithExpectedEntries(entries ...string) types.GomegaMatcher {
	return &zipfileMatcher{
		expectedEntries: entries,
	}
}

type zipfileMatcher struct {
	expectedEntries []string
	zr              *zip.Reader
}

// Match implements types.GomegaMatcher.
func (z *zipfileMatcher) Match(actual any) (success bool, err error) {
	if err = z.isActualZip(actual); err != nil {
		return false, nil
	}

	if len(z.expectedEntries) > 0 {
		entries := slices2.Map(z.zr.File, func(f *zip.File) string {
			return f.Name
		})

		for _, e := range z.expectedEntries {
			if !slices.Contains(entries, e) {
				return false, nil
			}
		}
	}

	return true, nil
}

// FailureMessage implements types.GomegaMatcher.
func (z *zipfileMatcher) FailureMessage(actual any) (message string) {
	return z.failureMessage(actual, false)
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (z *zipfileMatcher) NegatedFailureMessage(actual any) (message string) {
	return z.failureMessage(actual, true)
}

func (z *zipfileMatcher) isActualZip(actual any) error {
	switch t := actual.(type) {
	case string:
		return z.isFilenameZip(t)
	case []byte:
		return z.isByteSliceZip(t)
	case *os.File:
		return z.isFileZip(t)
	default:
		return fmt.Errorf("%T must be string, []byte, or *os.File", actual)
	}
}

func (z *zipfileMatcher) isFilenameZip(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	return z.isFileZip(fp)
}

func (z *zipfileMatcher) isFileZip(fp *os.File) error {
	st, err := os.Stat(fp.Name())
	if err != nil {
		return err
	}

	z.zr, err = zip.NewReader(fp, st.Size())
	return err
}

func (z *zipfileMatcher) isByteSliceZip(data []byte) error {
	r := bytes.NewReader(data)
	s := int64(len(data))

	var err error
	z.zr, err = zip.NewReader(r, s)

	return err
}

func (z *zipfileMatcher) failureMessage(actual any, negative bool) (message string) {
	message = fmt.Sprintf("Expected\n\t%#v\nto", actual)
	if negative {
		message = message + " not"
	}

	message = message + " represent a zip file"
	if len(z.expectedEntries) > 0 {
		message = fmt.Sprintf("%s that contained entries [%s]", message, strings.Join(z.expectedEntries, ", "))
	}

	return message
}
