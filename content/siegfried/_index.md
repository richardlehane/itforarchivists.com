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
### Windows 

  - download latest (v. {{% version %}}) binary: ({{% download win64 64bit %}})
  - copy to a location in your [system path](http://www.computerhope.com/issues/ch000549.htm)
  - run the `sf -update` command to download the latest signatures (got troubles? Try this [troubleshooting guide](https://github.com/richardlehane/siegfried/wiki/Getting-started#installing-the-latest-signature-file))
  - if you want to build your own signatures with `roy`, copy the latest signature data into a "siegfried" directory within your user home directory (e.g. c:\users\richardl\siegfried): {{% datadownload %}}
 
### Mac [Homebrew](http://brew.sh) (or [Homebrew on Linux](https://docs.brew.sh/Homebrew-on-Linux))

    brew install richardlehane/digipres/siegfried

(a fork of [mistydemeo/digipres/siegfried](https://github.com/mistydemeo/homebrew-digipres))

### Ubuntu/Debian (64 bit)

    curl -sL "http://keyserver.ubuntu.com/pks/lookup?op=get&search=0x20F802FE798E6857" | gpg --dearmor | sudo tee /usr/share/keyrings/siegfried-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/siegfried-archive-keyring.gpg] https://www.itforarchivists.com/ buster main" | sudo tee -a /etc/apt/sources.list.d/siegfried.list
    sudo apt-get update && sudo apt-get install siegfried

### FreeBSD

	pkg install siegfried

### Arch Linux

	git clone https://aur.archlinux.org/siegfried.git
	cd siegfried
	makepkg -si

---

## Usage 
### Identify files and directories

	sf file.ext // Identify a file
	sf *.ext // Identify groups of files using a glob pattern
	sf DIR // Identify all files in a directory and its subdirectories
	sf -nr DIR // Identify all files in a directory but don't recurse into subdirectories

### Save output

	sf PATH > my_results.yaml // Use a redirect (">") to save your results. PATH means file.ext, *.ext or DIR
	sf -csv PATH > my_results.csv // Get identification results in CSV format (default is YAML)
	sf -json PATH > my_results.json // Get identification results in JSON format

### Additional commands

	sf -z PATH // Scan within zip, tar, gzip, warc or arc files
	sf -hash sha1 PATH // Calculate md5, sha1, sha256, sha512, or crc hash
	sf -multi 32 PATH // Scan many files at once
	sf -setconf -multi 32 -hash md5 -csv // Save your preferred configuration
	sf -setconf -csv -conf csv.conf // Save (or load) named configurations with -conf

### Update your signature file

	sf -update

## User guide

Detailed information about installing siegfried, identifying file formats, as well as more advanced topics, is available on [the wiki](https://github.com/richardlehane/siegfried/wiki).

## Modify your signature file

The **roy** tool builds siegfried signature files. For help using this tool, see this [guide](https://github.com/richardlehane/siegfried/wiki/Building-a-signature-file-with-ROY).

## Examples

{{% asciicast ernm49loq5ofuj48ywlvg7xq6 %}}

{{% asciicast 39270 %}}

---

## Code, License, Issues

To view the source code and see the license details, go to the project page on [Github](https://github.com/richardlehane/siegfried). Please post any bugs or feature request to the [issues page](https://github.com/richardlehane/siegfried/issues).

## Announcements

Join the [Google Group](https://groups.google.com/d/forum/sf-roy) for updates, signature releases, and help.

