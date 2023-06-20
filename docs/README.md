# Prerequisites
- Install Kitex. Instructions can be found [here](https://github.com/cloudwego/kitex).
- Install Hertz. Instructions can be found [here](https://github.com/cloudwego/hertz).
- Make sure that your kitex server is configured with server options for supporting generic calls and service registry with Nacos. This helps for json-thrift genericCall and service registry.
- Run Nacos. The easiest way is to run it in a docker container using the image `nacos/nacos-server:2.0.3`. Use `centralx/nacos-server` instead if your chip is ARM.

# How to use the Gateway
1. Put the thrift files for the RPC server in the IDL folder that can be seen at the root directory of the project.![Image 1](https://raw.githubusercontent.com/yiwen101/CardWizards/main/images/image%201.png)
 Please do not create a subdirectory inside the IDL folder and/or change the project structure, or the gateway app might not be able to read the files. Also, if you wish to assign custom routes via annotation, make sure your thrift file's annotation does not contradict with the [IDL Definition Specification for Mapping between Thrift and HTTP](https://www.cloudwego.io/docs/kitex/tutorials/advanced-feature/generic-call/thrift_idl_annotation_standards/).
2. Be at the root directory of the project, run `go run .` in the terminal, then you can see the gateway running.

## Expected Behaviors
### Default Route
For all services and methods, it will be automatically registered with route `/:serviceName/:methodName` under Post method. For example, the `CreateNote` method of the NoteService will be registered under `Post localhost:8080/NoteService/CreateNote`.

### Customized Route
You could also assign customized routes by annotating the idl file like such:
Then the Add method will also be registered at `Get localhost:8080/arith/add`.

## Demonstration on Route
- A successful call to the `Add` method of the `arithmetic` service at path `/arithmetic/Add` under the POST method.
- The call to `/arith/add` under GET method is also successful, as it is a customized route as annotated in the IDL file.
- A case with no matching route.

## Parameter Validation and RPC Call
Parameters to the method as stipulated in the IDL files should be passed as a JSON-encoded body in your HTTP request. The gateway will automatically validate it before making the RPC call. Error messages will be displayed to the user to indicate any issues that lead to failure of validation. If validation is successful, it will forward the RPC response received from upstreams back.

## Demonstration on Validation and RPC Call
- Error message highlighting type mismatch for argument passed in.
- Error message highlighting field mismatch for argument passed in.
- Redundant fields do not hamper validation; can still get correct response.
