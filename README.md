# dccn-tools
protoc-gen-ankr:

​     This is a protobuffer code generation tool based on the actual needs of ankr. It is a secondary development based on the official protoc-gen-go code generation tool. It adds log tracking content and some usage optimizations, such as client calls, etc....



log:

​     The log package is based on the secondary encapsulation of the Uber's zap log package to meet the needs of Ankr distributed log tracking. Of course, in the log output, we add a lot of new elements, such as service name, host name, information type, etc.

