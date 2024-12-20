+++
title = "Privacy"
description = "how the itforarchivists.com site handles your data"
date = "2022-05-15"
+++

The aim of this page is to describe how your private information is used by this site. If I've missed anything glaring, or you have any questions, or you've inadvertently shared information on this site that you'd like taken down, please [contact me](mailto:richard@itforarchivists.com).

## General

  - the site is completely HTTPS - thanks [Let's Encrypt](https://www.letsencrypt.org)!
  - the site is hosted by [Google appengine](https://cloud.google.com/appengine) &, when I refer to server-side processing below, that's what I mean
  - the site's code is all open and published on [Github](https://www.github.com/richardlehane/itforarchivists.com).

## Use and storage of your data

There are three services on this site where you can upload data. All three services are available on the [siegfried](/siegfried) page (in the right hand panel).

### Sets tool

The sets service accepts very little input from users (just format IDs and sets). Shouldn't be anything private there. The sets service processes that input server-side. Nothing is stored after that process completes.

### Try siegfried

The "Try Siegfried" service allows users to upload files for identification. This is a server-side process so files are transmitted to Google appengine. The only thing the site does to the files is run the siegfried identification routine. Nothing is stored after that process completes.

### Charts your results

The charts service allows users to upload file identification results (raw results) for charting and analysis. Generating these charts is a server-side process so raw results are transmitted to Google appengine. The only thing that the site uses the raw results for is to generate a normalised view of the results that powers the charts. Nothing is stored after this process completes.

When a charts page has been generated, users have a further option to "publish" their results. This is the *only* point at which the site stores user data. The data stored is the normalised view of the results (and not the raw results file). Users have an option to "redact" their results before or during publication. This process runs the golang standard library function [.Ext](https://www.golang.org/pkg/path/filepath#Ext) to replace file names in the normalised results with their extension only. This process is irreversible and if applied then only redacted filenames (i.e. file extensions) are stored by the site.

The URLs generated by the publication process are pseudo-random and the site doesn't currently offer any way to discover other users' results (so they do need to be explicitly shared e.g. by tweeting them). But please don't rely on this practical obscurity for privacy! Your results may very well be indexed by search engines or other harvesters. Furthermore I may in future package published results as a downloadable corpus for analysis.

If you inadvertently publish your results and need them taken down, please [contact me](mailto:richard@itforarchivists.com). 
