# Note: These commands are for the development of bucketeer-find-code-refs.
# They are not intended for use by the end-users of this program.
build:
	go build ./cmd/...

init:
	pre-commit install

test: lint
	go test ./...

lint:
	pre-commit run -a --verbose golangci-lint

# Generate docs about GitHub Action inputs and updates README.md
github-action-docs:
	cd build/metadata/github-actions && npx action-docs -u --no-banner

# Strip debug informatino from production builds
BUILD_FLAGS = -ldflags="-s -w"

compile-macos-binary:
	GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} -o out/bucketeer-find-code-refs ./cmd/bucketeer-find-code-refs

compile-windows-binary:
	GOOS=windows GOARCH=amd64 go build ${BUILD_FLAGS} -o out/bucketeer-find-code-refs.exe ./cmd/bucketeer-find-code-refs

compile-linux-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/cmd/bucketeer-find-code-refs ./cmd/bucketeer-find-code-refs

compile-github-actions-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/github-actions/bucketeer-find-code-refs-github-action ./build/package/github-actions

compile-bitbucket-pipelines-binary:
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o build/package/bitbucket-pipelines/bucketeer-find-code-refs-bitbucket-pipeline ./build/package/bitbucket-pipelines

# Get the lines added to the most recent changelog update (minus the first 2 lines)
RELEASE_NOTES=<(GIT_EXTERNAL_DIFF='bash -c "diff --unchanged-line-format=\"\" $$2 $$5" || true' git log --ext-diff -1 --pretty= -p CHANGELOG.md)

echo-release-notes:
	@cat $(RELEASE_NOTES)

define publish_docker
	test $(1) || (echo "Please provide tag"; exit 1)
	docker build -t launchdarkly/$(3):$(1) build/package/$(4)
	docker push launchdarkly/$(3):$(1)
	# test $(2) && (echo "Not pushing latest tag for prerelease")
	test $(2) || docker tag launchdarkly/$(3):$(1) launchdarkly/$(3):latest
	test $(2) || docker push launchdarkly/$(3):latest
endef

# TODO: Remove all circleci publishing targets when we have a github owner token setup.
# Use ./ldrelease/publish-circleci.sh to publish to circleci orbs registry.
validate-circle-orb:
	test $(TAG) || (echo "Please provide tag"; exit 1)
	circleci orb validate build/package/circleci/orb.yml || (echo "Unable to validate orb"; exit 1)

publish-dev-circle-orb: validate-circle-orb
	circleci orb publish build/package/circleci/orb.yml launchdarkly/ld-find-code-refs@dev:$(TAG)

publish-release-circle-orb: validate-circle-orb
	circleci orb publish build/package/circleci/orb.yml launchdarkly/ld-find-code-refs@$(TAG)

publish-all: publish-release-circle-orb

# Configure Docker authentication for GitHub Container Registry
ghcr-login:
	@echo "Logging in to GitHub Container Registry..."
	@echo "Make sure you have a GitHub Personal Access Token with 'read:packages' and 'write:packages' scopes"
	@echo "Set your token as GITHUB_TOKEN environment variable"
	@echo $${GITHUB_TOKEN} | docker login ghcr.io -u $${GITHUB_USERNAME} --password-stdin
	@echo "Login successful!"

clean:
	rm -rf out/
	rm -f build/pacakge/cmd/bucketeer-find-code-refs
	rm -f build/package/github-actions/bucketeer-find-code-refs-github-action
	rm -f build/package/bitbucket-pipelines/bucketeer-find-code-refs-bitbucket-pipeline

GORELEASER_VERSION=v2.7.0

publish:
	curl -sL https://git.io/goreleaser | \
		VERSION=$(GORELEASER_VERSION) \
		GITHUB_TOKEN=$(GITHUB_TOKEN) \
		bash -s -- --clean --skip=validate

test-publish:
	curl -sL https://git.io/goreleaser | VERSION=$(GORELEASER_VERSION) bash -s -- --clean --skip-publish --skip-validate --snapshot

products-for-release:
	$(RELEASE_CMD) --skip-publish --skip-validate

.PHONY: init test lint compile-github-actions-binary compile-macos-binary compile-linux-binary compile-windows-binary compile-bitbucket-pipelines-binary echo-release-notes publish-dev-circle-orb publish-release-circle-orb publish-all clean build ghcr-login
