#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-sessions-api
  make audit
popd   