# Go bindings for the Brotli compression library

[![GoDoc](https://godoc.org/github.com/itchio/go-brotli/enc?status.svg)](https://godoc.org/github.com/itchio/go-brotli/enc)

This is a fork of Mike Houston's <https://github.com/kothar/brotli-go>. It
only provides the encoder. For decompression use a pure Go decoder such as
<https://github.com/itchio/dskompress>.

See <https://github.com/google/brotli> for the upstream C source, and the
`VERSION.md` file for the currently vendored version.

### Usage

Instead of including potentially-outdated examples in the README,
please refer to the `Examples` tests on the godoc page:
<https://godoc.org/github.com/itchio/go-brotli/enc>

### Bindings

The upstream `c/common`, `c/enc` and `c/include` directories are copied
unmodified into `c/`. cgo only compiles C files that sit directly in the
package directory, so `enc/` contains one `brotli_*.c` shim per upstream
source file that `#include`s it.

To update the vendored version, check out the desired upstream tag and run:

    ./vendor-brotli.sh path/to/google/brotli

### Links

  * original bindings: <https://github.com/kothar/brotli-go>
  * upstream cgo bindings (requires separate library compilation): <https://github.com/google/brotli/tree/master/go/cbrotli>

### License

Brotli and these bindings are open-sourced under the MIT License - see the LICENSE file.
