help:
	# swapgap-web
	#
	# make dev -- start live reload on port 8000

dev:
	goreload -x vendor -x src -x public -x static --all

style:
	npm run dev