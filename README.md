# URL Shortner

## Usage

To store a URL:

```bash
go run main.go http://www.example.com

> 2021/06/23 20:06:03 code: 4oAs25oz6pvcHwySaMvojR
```

To retrieve a URL from a shortcode:

```bash
go run main.go 4oAs25oz6pvcHwySaMvojR

> 2021/06/23 20:09:43 url: http://www.google.com
```

Data is stored in the `bolt.db` file. This is a [bbolt](https://github.com/etcd-io/bbolt) backed key/value store.

## Tests

There are tests for the bbolt store in [bolt_store_test.go](./internal/store/bolt_store_test.go), these use [testify](https://github.com/stretchr/testify) to provide a nicer suite with setup/teardown capabilities.

## Shortcode generation

Rather than write my own shortcode generator, I chose to use [shortuuid](https://github.com/lithammer/shortuuid), this generates based on a UUIDv4, running it via Base57 to keep it URL safe.
