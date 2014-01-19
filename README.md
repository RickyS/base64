Go Package to convert base64 to and from 'binary' and ascii-hex-text  
=======================================================

This was written as an exercise in Go, as it duplicates the existing Go library package to
do essentially the same thing, encoding/base64.  As a bonus, this package also converts to and from 
hex-ascii text.   This package was tested running Go 1.2 on Ubuntu 13.10 64-bit x86.

For example:  
      "Zm9vYmFyeQ==" <==> "666f6f62617279"   // Base64 <==> Ascii Hex Text  
      "Zm9vYmFyeQ==" <==> "foobary"          // Base64 <==> Binary, displayed here as %s text.  

To install:   
       $ cd &lt;Your src directory&gt;  
       $ go get -d github.com/RickyS/base64   
       $ mv base64/decode.go .   
       $ go install -v -x base64  
       $ go run decode.go  
       
     

Testing code is currently in the main program decode.go, which is included.  Be sure to remove the file decode.go from the
base64 directory.  

TODO:  Fix that decode.go nuisance.  Change to a proper test.
