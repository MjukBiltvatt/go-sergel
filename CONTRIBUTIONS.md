# Introduction
To contribute to this package, kindly follow the guidelines stated below.
# Guidelines
## Adding features
When adding more features, make sure to use the providers abstractions of HTTP calls which are named after the HTTP method that they concern. If no suiting abstraction is available you can simply create a new one. For example, since there is no GET abstraction you could add the following.
```go
func (p provider) get(uri string, callback providerCallback) error {
    [...]
}
```
The implementation within the function should be rather straight forward. If many abstractions are made which violates the DRY principles then actions to prevent this via refactoring can be done as long as the overall interface is not altered. 
## Changing existing features
When changing an existing feature of the package please keep in mind to try and remain backwards-compatible. This prevents an update of this module to break existing usages. The easiest way to ensure this is to make sure the exposed interfaces remain the same and to try and add to them instead of modifying them. 