#!/usr/bin/env bash

go version

packages=$(go list -mod vendor ./...);


for pkg in ${packages}; do
  skipTest=false;
  if [[ ${skipTest} = false ]]; then
    x=$(go test -mod=vendor -v "$pkg");
    # shellcheck disable=SC2181
    if [ ! "${x}" ]; then
      echo "\033[0m${x}";
      exit 1
    fi
    echo "\033[0;37m${pkg} \033[0;32mPASSED!.";
  fi
done
