# ********************************************************************
# * Makefile                                                         *
# *                                                                  *
# * 2020-03-17 First Version, JR                                     *
# *                                                                  *
# * File with dynamic make rules associated to test authorizer,      *
# * This file get all dirs inside test dir and execute the app       *
# * passing the "in" file to get the "result",                       *
# * If no diff between "result" and "out" then OK, else FAIL.        *
# *                                                                  *
# * Usage:                                                           *
# * $ make all                                                       *
# ********************************************************************

BIN := authorizer
TESTS = $(shell ls -d */ | cut -f1 -d'/')

all: $(TESTS)

$(TESTS):
	@$(BIN) < $@/in | tee $@/result > /dev/null
	@DIFF=`diff -w $@/out $@/result | tee` && \
	if [ -z "$$DIFF" ]; \
	then echo $@:OK.; \
	else echo $@:FAIL.; \
	fi

.PHONY: all $(TESTS)
