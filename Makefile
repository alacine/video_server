all:
	echo "TODO"
	echo "This will do nothing"

clean:
	find . -type f ! -regex '^\./\.git/.*' ! -regex '.+\..+' ! -name Makefile -delete
