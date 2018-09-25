---
title: "Continuous benchmarking"
date: 2018-07-30T20:00:00+10:00
categories: ["siegfried"]
---

The next siegfried release will be out shortly. I have been busy making changes to address two thorny issues: [verbose basis](https://github.com/richardlehane/siegfried/issues/111) and [missing results](https://github.com/richardlehane/siegfried/issues/112). I have some fixes in place but, when doing large scale testing against big sets of files, I noticed some performance and quality regressions. You can see these regressions on the new [develop benchmarks page](/siegfried/develop). There's also a new [benchmarks page](/siegfried/benchmarks) to measure siegfried against comparable file format identification tools (at this stage just [DROID](http://www.nationalarchives.gov.uk/information-management/manage-information/policy-process/digital-continuity/file-profiling-tool-droid/)). This post is the story of how these two pages came into being.

## Benchmarking pain

Up until this point I had always done large scale testing manually and it was a fiddly and annoying process: I would have to locate test corpora on whichever machine I'd last run a benchmark on and and copy to my current laptop; run the tests; do the comparison; interpret the results; look up how to do golang profiling again because I had forgot; run the profiler and check those results; etc. After going through all these steps, if I made changes to address issues, I'd have to repeat it all again in order to verify I'd fixed the issues. Obviously none of these results had much shelf life, as they depended on the vagaries of the machine I was running them on, how I had configured everything, and would be invalidated each time there was a new software or PRONOM release. I was also left with the uneasy feeling that I should be doing this kind of large scale testing much more regularly. Wouldn't it be nice to have some kind of continuous benchmarking process, just like the other automated tests and deployment workflows I run using [Travis CI](https://travis-ci.org/) and [appveyor](https://www.appveyor.com/)? And down that rabbit hole I went for the last few months... 

## The solution

Today, whenever I push code changes to siegfried's [develop branch](https://github.com/richardlehane/siegfried/tree/develop), or tag a new release on the master branch, the following happens:

  - siegfried's [Travis CI script](https://github.com/richardlehane/siegfried/blob/develop/.travis.yml) runs a little program called [provisioner](https://github.com/richardlehane/provisioner)
  - provisioner buys a machine for an hour or two from [packet.net](https://packet.net), feeding the new machine a [cloud init](https://cloud-init.io/) script that runs a series of install tasks (like downloading and installing siegfried)
  - the cloud init script concludes by starting another small program, [runner](https://github.com/richardlehane/runner), as a systemd service
  - runner downloads a list of jobs from the itforarchivists.com website (either the [develop jobs](/siegfried/jobs/develop) or the [benchmark jobs](/siegfried/jobs/bench)) and executes them
  - one of the early jobs is to use [rclone](https://rclone.org) to copy the test corpora over to the test machine from [backblaze.com](https://backblaze.com)
  - after running the jobs (and timing their duration), runner posts results back to the itforachivists.com website where they are stored and displayed on the [develop](/siegfried/develop) and [benchmarks](/siegfried/benchmarks) pages.

Benefits of this approach are:

  - it is completely transparent. Obviously as one of the tool makers there is potential for me to show some bias. But you can see exactly what machine has run the benchmark ([packet.net](https://packet.net) sells "baremetal" servers so you get a clear picture of the hardware used), what software has been installed and how it has been configured, and exactly what tasks have been run. 
  - it is cheap to run: Travis CI is free for open source projects, backblaze is super cheap, backblaze has partnered with packet.net for free data transport between their data centres, and the packet.net servers are competitively priced (particularly when you buy from the spot market)
  - it is routine. For me, this is the most important thing. I no longer have to go through benchmark pain and I can view the real world impacts of any changes I make to siegfried's code immediately after committing changes to github.

## Some reflections on the results

[DROID](http://www.nationalarchives.gov.uk/information-management/manage-information/policy-process/digital-continuity/file-profiling-tool-droid/) has got really fast recently! Kudos to the team at the National Archives for continuing to invest in the Aston Martin of format identification tools :). I'm particularly impressed by the DROID "no limit" (a -1 max bytes to scan setting in your DROID properties) results and wonder if it might make sense for future DROID releases to make that the default setting.

In order to make the most of SSD disks, you really need to use a `-multi` setting with sf to obtain good speeds. I'm making this easier in the next release of siegfried by introducting a config file where you can store your preferred settings (and not type them in each time you invoke the tool).

If you really care about speed, you can use `roy` to build signature files with built in "max bytes" limits. These test runs came out fastest in all categories. But this will impact the quality of your results.

Siegfried is a sprinter and wins on small corpora. DROID runs marathons and wins (in the no limit category) for the biggest corpus. This is possibly because of JVM start-up costs?

Speed isn't the only thing that matters. You also need to assess the quality of results and choose a tool that has the right affordances (e.g. do your staff prefer a GUI? What kind of reporting will you need? etc.). But speed is important, particularly for non-ingest workflows (e.g. consider re-scanning your repository each time there is a PRONOM update).

The tools differ in their outputs. This is because of differences in their matching engines (e.g. siegfried includes a text identifier), differences in their default settings (particularly that max byte setting), and differences in the way they choose to report results (e.g. if more than one result is returned just based on extension matches, then siegfried will return UNKNOWN with a descriptive warning indicating those possibilies; DROID and fido, on the other hand, will report multiple results).

Where's [fido](https://github.com/openpreserve/fido)? I had fido in early benchmarks but removed it because the files within the test corpora I've used cause fido to come to a grinding halt for some reason. I need to inspect the error messages and follow up. I hope to get fido back onto the scoreboard shortly!

I need more corpora. The corpora I've used reflect some use cases (e.g. scanning typical office type documents) but don't represent others (e.g. scanning large AV collections). Big audio and video files have caused problems for siegfried in the past and it would be great to include them in regular testing.

