#!/bin/bash
set -e

SOURCE_DIR=${SUB_DIR:-"gcr.io"}
PUBLIC_DIR=${PUBLIC_DIR:-"./publsh"}

if test -z "${SOURCE_DIR}"; then
  echo "SOURCE_DIR is nil, skip!"
  exit -1
fi

if test -z "${PUBLIC_DIR}"; then
  echo "PUBLIC_DIR is nil, skip!"
  exit -1
fi

echo "## generate json ##################"

gcrmirrors \
  -sourceDir "/github/workspace/${SOURCE_DIR}" \
  -publicDir ${PUBLIC_DIR}

echo "## Done. ##################"
