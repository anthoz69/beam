help:
	# swapgap-web
	#
	# make dev -- start live reload on port 8000

dev:
	goreload -x vendor -x src -x public -x static --all

style:
	npm run dev

cstyle:
	node-sass --output-style compressed --include-path node_modules src/scss/main.scss > assets/style.css