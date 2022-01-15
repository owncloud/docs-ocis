# oCis Documentation

## Building the Docs

The oCis documentation is not built independently. Instead, it is built together with the [server documentation](https://github.com/owncloud/docs/). However, you can build a local copy of the oCis documentation to preview changes you are making.

Whenever a Pull Request of this repo gets merged, it automatically triggers a full docs build.

## General Notes

To make life easier, most of the content written in [docs](https://github.com/owncloud/docs#readme) applies also here. For ease of reading, the most important steps are documented here too. For more information see the link provided. Only a few topics of this repo are unique like the branching.

## Prepare Your Environment

To prepare your local environment, some steps need to be made:

1.) Have the [necessary prerequisites](https://github.com/owncloud/docs/blob/master/docs/build-the-docs.md#install-the-prerequisites) installed.

2.) Clone this repository and run
```
yarn install
```
to setup all necessary dependencies.

## Building the oCis Documentation

Run the following command to build the client documentation locally

```
yarn antora-local
```

## Previewing the Generated Docs

Assuming that there are no build errors, the next thing to do is to view the result in your browser. In case you have already installed a web server to access local pages, you need to configure a virtual host (or similar) which points to the directory `public/`, located in the root directory of this repository. This directory contains the generated documentation. Alternatively, use the simple web server `serve` bundled with the current package.json, just execute the following command to serve the documentation at [http://localhost:8080/ocis/](http://localhost:8080/ocis/):

```
yarn serve
```

## Target Branch and Backporting

See the the [following section](https://github.com/owncloud/docs#target-branch-and-backporting) as the same rules and notes apply.

## Branching Workflow

Please refer to the [Branching Workflow for the oCis Documentation](https://github.com/owncloud/docs-ocis/blob/master/docs/the-branching-workflow.md) or more information.

## Create a New Version Branch for the oCis Documentation

Please refer to [Create a New Version Branch for the oCis Documentation](https://github.com/owncloud/docs-ocis/blob/master/docs/new-version-branch.md) for more information.
