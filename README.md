# XPReleaser

Like [Goreleaser](https://github.com/goreleaser/goreleaser) but for Crossplane packages.

# Usage

Build and push artifacts to a registry:

```bash
crossplanereleaser release
```

To just build artifacts:

```bash
crossplanereleaser build
```

# Requirements

Crossplanereleaser does not deal packages by itself but instead uses external
tools for that which need to be available in your system:

* `git` to generate package meta information.
* [`crank`](https://github.com/crossplane/crossplane/tree/master/cmd/crank)
  for package building
* [`crane`](https://github.com/google/go-containerregistry/tree/main/cmd/crane)
  for image publishing.
