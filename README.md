HDS to FHIR
===============================

This project converts JSON data for a patient in the structure used by
[health-data-standards](https://github.com/projectcypress/health-data-standards) and converts
it into the JSON format used by [FHIR](http://hl7.org/implement/standards/fhir/). It also provides
functionality for uploading the data to a FHIR compliant server using the FHIR HTTP interface.

Environment
-----------

This project currently uses Go 1.4 and is built using the Go toolchain.

To install Go, follow the instructions found at the [Go Website](http://golang.org/doc/install).

Following standard Go practices, you should clone this project to:

    $GOPATH/src/github.com/intervention-engine/hdsfhir

To get all of the dependencies for this project, run:

    go get

In this directory.

License
-------

Copyright 2014 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.