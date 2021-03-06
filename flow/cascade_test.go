package flow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cascading YAML templates", func() {
	It("flows through multiple templates", func() {
		source := parseYAML(`
---
foo: (( merge ))
baz: 42
`)

		secondary := parseYAML(`
---
foo:
  bar: (( merge ))
  xyz: (( bar ))
`)

		stub := parseYAML(`
---
foo:
  bar: merged!
`)

		resolved := parseYAML(`
---
foo:
  bar: merged!
  xyz: merged!
baz: 42
`)

		Expect(source).To(CascadeAs(resolved, secondary, stub))
	})

	Context("with multiple mutually-exclusive templates", func() {
		It("flows through both", func() {
			source := parseYAML(`
---
foo: (( merge ))
baz: (( merge ))
`)

			secondary := parseYAML(`
---
foo:
  bar: (( merge ))
`)

			tertiary := parseYAML(`
---
baz:
  a: 1
  b: (( merge ))
`)

			stub := parseYAML(`
---
foo:
  bar: merged!
baz:
  b: 2
`)

			resolved := parseYAML(`
---
foo:
  bar: merged!
baz:
  a: 1
  b: 2
`)

			Expect(source).To(CascadeAs(resolved, secondary, tertiary, stub))
		})
	})

	Describe("node annotation propagation", func() {

		Context("referencing a merged field", func() {
			It("is not handled as merge expression", func() {
				source := parseYAML(`
---
alice: (( merge ))
bob: (( alice ))
`)
				stub := parseYAML(`
---
alice: alice
bob: bob
`)
				resolved := parseYAML(`
---
alice: alice
bob: bob
`)
				Expect(source).To(CascadeAs(resolved, stub))
			})
		})
	})

	Describe("merging lists with specified key", func() {

		Context("auto merge with key tag", func() {
			It("overrides matching key entries", func() {
				source := parseYAML(`
---
list:
  - key:address: a
    attr: b
  - address: c
    attr: d
`)
				stub := parseYAML(`
---
list:
  - address: c
    attr: stub
  - address: e
    attr: f
`)
				resolved := parseYAML(`
---
list:
  - address: a
    attr: b
  - address: c
    attr: stub
`)
				Expect(source).To(CascadeAs(resolved, stub))
			})

			It("overrides matching key entries with key tag", func() {
				source := parseYAML(`
---
list:
  - key:address: a
    attr: b
  - address: c
    attr: d
`)
				stub := parseYAML(`
---
list:
  - key:address: c
    attr: stub
  - address: e
    attr: f
`)
				resolved := parseYAML(`
---
list:
  - address: a
    attr: b
  - address: c
    attr: stub
`)
				Expect(source).To(CascadeAs(resolved, stub))
			})
		})
	})
})
