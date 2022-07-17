project := "github.com/stephenwilliams/semver-cli"

default: format test snapshot

format: format-imports
    go fmt ./...

format-imports: tools
    #!/usr/bin/env bash
    set -euo pipefail
    for file in $(find . -name '*.go'); do
        ./bin/goimports-reviser -file-path $file -project-name {{project}}
    done

test:
    go test ./... -coverpkg=./... -coverprofile cover.out
    go tool cover -func cover.out

tools:
    #!/usr/bin/env bash
    set -e pipefail
    TMP_DIR=$(mktemp -d)
    cd $TMP_DIR
    go mod init tmp > /dev/null 2>&1
    GOBIN="{{justfile_directory()}}/bin"
    go install github.com/incu6us/goimports-reviser/v2@v2.5.1
    rm -rf $$TMP_DIR

snapshot:
    goreleaser release --snapshot --skip-publish --rm-dist

release:
    #!/usr/bin/env bash
    set -euxo pipefail
    git tag "$(svu next)"
    git push --tags
    goreleaser --rm-dist
