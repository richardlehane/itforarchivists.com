---
title: "Charts"
date: 2017-09-14
categories: ["siegfried"]
---

In my recent updates to this site I've added a new "Chart your results" tool on the [siegfried](/siegfried) page (in the right hand panel under "Try Siegfried"). This tool produces single page reports like this: [/siegfried/results/ea1zaj](https://www.itforarchivists.com/siegfried/results/ea1zaj).

Before covering this tool in detail let's recap some of the existing ways you can already analyse your results.

## Other ways of charting and analysing your results

### Command-line charting

I appreciate that not everyone is a command-line junkie, but the way I inspect results is just to use sf's -log flag. If you do `sf -log chart` (or `-log c`) you can make simple format charts:

![sf chart](/img/sf-chart.png)

*(In these examples I add "o" to my log options to direct logging output to STDOUT... otherwise you'll see it in STDERR).*

A chart can be a starting point for deeper analysis e.g. inspecting lists of files of a particular format: 

![sf log fmt/61](/img/sf-fmt61.png)

You can also inspect lists of unknowns with `-log u` and warnings with `-log w`. 

Rather than re-run the format identification job with every step, you can pair these commands with the `-replay` flag to run them against a pre-generated results file instead. I cover this workflow in detail in the [siegfried wiki](https://github.com/richardlehane/siegfried/wiki/Identifying-file-formats#replaying-a-scan-from-results-files).

### Standalone tools

It would be remiss of me not to mention the two great standalone tools that Ross Spencer and Tim Walsh have written for analysing your results: [DROID-SF sqlite analysis](https://github.com/exponential-decay/droid-siegfried-sqlite-analysis-engine) and [Brunnhilde](https://github.com/timothyryanwalsh/brunnhilde).

These tools both do a lot more than simple chart generation. E.g. DROID-SF can create a "Rogues Gallery" of all your problematic files. Brunnhilde has a GUI, does virus scanning, and can also run bulk_extractor against your files. I'd definitely encourage you to check both of these tools out!

## Chart your results

If your needs are a little bit simpler, and you just want a chart, then my new "Chart your results" tool might be a good fit.

To try this tool, go to the [siegfried](/siegfried) page and upload a results file in the "Chart my results" form in the right-hand panel.

Let's run through some of its features:

  - it can handle siegfried, DROID and fido results files
  - it gives you a single page, interactive report
  - pie charts for format IDs and MIME-types
  - highlights unknowns, errors, warnings, multiple IDs and duplicates (if you have checksums)
  - you can drill-down on all of those features and formats to generate lists that you can export as CSV or to your clipboard
  - supports multiple identifiers (a siegfried feature): each identifier is a separate single page report.

Probably the distinguishing feature of this tool is that you can easily share your analysis with colleagues, or with the digital preservation community broadly, by "publishing" your results. This gives you a permanent URL (like [https://www.itforarchivists.com/siegfried/results/ea1zaj](https://www.itforarchivists.com/siegfried/results/ea1zaj)) and stores your results on the site. Prior to publication you can opt to "redact" your filenames if they contain sensitive information. I've added a [privacy](/about/privacy) section to this site to address some of the privacy questions raised by this feature in a little more detail.

That's it, please use it, and if you like it tweet your results!
