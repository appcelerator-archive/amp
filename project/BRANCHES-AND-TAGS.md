Branches and tags
=================

Note: details of the release process for AMP are documented in the
[RELEASE-CHECKLIST](RELEASE-CHECKLIST.md).

# Branches

The appcelerator/amp repository should normally have only three living branches at all time, including
the regular `master` branch:

## `docs` branch

The `docs` branch supports documentation updates between product releases. This branch allow us to
decouple documentation releases from product releases.

## "release" branch

The "release" branch is not actually called `release`, but instead indicates the target release milestone
(ex: `1.6.x`).

The release branch contains the last *released* version of the code for the project (ex: `1.5.3`)
and any commits that have been cherry-picked from `master`.

See [Release Checklist](RELEASE-CHECKLIST.md) for details.

# Tags

Any public release of a compiled binary, with the logical exception of nightly builds, should have
a corresponding tag in the repository.

The general format of a tag is `vX.Y.Z[-suffix[N]]`:

- All of `X`, `Y`, `Z` must be specified (example: `v1.0.0`)
- First release candidate for version `1.8.0` should be tagged `v1.8.0-rc1`
- Second alpha release of a product should be tagged `v1.0.0-alpha1`
