#!/bin/bash -eux

pushd dp-sessions-api
  make build
  cp build/dp-sessions-api Dockerfile.concourse ../build
popd
