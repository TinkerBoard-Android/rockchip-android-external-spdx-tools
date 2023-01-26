// Package parser2v3 contains functions to read, load and parse
// SPDX tag-value files, version 2.3.
// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later
package parser2v3

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_3"
	"github.com/spdx/tools-golang/tvloader/reader"
)

// ParseTagValues takes a list of (tag, value) pairs, parses it and returns
// a pointer to a parsed SPDX Document.
func ParseTagValues(tvs []reader.TagValuePair) (*v2_3.Document, error) {
	parser := tvParser2_3{}
	for _, tv := range tvs {
		err := parser.parsePair2_3(tv.Tag, tv.Value)
		if err != nil {
			return nil, err
		}
	}
	if parser.file != nil && parser.file.FileSPDXIdentifier == nullSpdxElementId2_3 {
		return nil, fmt.Errorf("file with FileName %s does not have SPDX identifier", parser.file.FileName)
	}
	if parser.pkg != nil && parser.pkg.PackageSPDXIdentifier == nullSpdxElementId2_3 {
		return nil, fmt.Errorf("package with PackageName %s does not have SPDX identifier", parser.pkg.PackageName)
	}
	return parser.doc, nil
}

func (parser *tvParser2_3) parsePair2_3(tag string, value string) error {
	switch parser.st {
	case psStart2_3:
		return parser.parsePairFromStart2_3(tag, value)
	case psCreationInfo2_3:
		return parser.parsePairFromCreationInfo2_3(tag, value)
	case psPackage2_3:
		return parser.parsePairFromPackage2_3(tag, value)
	case psFile2_3:
		return parser.parsePairFromFile2_3(tag, value)
	case psSnippet2_3:
		return parser.parsePairFromSnippet2_3(tag, value)
	case psOtherLicense2_3:
		return parser.parsePairFromOtherLicense2_3(tag, value)
	case psReview2_3:
		return parser.parsePairFromReview2_3(tag, value)
	default:
		return fmt.Errorf("parser state %v not recognized when parsing (%s, %s)", parser.st, tag, value)
	}
}

func (parser *tvParser2_3) parsePairFromStart2_3(tag string, value string) error {
	// fail if not in Start parser state
	if parser.st != psStart2_3 {
		return fmt.Errorf("got invalid state %v in parsePairFromStart2_3", parser.st)
	}

	// create an SPDX Document data struct if we don't have one already
	if parser.doc == nil {
		parser.doc = &v2_3.Document{ExternalDocumentReferences: []v2_3.ExternalDocumentRef{}}
	}

	switch tag {
	case "DocumentComment":
		parser.doc.DocumentComment = value
	case "SPDXVersion":
		parser.doc.SPDXVersion = value
	case "DataLicense":
		parser.doc.DataLicense = value
	case "SPDXID":
		eID, err := extractElementID(value)
		if err != nil {
			return err
		}
		parser.doc.SPDXIdentifier = eID
	case "DocumentName":
		parser.doc.DocumentName = value
	case "DocumentNamespace":
		parser.doc.DocumentNamespace = value
	case "ExternalDocumentRef":
		documentRefID, uri, alg, checksum, err := extractExternalDocumentReference(value)
		if err != nil {
			return err
		}
		edr := v2_3.ExternalDocumentRef{
			DocumentRefID: documentRefID,
			URI:           uri,
			Checksum:      common.Checksum{Algorithm: common.ChecksumAlgorithm(alg), Value: checksum},
		}
		parser.doc.ExternalDocumentReferences = append(parser.doc.ExternalDocumentReferences, edr)
	default:
		// move to Creation Info parser state
		parser.st = psCreationInfo2_3
		return parser.parsePairFromCreationInfo2_3(tag, value)
	}

	return nil
}
