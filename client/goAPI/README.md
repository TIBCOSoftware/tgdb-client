# GO Client API for TIBCO(R) Graph Database

This is a comprehensive collection of client API referred as 'TGDB-GO-Client' that provides the following capabilities:
* Support for communication with TIBCO Graph Database server via TCP protocol
* Provision to enhance/customize additional protocols such as SSL, HTTP or HTTPS later
* Support for CRUD operations (Insert, Update, Query and Delete) for entity types - Node, Edge, Graph
* Provision to enhance/customize additional types of connections for specific use cases w/ their own connection pool
* Support for 13 types of attributes - bool, byte, short, int, long, float, double, date, time, timestamp, string, BLOB and CLOB
* Support for 6 levels of logging - TRACE, DEBUG, INFO, WARNING, ERROR, and FATAL
* Provision to enhance/customize some other logging mechanism

The entire GO API client implementation has been organized in an easy-to-understand folder structure.

## Folder Structure Overview
* `channel` - A folder that hosts various channel implementations
* `connection` - A folder where bulk of the connection functionality is consolidated
* `exception` - A folder that has various error message types have been implemented
* `iostream` - A folder that implements the serialization and deserialization of messages into byte format
* `logging` - A folder with default log manager implementation, that can be enhanced / augmented
* `model` - All the required data model objects necessary to interact with server
* `pdu` - Various request and response message types that the server recognizes
* `query` - A folder with basic implementation of query API
* `samples` - Various examples demonstrating how to use various APIs
* `types` - A folder that consolidates various interfaces that are implemented and can be enhanced/augmented
* `utils` - Various utilities used by various API
* `main.go` - A driver program that can be used to test, execute different samples
* `README.md` - This file itself

This set of API are compatible with TIBCO Graph Database v2.0.0 and above.

Please see [TIBCO Graph Database Community](https://community.tibco.com/products/tibco-graph-database) for more information on Releases and Reference material.

## Frequently Asked Questions (FAQ)

###(1) What are the pre-requisites to use these API?

There are no third party software and/or components required to use these APIs. These APIs have been implemented using
[GO language packages](https://golang.org/pkg/).
Please make sure you have the following appropriately set:
 * GO environment a.k.a. $GOROOT, and 
 * GO Repositories top level folder a.k.a. $GOPATH/{src|bin|pkg} - where you can sync from github or any other source code repositories
In other words, all you need to have is GO language environment - either IDE supporting GO language development OR GO language 
toolkit (command line tools) available for use on a console/command prompt.

###(2) Are there any samples or examples as a reference?

The API source code folder structure has a 'samples' folder that hosts various sample programs and associated
initdb/tgdb configuration files.
 
All command-line instructions in this document assume `bash` as the target shell, running in Linux or Mac OS X. 
If you're using Windows, you'll have to work out the PowerShell equivalents.

###(3) How do I execute / test the sample programs?

As mentioned earlier, the driver program to execute is 'main.go'
   
    $ ls -l main.go
    -rw-r--r--  1 tgdb  staff  1643 Mar  1 19:30 main.go
    $ go run main.go | tee SampleTest.log

You can edit 'main.go' to comment / uncomment required sample execution.
Each of the sample execution is documented within 'main.go' with appropriate and required steps. If you need to 
add / edit / modify new tests, 'main.go' is the starting point before moving on to bigger integrations and solutions.

###(4) How do I compile TGDB-GO-Client API?

It is standard GO compilation as shown below:
   
    $ cd client/go
    $ go build ./...

###(5) How do I verify any functional changes / enhancements / augmentation via unit tests?

It is standard GO testing as shown below:
   
    $ cd client/go
    $ go test ./...

###(6) Is there a provision to override/change log level?

Yes, the default log level is set at 'DEBUG' state, as part of defaultLogger(). However, it can be overridden in 3 ways:

* Edit/Modify SimpleLogger.go to change the defaultLogger()
* Instead of using defaultLogger(), use NewLogger() that accepts the log level as one of the function parameters
* Develop CustomLogger.go (similar to SimpleLogger.go) that implements all the interfaces from types/TGLogger.go

###(7) Can I implement/use my own custom connection factory with custom functionality?

Yes, it is possible for you to develop and implement your own custom connection functionality. 
Please take a look at connection/defaultConnectionFactory.go as a reference, and make sure that your custom
connection factory implementation is in-line with the interface implementations from types/TGConnection.go

###(8) Are there any known issues / gotchas that we should be aware of?

As always, this set of API implemented in GO language is work-in-progress, and is continuously being enhanced to add
new features / options / capabilities / functionality that may be needed to effectively and efficiently communicate
with TIBCO Graph Database server engine.

The following are still being implemented and tested:
* Support for communication via SSL (TLSv1.2+) over TCP
* Support for administrative API that will mimic 'tgdb-admin' functionality
