# Releasing a new TCR version

We use [GoReleaser](https://goreleaser.com/) for releasing new TCR versions.

## Versioning Rules

TCR release versions comply with [Semantic Versioning rules](https://semver.org/).

## Release Branch

All TCR releases are published on GitHub's `main` branch.

## Release Preparation

- [ ] Update dependencies and run sanity checks: `make prepare`
- [ ] Commit all changes on the `main` branch
- [ ] Push the changes to GitHub and [wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Create the release tag: `git tag -m "" -a vX.Y.Z`
- [ ] Verify that everything is ready for GoReleaser: `make snapshot`

## Releasing

The creation of the new release is triggered by pushing the newly created release tag to GitHub repository

- [ ] Push the release tag: `git push origin vX.Y.Z`
- [ ] [Wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Open [TCR Release page](https://github.com/murex/TCR/releases) and verify that the new release is there
- [ ] Edit the release notes document, and insert a `Summary` section at the top, listing the main changes included in
  this release. You may take a look at previous release notes if unsure what should go in there.
