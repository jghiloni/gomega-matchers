package matchers_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/jghiloni/gomega-matchers"
)

var _ = Describe("Zipfile", func() {
	It("Detects a string file as a zip file", func() {
		Expect("testdata/test.zip").To(BeAZipFile())
	})

	It("Detects a byte slice as a zip file", func() {
		b, err := os.ReadFile("testdata/test.zip")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(b).To(BeAZipFile())
	})

	It("Detects a file object as a zip file", func() {
		fp, err := os.Open("testdata/test.zip")
		Expect(err).ShouldNot(HaveOccurred())
		defer fp.Close()
		Expect(fp).To(BeAZipFile())
	})
})
