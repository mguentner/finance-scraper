#!/bin/bash

BUILD_TAGS_FILE=build.tag
SOURCE=finance.go
DESTINATION=test/mock_finance.go
PACKAGE=test

mockgen -source=${SOURCE} -destination=${DESTINATION} -package=${PACKAGE} -copyright_file ${BUILD_TAGS_FILE}