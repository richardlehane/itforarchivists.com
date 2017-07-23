+++
title = "siegfried"
description = "siegfried is a signature-based file format identification tool"
date = "2017-07-20"
+++
# Siegfried

Siegfried is a signature-based file format identification tool.

It implements:

  - the National Archives UK's [PRONOM](http://www.nationalarchives.gov.uk/pronom) file format signatures 
  - freedesktop.org's [MIME-info](https://freedesktop.org/wiki/Software/shared-mime-info) file format signatures
  - the Library of Congress's [FDD](http://www.digitalpreservation.gov/formats/fdd/descriptions.shtml) file format signatures (**beta**).

---

## Install
### Windows {{% fontawesome windows %}}

  - download latest (v. {{% version %}}) binary: ({{% download win64 64bit %}} | {{% download win32 32bit %}})
  - copy to a location in your [system path](http://www.computerhope.com/issues/ch000549.htm)
  - run the `sf -update` command to download the latest signatures (got troubles? Try this [troubleshooting guide](https://github.com/richardlehane/siegfried/wiki/Getting-started#installing-the-latest-signature-file))
 
### Mac Homebrew (or [Linuxbrew](http://brew.sh/linuxbrew/)) {{% fontawesome apple %}}

- `brew install richardlehane/digipres/siegfried` (a fork of [mistydemeo/digipres/siegfried](https://github.com/mistydemeo/homebrew-digipres))

### Ubuntu/Debian (64 bit) {{% fontawesome linux %}}

- `wget -qO - https://bintray.com/user/downloadSubjectPublicKey?username=bintray | sudo apt-key add -`
- `echo "deb http://dl.bintray.com/siegfried/debian wheezy main" | sudo tee -a /etc/apt/sources.list`
- `sudo apt-get update && sudo apt-get install siegfried`

### FreeBSD {{% fontawesome mercury %}}

- `pkg install siegfried`

### Arch Linux {{% fontawesome linux %}}

- `git clone https://aur.archlinux.org/siegfried.git`
- `cd siegfried`
- `makepkg -si`

---

## Usage 
### Identify a file

`sf file.ext`

### Identify all files in a directory and its subdirectories

`sf DIR`

Identify all files in a directory but not subdirectories

`sf -nr DIR`

### Use a redirect (">") to save your results

`sf file.ext or DIR > my_results.yaml`

### Get identification results in CSV format (default is YAML)

`sf -csv file.ext or DIR > my_results.csv<`

### Get identification results in JSON format

`sf -json file.ext or DIR > my_results.json`

### Scan within zip, tar, gzip, warc or arc files

`sf -z file.ext or DIR`

### Calculate md5, sha1, sha256, sha512, or crc hash

`sf -hash sha1 file.ext or DIR`

### Update your signature file

`sf -update`

## Additional commands

For more detail on using **sf**, see this [guide](https://github.com/richardlehane/siegfried/wiki/Identifying-file-formats).

## Modify your signature file

The **roy** tool builds siegfried signature files. For help using this tool, see this [guide](https://github.com/richardlehane/siegfried/wiki/Building-a-signature-file-with-ROY).

## Examples

{{% asciicast ernm49loq5ofuj48ywlvg7xq6 %}}

{{% asciicast 39270 %}}

## User Guide

Detailed information about installing siegfried, identifying file formats, as well as more advanced topics, is available on [the wiki](https://github.com/richardlehane/siegfried/wiki).

---

## Code, License, Issues

To view the source code and see the license details, go to the project page on [Github](https://github.com/richardlehane/siegfried). Please post any bugs or feature request to the [issues page](https://github.com/richardlehane/siegfried/issues).

## Announcements

Join the [Google Group](href="https://groups.google.com/d/forum/sf-roy") for updates, signature releases, and help.

