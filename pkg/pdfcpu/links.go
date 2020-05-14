/*
Copyright 2018 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pdfcpu

import (
	"strings"
)

// RelativizeFileLinks converts absolute links to file:/// urls to relative urls.
// Helpful when generating a bunch of interconnected PDF documents from HTML pages.
func RelativizeFileLinks(ctx *Context, selectedPages IntSet) error {
	rootDict, err := ctx.Catalog()
	if err != nil {
		return err
	}
	pages, _ := rootDict.Find("Pages")

	dp, err := ctx.DereferenceDict(pages)
	if err != nil {
		return err
	}

	a := dp.ArrayEntry("Kids")
	for pageIndex, pageRef := range a {

		if selectedPages != nil && !selectedPages[pageIndex+1] {
			continue
		}

		page, err := ctx.DereferenceDict(pageRef)
		if err != nil {
			return err
		}

		annots := page.ArrayEntry("Annots")
		for _, annotRef := range annots {

			annot, err := ctx.DereferenceDict(annotRef)
			if err != nil {
				return err
			}

			a, found := annot.Find("A")
			if found {
				uri, found := a.(Dict).Find("URI")
				uriStr := uri.(StringLiteral).Value()
				if found {
					if strings.HasPrefix(uriStr, "file:///") {
						uriStr = strings.TrimPrefix(uriStr, "file:///")
						a.(Dict).Update("URI", StringLiteral(uriStr))
					}
				}
			}
		}
	}

	return nil
}
