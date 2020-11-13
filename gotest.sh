#!/usr/bin/env bash

go version

skipedPackages=(
  "salt"
)
packages=$(go list -mod vendor ./...);


for pkg in ${packages}; do
  skipTest=false;
  for spkg in "${skipedPackages[@]}"; do
    x=$(echo "$pkg" | grep -c "$spkg");
    if [[ "$x" -eq 1 ]]; then
      skipTest=true
    fi
  done
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
