GORELEASER_BIN 	:= goreleaser

releaseLocal:
	@$(GORELEASER_BIN) release --clean --snapshot