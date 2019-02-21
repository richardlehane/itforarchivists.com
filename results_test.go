package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResults(t *testing.T) {
	res, err := results(strings.NewReader(testRes), "testdata")
	if err != nil {
		t.Fatal(err)
	}
	if !res.validate() {
		t.Fatal("results don't validate")
	}
	byt, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	res2 := &Results{}
	err = json.Unmarshal(byt, res2)
	if err != nil {
		t.Fatal(err)
	}
	if !res2.validate() {
		t.Fatal("results don't validate")
	}
	byt2, err := json.Marshal(res2)
	if err != nil {
		t.Fatal(err)
	}
	if len(byt) != len(byt2) {
		t.Fatalf("expecting results roundtrip to result in a match:\n%s\n\n---------\n\n%s", string(byt), string(byt2))
	}
}

var testRes = `---
siegfried   : 1.7.4
scandate    : 2017-07-31T14:08:02+10:00
signature   : default.sig
created     : 2017-07-14T15:56:40+10:00
identifiers : 
  - name    : 'pronom'
    details : 'DROID_SignatureFile_V90.xml; container-signature-20170330.xml'
  - name    : 'loc'
    details : 'fddXML.zip (2017-06-10, DROID_SignatureFile_V90.xml, container-signature-20170330.xml)'
---
filename : 'procedures\Briefing note template for Executive Director.docx'
filesize : 223627
modified : 2016-11-16T12:24:58+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 378, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Briefing note template.docx'
filesize : 55417
modified : 2015-06-09T15:36:46+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 378, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Contractors Brochure  2 WSRC .pdf'
filesize : 448021
modified : 2015-03-05T12:27:33+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/19'
    format  : 'Acrobat PDF 1.5 - Portable Document Format'
    version : '1.5'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at [[[0 8]] [[448016 5]]]'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000030'
    format  : 'PDF (Portable Document Format) Family'
    full    : 'PDF (Portable Document Format) Family'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at 0, 4'
    warning : 
---
filename : 'procedures\DA and ACM workflows V0.2.docx'
filesize : 20752
modified : 2016-01-21T10:53:05+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 327, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\DFS - Code of Conduct - Sep 13.pdf'
filesize : 1182010
modified : 2015-01-22T15:44:32+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/19'
    format  : 'Acrobat PDF 1.5 - Portable Document Format'
    version : '1.5'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at [[[0 8]] [[1182003 5]]]'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000030'
    format  : 'PDF (Portable Document Format) Family'
    full    : 'PDF (Portable Document Format) Family'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at 0, 4'
    warning : 
---
filename : 'procedures\Direct Report Template.docx'
filesize : 77439
modified : 2015-02-25T17:16:14+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 379, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\File Note.dot'
filesize : 63488
modified : 2011-06-17T10:54:48+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'x-fmt/45'
    format  : 'Microsoft Word Document Template'
    version : '97-2003'
    mime    : 
    basis   : 'extension match dot; container name CompObj with byte match at 77, 20; name WordDocument with byte match at 10, 1'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000132'
    format  : 'Macromedia Flash FLA Project File Format'
    full    : 'Macromedia Flash FLA Project File Format'
    mime    : 
    basis   : 'byte match at 0, 24'
    warning : 'extension mismatch'
---
filename : 'procedures\Fin 11 SRA Intangible Assets Policy v4.pdf'
filesize : 421503
modified : 2015-02-19T14:48:47+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/19'
    format  : 'Acrobat PDF 1.5 - Portable Document Format'
    version : '1.5'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at [[[0 8]] [[421496 5]]]'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000030'
    format  : 'PDF (Portable Document Format) Family'
    full    : 'PDF (Portable Document Format) Family'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at 0, 4'
    warning : 
---
filename : 'procedures\IDG PLAN AUGUST 2016.docx'
filesize : 32145
modified : 2016-08-25T16:51:07+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 327, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Letterhead.docx'
filesize : 65916
modified : 2015-11-03T13:42:44+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 377, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Petty Cash Claim Form 2014 11 28.xlsx'
filesize : 361693
modified : 2015-01-19T13:09:42+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/214'
    format  : 'Microsoft Excel for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    basis   : 'extension match xlsx; container name [Content_Types].xml with byte match at 642, 88 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000398, fdd000401'
---
filename : 'procedures\RL Direct Report 12 May.docx'
filesize : 77987
modified : 2015-05-11T15:23:51+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 379, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\RL Direct Report 31 Mar 15.docx'
filesize : 77246
modified : 2015-03-30T11:51:28+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 379, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\SRA Vendor Relations Protocol and Declaration of Interest.docx'
filesize : 67471
modified : 2015-09-01T11:29:54+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 379, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Smartsheet Training v2.0 for Delegates.docx'
filesize : 126076
modified : 2015-03-03T13:34:28+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 503, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\Smartsheet Training v2.0.docx'
filesize : 632434
modified : 2015-03-03T13:34:28+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 503, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\State Records LBCP 2015 V1.6.docx'
filesize : 1681817
modified : 2015-03-26T15:11:34+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/412'
    format  : 'Microsoft Word for Windows'
    version : '2007 onwards'
    mime    : 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    basis   : 'extension match docx; container name [Content_Types].xml with byte match at 481, 94 (signature 1/3)'
    warning : 
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match; possibilities based on extension are fdd000395, fdd000397, fdd000400'
---
filename : 'procedures\State Records Purchasing Card Guidelines - signed copy (18 Aug 2016).pdf'
filesize : 11155087
modified : 2016-09-13T17:37:04+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/17'
    format  : 'Acrobat PDF 1.3 - Portable Document Format'
    version : '1.3'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at [[[0 8]] [[11155080 5]]]'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000030'
    format  : 'PDF (Portable Document Format) Family'
    full    : 'PDF (Portable Document Format) Family'
    mime    : 'application/pdf'
    basis   : 'extension match pdf; byte match at 0, 4'
    warning : 
---
filename : 'procedures\StateArchives&Records_LOGO_Black_CMYK.JPG'
filesize : 670992
modified : 2016-11-16T12:25:03+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/41'
    format  : 'Raw JPEG Stream'
    version : 
    mime    : 'image/jpeg'
    basis   : 'extension match jpg; byte match at [[[0 3]] [[670990 2]]] (signature 1/2)'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000017'
    format  : 'JPEG Image Encoding Family'
    full    : 'JPEG Image Encoding Family'
    mime    : 
    basis   : 'byte match at 0, 2'
    warning : 
---
filename : 'procedures\StateArchives&Records_LOGO_Two colour_CMYK.JPG'
filesize : 831508
modified : 2016-11-16T12:25:01+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/41'
    format  : 'Raw JPEG Stream'
    version : 
    mime    : 'image/jpeg'
    basis   : 'extension match jpg; byte match at [[[0 3]] [[831506 2]]] (signature 1/2)'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000017'
    format  : 'JPEG Image Encoding Family'
    full    : 'JPEG Image Encoding Family'
    mime    : 
    basis   : 'byte match at 0, 2'
    warning : 
---
filename : 'procedures\TRIM files.md'
filesize : 171
modified : 2014-12-17T14:22:41+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'x-fmt/111'
    format  : 'Plain Text File'
    version : 
    mime    : 'text/plain'
    basis   : 'text match ASCII'
    warning : 'match on text only; extension mismatch'
  - ns      : 'loc'
    id      : 'UNKNOWN'
    format  : 
    full    : 
    mime    : 
    basis   : 
    warning : 'no match'
---
filename : 'procedures\Team meeting agenda.txt'
filesize : 346
modified : 2015-03-16T10:46:14+11:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'x-fmt/111'
    format  : 'Plain Text File'
    version : 
    mime    : 'text/plain'
    basis   : 'extension match txt; text match ASCII'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000284'
    format  : 'ESRI ArcInfo Coverage'
    full    : 'ESRI ArcInfo Coverage'
    mime    : 
    basis   : 'extension match txt'
    warning : 'match on extension only'
---
filename : 'procedures\selection committee reports.doc'
filesize : 74752
modified : 2015-06-09T16:15:30+10:00
errors   : 
matches  :
  - ns      : 'pronom'
    id      : 'fmt/40'
    format  : 'Microsoft Word Document'
    version : '97-2003'
    mime    : 'application/msword'
    basis   : 'extension match doc; container name CompObj with byte match at 78, 20; name WordDocument with name only'
    warning : 
  - ns      : 'loc'
    id      : 'fdd000132'
    format  : 'Macromedia Flash FLA Project File Format'
    full    : 'Macromedia Flash FLA Project File Format'
    mime    : 
    basis   : 'byte match at 0, 24'
    warning : 'extension mismatch'
`
