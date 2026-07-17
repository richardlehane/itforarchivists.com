# IT for archivists website (https://itforarchivists.com)

Some features lost in this migration including:

- automatic benchmarks
- share your results
- sets tool.

I will migrate legacy content for the benchmarks and results over the coming weeks (see below).

These features never saw much uptake & removing them allows me to greatly simplify maintenance burden by replacing the old appengine application with a fully static (using Zine ssg) site.

As part of the migration, I've moved "Try Siegfried" to the WASM build which makes it much faster, more private/secure, and more fully featured (checksums, directory scanning, -z flag).

## Remaining migration tasks

- old blog: move to /attic/blog
- benchmarks: move to /attic
- results: move to /attic & include redirects or cf functions to preserve links
- extract results generation code and include it in sf
- extract update.json generation code and create a script for new site
- migrate /scripts to ease updates
- decommission old server and bucket

## Release workflow

1. Run a script to copy over latest signatures to assets/

2. Update the ziggy custom metadata (site or sf content level) with version and release date.

3. Copy deb into the site and:

- run reprepro, 
- then copy the pool and dists folder to assets.

4. Update these files in Cloudflare R2:

- wasmexec
- 1_ll deb