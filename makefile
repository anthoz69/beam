help:
	# swapgap-web
	#
	# make dev -- start live reload on port 8000

dev:
	gin -p 8000 -a 8080 -x vendor -x src -x public -x static --all -i

style:
	npm run dev