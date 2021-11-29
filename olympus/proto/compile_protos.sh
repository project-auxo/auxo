#!/usr/bin/env bash

if [ $# -eq 0 ]
  then
    # Get current directory.
  DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
else
  # Proto directiory passed as arg to script.
  DIR=$1
fi

# Find all directories containing at least one prototfile.
# Based on: https://buf.build/docs/migration-prototool#prototool-generate.
for dir in $(find ${DIR} -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq); do
  files=$(find "${dir}" -name '*.proto')

  # Generate all files with protoc-gen-go.
  protoc -I ${DIR} --go-grpc_out=paths=source_relative:${DIR} --go_out=paths=source_relative:${DIR} ${files}
done

echo "Succesfully compiled protos within ${DIR}"