#!/bin/sh
# Usage: ./vendor-brotli.sh path/to/google/brotli
#
# Copies the upstream C sources into c/ unmodified and regenerates the enc/
# shims. cgo only compiles C files that sit directly in the package directory,
# so each shim #includes one upstream source file.
set -e

src="$1"
if [ ! -f "$src/c/include/brotli/encode.h" ]; then
	echo "usage: $0 path/to/google/brotli" >&2
	exit 1
fi

cd "$(dirname "$0")"

rm -rf c
mkdir -p c
cp -r "$src/c/common" "$src/c/enc" "$src/c/include" c/
cp "$src/LICENSE" LICENSE

rm -f enc/brotli_*.c
for f in c/common/*.c c/enc/*.c; do
	dir=$(basename "$(dirname "$f")")
	name=$(basename "$f")
	echo "#include \"../$f\"" > "enc/brotli_${dir}_${name}"
done

commit=$(git -C "$src" rev-parse HEAD 2>/dev/null || echo unknown)
version=$(git -C "$src" describe --tags --match 'v*' 2>/dev/null || echo unknown)
cat > VERSION.md <<EOF
Brotli Version
---

Vendored Brotli encoder from https://github.com/google/brotli

Upstream version: \`$version\`

Upstream commit: \`$commit\`
EOF
