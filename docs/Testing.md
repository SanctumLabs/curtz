# Testing

Testing is important & that is no different to this application. There are different types of tests setup in this application:

1. Unit Tests
2. Integration Tests

## Unit Tests

By far the simplest to run & the easiest to setup. These test a unit of work(a function or a method). If a function or a method depends on executing another function or method that is not within the current functions or method's scope, then it can be mocked. This allows finer grained control of the test.

Mocking is done using [mock](github.com/golang/mock) & allows unit testing to be simpler and quicker. Documentation is already provided in the mock site, but in general, mocking is done on interfaces and the return of these interfaces is controlled within the test.

To generate mocks you can run this command:

``` bash
mockgen -source <FILE_TO_MOCK> -destination <DESTINATION_OF_MOCK> -package mocks
```

> Ensure the <DESTINATION_OF_MOCK> is within the [test mocks](../app/test/mocks) directory
