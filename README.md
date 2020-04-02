<!---
// ********************************************************************
// * README.md                                                        *
// *                                                                  *
// * 2020-03-15 First Version, JR                                     *
// * 2020-03-16 Adds authorizer summary, JR                           *
// * 2020-03-17 Adds unit testing instructions, JR                    *
// * 2020-03-18 Adds e2e tests instructions, JR                       *
// *                                                                  *
// * Contains a brief summary of the project and its instructions,    *
// * to build, execute and run relevant commands related.             *
// ********************************************************************
-->

## Authorizer

The auhorizer app is writtern in [GoLang](https://golang.org/), language which have its own [Docker official image](https://hub.docker.com/_/golang), so is strongly recommended to usage of the pre-built image to have the most possible similar setup and be consistent with commands used for `build`, `test` and `execute` the application. This project contains their own `Dockerfile` to create a container with all dependencies and the application binary already built and ready to be used.

As a part to have mostly automated this project also contains a main `Makefile` with rules that executes all commands required to build, run unit testing, run code coverage, and end to end testing, same steps that below in this documentation are explained.

### Requirements

* [Docker](https://docs.docker.com/install/)

### Build

First step and the most important one is `Create` a docker image, since next steps have an strong dependency from this image, to create our docker image, locate your terminal in root authorizer directory and execute docker build instruction:

* $`docker image build -t authorizer:go .`

### Unit testing

All project packages includes its own unit testing file with their test's scenarios, `Golang` provides native testing suit to do your own unit testing but does not support `Assert` or `Verify` statements, To handle a more comprehensive unit testing in comparison with most used programming languages the app uses a [Third party implementation for asserts](https://github.com/stretchr/testify/assert).

To retrieve unit testing results, execute the go native test suit command, as shown below:

* $`docker run -i authorizer:go make unit-test`

### Code coverage

As same as go handles unit tasting natively, go can handle code coverage with same tool, to see percentiles related with the current code coverage from our application we can run the next command:

* $`docker run -i authorizer:go make cover`

### e2e tests

This project contains scripts to test the authorizer application from end to end (e2e), those scripts automates the application execution with a defined `in` sample in order to get an expected `out`, if the `result` of a test is different from the `out` the the test report that test scenario as `FAIL`, otherwise the test scenario reports an `OK`, to run the e2e tests is necessary to execute next command line:

* $`docker run -i authorizer:go make test`

Test scenarios can be added easily to the e2e, this can be done adding a new directory inside `./test`, the test scenario only needs 2 files inside, `in` which contains all our `json` lines input to execute, and `out` that represents the expected output.

NOTE: ** Take in consideration: all changes to the app needs a new docker build. **

### Run binary

To run the application we only need to run the next single command, replacing `$FILE` for the real path of the file which contains the `json` input operations to execute:

* $`docker run -i authorizer:go authorizer < $FILE`
