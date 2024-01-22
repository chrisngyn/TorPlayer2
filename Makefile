.PHONY: tailwindcss
tailwindcss:
	@echo "Building Tailwind CSS..."
	npx tailwindcss -i static/style/input.css -o static/style/output.css
	@echo "Done."

.PHONY: templ
templ:
	@echo "Building templates..."
	templ generate
	@echo "Done."

.PHONY: generate
generate: tailwindcss templ
