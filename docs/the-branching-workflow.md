# The Branching Workflow

Only three branches are maintained at any one time; these are `master`, the current, and the former oCIS release series. Any change to the documentation is made in a branch based off of `master`. Once the branch's PR is approved and merged, the PR is backported to the branch for the **current** Desktop release and the **former** release but only if it applies to it.

When a new ownCloud major or minor oCis version is released, a new branch is created to track the changes for that release. The branch for the oldest release is frozen, taken off the active maintained branch list and is no longer maintained.
