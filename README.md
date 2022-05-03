## ead-stream-parser

#### Overview
This repo contains EXPERIMENTAL code.  I am exploring the use of  
stream parsing with EAD XML files.

##### Sample usage on *nix:
* To generate JSON:
`go run gen-json.go Omega-EAD.xml | tee omega-stream-parser-output.json`
* To generate text:
`go run gen-text.go Omega-EAD.xml | tee omega-stream-parser-output.text`
