# manioc

<h2 align="center">
<img src="./docs/manioc.svg" width="200" />

**manioc** 📦 IoC/DI Container for Golang 1.18+

</h2>

|⚠️ This library is currently in pre-alpha stage. Breaking changes may occur at any time. Also, we are in a technical research phase. Not provided for production use.|
|-|

## Introduction

Golang has introduced generics feature since version [1.18](https://tip.golang.org/doc/go1.18#generics). To take advantage of this feature, this library provides DI (Dependency Injection) container functionality to enable IoC (Inversion of Control) design patterns on Golang, similar to those commonly used in other class-based languages (such as TypeScript or C#).

Currently (v1.18), Golang's generics feature is rather limited compared to powerful type systems like TypeScript. Due to this, static typing is only effective in a small part of the library, while the rest of the library uses reflection at runtime. Performance may be slightly limited by the use of generics and reflection. Resolving these limitations will require an evolution of Golang itself. We intend to follow golang's upgrades and work on such optimizations as a technical experiment.

## Getting Started

### 0. Installation

Import `github.com/fuzmish/manioc` package in your code:
```go
package main

import (
    "github.com/fuzmish/manioc"
)

func main() {
  // example
  ctr := manioc.NewContainer()
  // ...
}
```
Then run `$ go mod tidy`. Or, run or build your app. [Try with Go Playground](https://go.dev/play/p/goz57y6vpcB).

### 1. Register and Resolve

Here we assume that the struct type `MyService` implements the interface `IMyService`:
```go
type IMyService interface { ... }
// MyService implements IMyService
type MyService struct {}
// ...
```
To register `MyService` as an implementation of `IMyService` in the global container, use the `Register` function:
```go
manioc.Register[IMyService, MyService]()
```
To instantiate a registered implementation of the `IMyService` interface, use the `Resolve` function:
```go
ret, err := manioc.Resolve[IMyService]()
// ret is an instance of MyService
```

Suppose `MyService` has a constructor:
```go
func NewMyService() *MyService { ... }
```
To register a constructor to be called during resolution, use the `RegisterConstructor` function:
```go
manioc.RegisterConstructor[IMyService](NewMyService)
```

To register an existing `MyService` instance as a singleton instance of `IMyService`, use the `RegisterInstance` function:
```go
var myService MyService = //...
manioc.RegisterInstance[IMyService](&myService)
```

### 2. Container

In the above examples, we have used an implicitly prepared global container for registering and resolving dependencies. However, there may be cases where it is inappropriate to use a global container. The `NewContainer` function can be used to create a new container:
```go
ctr := manioc.NewContainer()
```
Use the `WithContainer` option to specify the container to use when registering:
```go
manioc.Register[IMyService, MyService](manioc.WithContainer(ctr))
```
Use the `WithScope` option to specify the container to use when resolving:
```go
ret, err := manioc.Resolve[IMyService](manioc.WithScope(ctr))
```
We will explain about scopes later.

### 3. Cache Policy

When resolving dependencies, you can cache instances in the container. The library provides three types of cache policies:
- `NeverCache`: The instance is always newly created and is not cached. This behavior is known as Transient.
  - This is the default cache policy for the `Register` and `RegisterConstructor`.
- `ScopedCache`: The instance is created only once for the same scope, and then it is used repeatedly. For different scopes, new instances are created. We will explain about scopes later.
- `GlobalCache`: The instance is created only once for the same container, and then it is used repeatedly. This behavior is known as Singleton.
  - This is the default cache policy for the `RegisterInstance` function. Note that the `RegisterInstance` function cannot use any other policy.

The cache policy can be specified at registration using the `WithCachePolicy` option:
```go
manioc.Register[IMyService, MyService](manioc.WithCachePolicy(manioc.ScopedCache))
```
In addition, the following helper functions are provided to make it easy to set policies:
- `RegisterSingleton`, `RegisterSingletonConstructor`: It is equivalent to setting a `GlobalCache` policy and calling `Register` or `RegisterConstructor` respectively.
- `RegisterScoped`, `RegisterScopedConstructor`: It is equivalent to setting a `ScopedCache` policy and calling `Register` or `RegisterConstructor` respectively.
- `RegisterTransient`, `RegisterTransientConstructor`: It is equivalent to setting a `NeverCache` policy and calling `Register` or `RegisterConstructor` respectively.
  - That is, it is equivalent to the default options for `Register` and `RegisterConstructor`.

### 4. Scope

A scope only affects resolution if the cache policy is `ScopedCache`. By using scopes, you can control the range in which instances are cached. A new scope can be created from an existing container or scope using the `OpenScope` function:
```go
// open scope from global container
newScope1, _ := manioc.OpenScope()
// open scoep from specific container
var ctr Container = //...
newScope2, _ := manioc.OpenScope(manioc.WithParentScope(ctr))
// open scope from specific scope
var scope Scope = //...
newScope3, _ := manioc.OpenScope(manioc.WithParentScope(scope))
```
A container can use itself as the most global scope.

The second return value of the `OpenScope` function is a `cleanup` function. By calling this function, you can explicitly close the corresponding scope:
```go
// example
func handler() {
    // open new scope
    scope, cleaup := manioc.OpenScope()
    // Call cleanup before the `handler` function returns.
    // cf. https://go.dev/ref/spec#Defer_statements
    defer cleanup()

    // resolve in this scope
    ret, err := manioc.Resolve[IMyService](WithScope(scope))
    // ...
}
```

### 5. Constructor / Field Injection

In this library, dependency injection is performed on constructors or fields.

The container resolves and injects these arguments when a constructor function registered with `RegisterConstructor` has arguments:
```go
type IFooService interface { ... }
type IBarService interface { ... }

// FooService implements IFooService
type FooService struct {}
// ...

// BarService implements IBarService
type BarService struct {
    Foo IFooService
}
// ...

// constructor
func NewBarService(foo IFooService) *BarService {
    // the argument foo will be injected by the container
    return &BarService{Foo: foo}
}

func main() {
    // register
    manioc.Register[IFooService, FooService]()
    manioc.RegisterConstructor[IBarService](NewBarService)
    // resolve
    ret, _ := manioc.Resolve[IBarService]()
    // then, ret is an instance of BarService,
    // and the field ret.Foo holds an instance of FooService
    // ...
}
```

Field injection is available by tagging fields of a struct:
```go
type IFooService interface { ... }
type IBarService interface { ... }

// FooService implements IFooService
type FooService struct {}
// ...

// BarService implements IBarService
type BarService struct {
    // Setting "inject" value to "manioc" tag will trigger field injection.
    Foo IFooService  `manioc:"inject"`
}
// ...

func main() {
    // register
    manioc.Register[IFooService, FooService]()
    manioc.Register[IBarService, BarService]()
    // resolve
    ret, _ := manioc.Resolve[IBarService]()
    // then, ret is an instance of BarService,
    // and the field ret.Foo holds an instance of FooService
    // ...
}
```

Constructor and field injections can be used together. Also, if the dependencies are registered correctly, resolution is performed recursively.

### 6. Service Key

If multiple implementations are to be registered, they can be keyed with arbitrary values to distinguish them. Use the `WithRegisterKey` option when registering:
```go
manioc.Register[IMyService, MyService]()
manioc.Register[IMyService, MyAnotherService](manioc.WithRegisterKey("another"))
```
When resolving, use the `WithResolveKey` option to specify the key:
```go
ret, _ := manioc.Resolve[IMyService]()
// ret is an instance of MyService
ret2, _ := manioc.Resolve[IMyService](manioc.WithResolveKey("another"))
// ret2 is an instance of MyAnotherService
```
Currently, keys cannot be specified in constructor injections. In field injections, you can specify a key for resolution by appending the `key` option to the tag. However, only string keys can be specified:
```go
type BarService struct {
    Foo        IFooService  `manioc:"inject"`
    FooAnother IFooService  `manioc:"inject,key=another"`
}
```

### 7. Multiple Registration / Resolution

It is possible to register multiple implementations for the same interface, but if there are more than two implementations, `Resolve` will fail because it cannot determine which implementation to use. Instead, by using the `ResolveMany` helper function, you can get a list of resolved instances for all implementations:
```go
// register
manioc.Register[IMyService, MyService1]()
manioc.Register[IMyService, MyService2]()
manioc.Register[IMyService, MyService3](manioc.WithRegisterKey("another"))
manioc.Register[IMyService, MyService4](manioc.WithRegisterKey("another"))

// resolve many
ret1, _ := manioc.ResolveMany[IMyService]()
// ret1 contains an instance of MyService1 and MyService2
ret2, _ := manioc.ResolveMany[IMyService](manioc.WithResolveKey("another"))
// ret2 contains an instance of MyService3 and MyService4
```
Note that `ResolveMany[T]` is equivalent to `Resolve[[]T]`:
```go
ret3, _ := manioc.Resolve[[]IMyService]()
// ret3 contains an instance of MyService1 and MyService2
```
For constructor or field injections, a similar resolution can be performed by making the type of the injected argument or field a slice:
```go
// constructor injection with resolve many
func NewBarService(foos []IFooService) *BarService {...}

// field injection with resolve many
type BarService struct {
    Foos []IFooService  `manioc:"inject"`
}
```

### 8. Must Resolve

The `MustResolve` and `MustResolveMany` functions are variants of the API that can omit error handling. They basically do the same as `Resolve` and `ResolveMany`, but they do not have `error` as a return value, and they will cause `panic` if the dependency cannot be resolved.
```go
// If the resolution is successful, an instance is obtained.
// If it cannot be resolved, it will cause panic.
var instance IFooService = MustResolve[IFooService]()
```

### 9. Query the registry

To check if a dependency on a given interface is registered with a container, use the `IsRegistered` function:
```go
if manioc.IsRegistered[IMyService]() {
    fmt.Println("Registered!")
}
```
This search is affected by the service key. To check for registered dependencies by specifying a key, specify the key with the `WithRegisterKey` option to the `IsRegistered` function:
```go
if manioc.IsRegistered[IMyService](manioc.WithRegisterKey("another")) {
    fmt.Println("Registered! (with key=another)")
}
```
To remove registered dependencies, use the `Unregister` function. This operation deletes all registrations, even if multiple dependencies are registered. Also, the service key must match the key at registration as well as `IsRegistered`:
```go
// Remove all dependencies for IMyService registered with key="another"
manioc.Unregister[IMyService](manioc.WithRegisterKey("another"))
```
The `Unregister` function returns `true` if one or more registrations were deleted, or `false` if none existed.

## Tips

### Known Issues / Limitations

- The following features are being investigated for implementation:
  - Instance disposal (something like `IDisposable` in C#)
  - More flexible and configurable dependency resolution rules (e.g., handling of duplicate keys, handling of keys at ResolveMany, etc.)
  - Instance cache inheritance at scope creation (currently always inherits nothing)
  - We are also considering other nice features that other DI containers have, but due to technical and development resource limitations, there are currently no features that we are considering as a high priority. We are open to your suggestions at any time.
- Currently, the test is very crude. We would like to write clean and comprehensive tests. We also plan to set up CI with GitHub Actions.
- At this moment (v1.18), a type parameter cannot be used as a constraint on another type parameter[[1]](#x-fn1-generics-limitations). For this reason, some APIs cannot perform static type checking and are implemented with runtime reflection. For example:
  - The `Register[TInterface any, TImplementation any](...)` function will compile successfully for any type parameter. We want to constrain the `TImplementation` type by `TInterface` (i.e., we want to do something like `Register[TInterface any, TImplementation TInterface](...)`). It is unclear if this will be possible in future versions of Golang.
  - `RegisterConstructor[TInterface any, TConstructor any](ctor TConstructor, ...)` function will compile successfully for any type parameter. We want to constrain the `TConstructor` type to be any function with a return value of type `TInterface`. It is unclear if this will be possible in future versions of Golang.
  - In summary, the following APIs do not perform static type checking for the reasons discussed above:
    - `RegisterConstructor` (`RegisterSingletonConstructor`, `RegisterScopedConstructor`, `RegisterTransientConstructor`)
    - `Register` (`RegisterSingleton`, `RegisterScoped`, `RegisterTransient`)
  - In contrast, the following APIs are type-checked statically:
    - `RegisterInstance`
  - Also, the following APIs do not need to use type parameter constraints, so the above problem does not occur:
    - `IsRegistered`
    - `Unregister`
    - `Resolve` (`ResolveMany`)
- At this moment (v1.18), Golang does not allow type parameters to be set for individual methods of struct. It is because of this limitation that the API is in the form `Resolve[T](WithScope(scope))` instead of the form `scope.Resolve[T]()`. It is unclear if this can be improved in future Golang releases.
- At this moment (v1.18), Golang does not have optional arguments, so we employ the Functional Options Pattern. It is unclear if this can be improved in future Golang releases.

<small>

<span id="x-fn1-generics-limitations">[1]: cf. https://stackoverflow.com/a/71568095</span>

</small>

### Name and Logo

- In French, the word *manioc* means [cassava](https://en.wikipedia.org/wiki/Cassava). It was chosen as a unique name containing "IoC."
- The logo was created by [@fuzmish](https://github.com/fuzmish). The materials are obtained from [here](https://www.freepik.com/free-vector/hand-drawn-tapioca-illustration_9924650.htm) and [here](https://go.dev/images/gophers/ladder.svg). We believe there are no legal issues, but if you notice any problems, please let us know.

## Contribution

This library is in the pre-alpha stage, so there may be various problems. If you would like to contribute to the project, please let us know your opinions through [GitHub](https://github.com/fuzmish/manioc) issues and pull requests.

## License

[MIT License](https://github.com/fuzmish/manioc/blob/main/LICENSE)
