HDS to FHIR
===============================

This project converts JSON data for a patient in the structure used by
[health-data-standards](https://github.com/projectcypress/health-data-standards)
into the [HL7 FHIR](http://hl7.org/implement/standards/fhir/) models defined
in the Intervention Engine [fhir](https://github.com/interventionengine/fhir)
project.

Environment
-----------

This project currently uses Go 1.4 and is built using the Go toolchain.

To install Go, follow the instructions found at the [Go Website](http://golang.org/doc/install).

Following standard Go practices, you should clone this project to:

    $GOPATH/src/github.com/intervention-engine/hdsfhir

To get all of the dependencies for this project, run:

    go get

in this directory.

To run all of the tests for this project, run:

    go test ./...

in this directory.

License
-------

Copyright 2015 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.