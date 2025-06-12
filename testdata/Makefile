SHELL := bash

.PHONY: all

go := go
png-size := $(go) run ../_tools/png-size.go

all: go-logo-blue.png go-logo-blue-2.png

go-logo-blue-2.png: go-logo-blue.svg
	magick -background none $< -colors 2 $@

go-logo-blue.png: go-logo-blue.svg
	magick -background none $< $@

go-logo-blue.svg:
	curl -O https://go.dev/images/go-logo-blue.svg

# Doesn't work: outputs zTXt instead of tEXt
# magick $< -define PNG:text-type=text -set comment $$(perl -e 'print q[x] x (3072 - 20 - -s q[$@])') $@

go-favicon-3069.png: go-favicon-1.png
	$(png-size) -size=3069 $< $@

go-favicon-3071.png: go-favicon-1.png
	$(png-size) -size=3071 $< $@

go-favicon-3072.png: go-favicon-1.png
	$(png-size) -size=3072 $< $@

go-favicon-3073.png: go-favicon-1.png
	$(png-size) -size=3073 $< $@
