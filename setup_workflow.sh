#!/bin/bash

# Colors
ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_RED=$ESC_SEQ"31;01m"
COL_GREEN=$ESC_SEQ"32;01m"
COL_YELLOW=$ESC_SEQ"33;01m"
COL_BLUE=$ESC_SEQ"34;01m"
COL_MAGENTA=$ESC_SEQ"35;01m"
COL_CYAN=$ESC_SEQ"36;01m"

echo -e "$COL_BLUE""Installing code generator...$COL_RESET"
go get -v github.com/clipperhouse/gen

echo -e "$COL_BLUE""Installing live code reloader...$COL_RESET"
go get -v github.com/codegangsta/gin

echo -e "$COL_BLUE""Installing dependencies manager...$COL_RESET"
go get -v github.com/tools/godep

echo -e "$COL_BLUE""Installing Atom plugins dependencies...$COL_RESET"
go get -v github.com/redefiance/go-find-references
go get -v golang.org/x/tools/cmd/gorename
go get -v code.google.com/p/rog-go/exp/cmd/godef

echo -e "$COL_GREEN""Done.$COL_RESET"
