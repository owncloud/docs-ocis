# Create a New Version Branch for Infinite Scale

Note that at the moment **no branching** is necessary for the ocis repo as only the `master` branch is used that appears a `next` in the documentation.

<!--
When doing a new release for Infinite Scale like `1.x`, a new version branch must be created based on `master`. It is necessary to do this in four steps. Please set the new and former version numbers accordingly

**Step 1: Create and configure the new `1.x` branch**

1.  Create a new `1.x` branch based on latest `origin/master`
2.  Copy the `.drone.star` file from the _former_ `1.x-1` branch
    (it contains the correct branch specific setup rules and replaces the current one coming from master)
3.  In `.drone.star` set `latest_version` to `1.x` (on top in section `def main(ctx)`)
4.  In `site.yml` adjust all `-version` keys according the new and former releases
    (in section `asciidoc.attributes`)
5.  In `antora.yml` change the version from `next` to `1.x`
6.  Run a build by entering `yarn antora-local`. No errors should occur
7.  Commit the changes and push the new `1.x` branch. **DO NOT CREATE A PR!**

**Step 2: Configure the master branch to use the new `1.x` branch**

9.  Create a new `changes_necessary_for_1.x` branch based on latest `origin/master`
10.  In `.drone.star` set `latest_version` to `1.x` (on top in section `def main(ctx)`)
11. In `site.yml` in section `asciidoc.attributes`, adjust all `-version` keys related to this repo according the new and former releases. Note if those attributes exist in other content sources, they must be set to the identical value to create consistent test builds.
12. No changes in `antora.yml` but check if the version is set to `next`
13. Run a build by entering `yarn antora-local`. No errors should occur
14. Commit changes and push it
15. Create a Pull Request. When CI is green, all is done correctly. Merge the PR to master.

**Step 3: Set the correct Branding build branches in the docs repo**

16. In `site.yml` of [docs](https://github.com/owncloud/docs/blob/master/site.yml) adjust the last **two** branches at `url: https://github.com/owncloud/docs-client-branding.git` accordingly
    (in section `content.sources.url.branches`)
17. In `site.yml` of [docs](https://github.com/owncloud/docs/blob/master/site.yml) adjust all `-version` keys in section `attributes` related to this repo according the new and former releases.

**Step 4: Protection and Renaming**

18. Go to the settings of the this repository and change the protection of the branch list (Settings > Branches) so that the `1.x` branch gets protected and the `1.x-2` branch is no longer protected.
19. Rename the `1.x-2` branch to `x_archived_1.x-2`

**Text Suggestion for Step 2**

The following text is a copy/paste suggestion for the PR in step 2, replace the branch numbers accordingly:
```
These are the changes necessary to finalize the creation of the 1.x branch.

The 1.x branch is already pushed and prepared and is included in the branch protection rules.

When 1.x (Infinite Scale) is finally out, the 1.x-2 branch can be archived,
see step 4 in https://github.com/owncloud/docs-client-branding/blob/master/docs/new-version-branch.md

Note, that the 1.x branch in this repo is already created, but the `latest` pointer on the web
will be set to it automatically when the tag in Branding is set. This means, that in the docs homepage,
`latest` will point to 1.x-1 until the tag in Branding Clients is set accordingly. When merging this PR,
1.x-2 will be dropped from the web but is available via pdf as usual.

Note, this PR must be merged before the 1.x tag in the Branding Clients repo is set to avoid a 404 for `latest`.

Note that a PR in docs must be made to announce the 1.x branch. The docs PR must be merged AFTER this PR is merged to avoid a CI error in docs.

Before merging this PR, we should take care that 1.x-2 has all changes necessary merged as post
merging the 1.x-2 pdf is fixed.

@michaelstingl @jesmrec fyi

@mmattel @EParzefall @phil-davis
post merging this, we need to backport all relevant changes to 1.x
```

-->
