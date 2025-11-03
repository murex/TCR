# Releasing a new TCR version

We use [GoReleaser](https://goreleaser.com/) for releasing new TCR versions.

## Versioning Rules

TCR release versions comply with [Semantic Versioning rules](https://semver.org/).

## Release Branch

All TCR releases are published on GitHub's `main` branch.

## Release Preparation

- [ ] Update dependencies and run sanity checks:
    ```bash
    make prepare
    ```
- [ ] Commit all changes on the `main` branch
- [ ] Push the changes to GitHub and [wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Run the release preparation script:
    ```bash
    ./tools/scripts/prepare-release.sh vX.Y.Z
    ```
  - This script will automatically update version files, commit changes, and create the release tag
- [ ] Verify that everything is ready for GoReleaser:
    ```bash
    make snapshot
    ```

## Releasing

The creation of the new release is triggered by pushing the newly created release tag to GitHub repository

- [ ] Push the release tag:
    ```bash
    git push origin vX.Y.Z
    ```
- [ ] [Wait until all GitHub Actions are green](https://github.com/murex/TCR/actions)
- [ ] Open [TCR Release page](https://github.com/murex/TCR/releases) and verify that the new release is there
- [ ] Edit the release notes document, and insert a `Summary` section at the top, listing the main changes included in
  this release. You may take a look at previous release notes if unsure what should go in there.
