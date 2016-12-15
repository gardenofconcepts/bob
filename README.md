# Builder

## Run

    builder -pattern "*.build.yml" this/is/my/optional/search/path


## Configuration

Expected environment variables:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`

Alternative: The AWS credentials file – located at `~/.aws/credentials` on
Linux, OS X, or Unix, or at `C:\Users\USERNAME\.aws\credentials` on Windows.

## Arguments

      -debug
            Enable debug mode (Log level: debug)
      -exclude string
            Excludes directories with this pattern (e.g. **/node_modules/**,.git)
      -force
            Rebuild data without checking remote
      -include string
            Pattern for directory traversal (default "**")
      -pattern string
            File pattern for build files (default "*.build.yml")
      -s3-bucket string
            Specify S3 bucket name (default "cmsbuild")
      -s3-region string
            Specify S3 region (default "eu-central-1")
      -skip-download
            Don't download builds
      -skip-upload
            Don't upload builds
      -verbose
            Enable verbose mode (Log level: info)
      -version
            Show version

## Build

Run ``make`` to build binaries:

    export GOPATH=`pwd`
    make init
    make build

## Golang setup (OS X)

    brew update
    brew install go
