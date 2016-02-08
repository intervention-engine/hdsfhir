HDS to FHIR [![Build Status](https://travis-ci.org/intervention-engine/hdsfhir.svg?branch=master)](https://travis-ci.org/intervention-engine/hdsfhir)
=====================================================================================================================================================

This project converts JSON data for a patient in the structure used by[health-data-standards](https://github.com/projectcypress/health-data-standards) into the [HL7 FHIR DSTU2](http://hl7.org/fhir/DSTU2/index.html) models defined in the Intervention Engine [fhir](https://github.com/interventionengine/fhir) project.

*NOTE: The HDS-to-FHIR conversion focuses only on those data elements that are needed by Intervention Engine. It is not a complete and robust conversion.*

Building hdsfhir Locally
------------------------

For information on installing and running the full Intervention Engine stack, please see [Building and Running the Intervention Engine Stack in a Development Environment](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md).

The *hdsfhir* project is a Go library. For information related specifically to building the code in this repository (*hdsfhir*), please refer to the following sections in the above guide:

-	(Prerequisite) [Install Git](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#install-git)
-	(Prerequisite) [Install Go](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#install-go)
-	[Clone hdsfhir Repository](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#clone-hdsfhir-repository)

To build the *hdsfhir* library, you must install its dependencies via `go get` first, and then build it:

```
$ cd $GOPATH/src/github.com/intervention-engine/hdsfhir
$ go get
$ go build
```

For information on using the *uploadhds* tool to convert HDS patients to FHIR and upload them to a FHIR server, please refer to the [uploadhds](https://github.com/intervention-engine/tools#uploadhds) section of the [tools](https://github.com/intervention-engine/tools) repository README.

Using hdsfhir as a library
--------------------------

The following is a simple example of converting a HDS patient to the corresponding slice of FHIR models:

```go
import (
	"encoding/json"
	"io/ioutil"
	"github.com/intervention-engine/hdsfhir"
)

func ExampleConversion(pathToHdsPatient string) ([]interface{}, error) {
	data, err := ioutil.ReadFile(pathToHdsPatient)
	if err != nil {
		return nil, err
	}

	p := &hdsfhir.Patient{}
	err = json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	return p.FHIRModels(), nil
}
```

License
-------

Copyright 2016 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
