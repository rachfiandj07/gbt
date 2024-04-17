usage: FORCE
	# See targets in Makefile (e.g. "buildlet.darwin-amd64")
	exit 1

FORCE:

.PHONY: gbt
gbt: FORCE
	@echo " >> building gbt binaries..."
	@go build -o cmd/gbt/gbt cmd/gbt/main.go
	@echo " >> gbt built."
	@echo "executing gbt..."
	@./cmd/gbt/gbt
	@echo "gbt is running."

bg: FORCE
	@echo " >> building gbt bg binaries..."
	@go build -o cmd/gbt/gbt_bg cmd/gbt/main.go
	@echo " >> gbt bg built."
	@echo "executing gbt bg..."
	@./cmd/gbt/gbt_bg
	@echo "gbt bg is running."

.PHONY: default
default: gbt

.DEFAULT_GOAL := gbt