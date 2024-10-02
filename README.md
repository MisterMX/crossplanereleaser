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

An example `.crossplanereleaser.yaml` config can look like this

```yaml
project_name: my-project
dist: dist/my-project

builds:
  - id: composition-package
    dir: package/compositions
    examples: examples
  - id: function-package
    dir: package/function
    examples: "IGNORE"
    # Use a prebuilt image tar that contains the function binary, i.e. by Ko
    runtime_image_tar: dist/function-base-image.tar

pushes:
  - build: composition-package
    image_templates:
      - "my-registry.com/my-project/package-compositions:{{ .Tag }}"
      - "my-registry.com/my-project/package-compositions:{{ .FullCommit }}"
  - id: function # Used to manually filter for images to be pushed
    build: function-package
    image_templates:
      - "my-registry.com/my-project/package-function:{{ .Tag }}"
      - "my-registry.com/my-project/package-function:{{ .FullCommit }}"
```

# Requirements

Crossplanereleaser does not deal packages by itself but instead uses external
tools for that which need to be available in your system:

* `git` to generate package meta information.
* [`crank`](https://github.com/crossplane/crossplane/tree/master/cmd/crank)
  for package building
* [`crane`](https://github.com/google/go-containerregistry/tree/main/cmd/crane)
  for image publishing.
