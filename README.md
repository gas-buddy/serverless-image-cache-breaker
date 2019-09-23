# serverless-image-cache-breaker
Command line tool to deleted resized images - usually after an image has been updated

## Origin
The [serviceless-image-resizing](https://github.com/gas-buddy/serverless-image-resizing) tool (at the time of this writting) does not support breaking the resized imaged cache when an image is updated.  This tool is to address that.  In hindsight the effort to build this tool would have been better placed automating this functionality, but it is already done.  Consider using this as a starting point whenever the time comes to automate this process.

## How to Use
1. You will need to have Go setup on your machine to build the executable. https://golang.org/doc/install
2. Setup your local AWS configuration to point to the correct environment via `aws configure`
3. Run the command, for example:
```go
go run main.go -bucket gb-images-stg -file parking.png -ignore home
```

**-help** will provide some guidance and describe the following options.

**-bucket** specifies which bucket to scan.  For staging you probably want gb-images-stg

**-file** the name of the image file whose cache you want to bust

**-ignore** directories you do not want to delete the image from, you will probably want to specify the source directory here, which is "home"

Upon execution you will be previed the action about the be taken, including which files will be removed and which ones will be ignored.  You can safely abort at this point.
