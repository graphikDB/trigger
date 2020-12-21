version := "0.0.0"

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "----------------------------------------------------------------"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo "----------------------------------------------------------------"

patch: ## bump sem version by 1 patch
	bumpversion patch --allow-dirty

minor: ## bump sem version by 1 minor
	bumpversion minor --allow-dirty

tag: ## tag the repo (remember to commit changes beforehand)
	git tag v$(version)

push:
	git push origin v$(version)
