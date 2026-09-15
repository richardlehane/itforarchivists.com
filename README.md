# IT for archivists website (https://itforarchivists.com)

New codebase for itforarchivists, replacing the old go appengine/hugo combo with Zine SSG.

Regretfully, some features lost in the migration, including:

- automatic benchmarks
- share your results
- sets tool.

*I will migrate legacy content for the benchmarks and results over the coming weeks (see below).*

Rationale: these features never saw much uptake and removing them allows me to greatly simplify site maintenance.

On the plus side: as part of the migration, I've moved "Try Siegfried" to the WASM build which makes it much faster, more private/secure, and more fully featured (checksums, directory scanning, -z flag).

## Remaining migration tasks

  - [X] old blog: move to /attic/blog
  - [X] benchmarks: move to /attic
  - [X] results: move to /attic & include redirects or cf functions to preserve links
  - [ ] extract results generation code and include it in sf
  - [ ] extract update.json generation code and create a script for new site
  - [ ] migrate /scripts to ease updates
  - [X] decommission old server and bucket

## Release workflow

1. Copy over latest signatures to assets/ (from Code for latest and ~/siegfried if build with roy10 (roy10 build and roy10 build -extend @archivematica archivematica))

2.  Go into scripts/update folder and run "go run run.go" 

3. Update the ziggy custom metadata (site or sf content level) with version and release date.

4. Copy deb into the site and:

- run reprepro, "reprepro -b repositories includedeb buster siegfried_1.X.X-1_amd64.deb"
- then copy the dists folder to assets.

5. Delete the "public" folder and run "zine release"

5. Update these files in Cloudflare R2:

- sf.wasm
- siegfried_1_ll.X.deb