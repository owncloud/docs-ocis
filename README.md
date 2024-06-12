# Infinite Scale Documentation

**Table of Contents**

* [Building the Infinite Scale Docs](#building-the-infinite-scale-docs)
* [General Notes](#general-notes)
* [Generating the Documentation](#generating-the-documentation)
* [Important Notes](#important-notes)
* [Target Branch and Backporting](#target-branch-and-backporting)
* [Branching Workflow](#branching-workflow)
* [Create a New Version Branch for Infinite Scale](#create-a-new-version-branch-for-infinite-scale)

## Building the Infinite Scale Docs

The Infinite Scale documentation is not built independently. Instead, it is built together with the [main documentation](https://github.com/owncloud/docs/). However, you can build a local copy of the Infinite Scale documentation to preview changes you are making.

Whenever a Pull Request of this repo gets merged, it automatically triggers a full docs build.

## General Notes

To make life easier, most of the content written in [docs](https://github.com/owncloud/docs#readme) applies also here. For ease of reading, the most important steps are documented here too. For more information see the link provided. Only a few topics of this repo are unique like the branching.

## Generating the Documentation

See the [Generating the Documentation](https://github.com/owncloud/docs#generating-the-documentation) in the docs repo for more details as it applies to all documentation repositories.

## Important Notes

* When the Infinite Scale dev team creates a new service and merges the code, you **must** add a new service page in the services folder using the service name as document name. If this is omitted, **ALL** new and pending doc (!!) PR's ()will error with `target of xref not found` because of missing reference targets. These references originate in the `env-vars-special-scope.adoc` document which uses sources from the `ocis` repo containing automatically generated content where the referenced target is missing in the admin docs.

## Target Branch and Backporting

See the [following section](https://github.com/owncloud/docs#target-branch-and-backporting) as the same rules and notes apply.

## Branching Workflow

Please refer to the [Branching Workflow for the Infinite Scale](./docs/the-branching-workflow.md) for more information.

## Create a New Version Branch for Infinite Scale

Please refer to [Create a New Version Branch for Infinite Scale](./docs/new-version-branch.md) for more information.
