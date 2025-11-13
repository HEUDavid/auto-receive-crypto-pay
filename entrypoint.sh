#!/usr/bin/env sh

set -e
envsubst < conf/conf.toml.tpl > conf/conf.toml
exec ./receivepay
