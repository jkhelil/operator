#!/usr/bin/env bash
# Print operator release version from config (devel or M.m.p without v).
# Source: first operator.tekton.dev/release under metadata in
# config/kubernetes/base/300-operator_v1alpha1_dashboard_crd.yaml
#
# Usage: ./hack/infer-operator-release-from-config.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
FILE="${REPO_ROOT}/config/kubernetes/base/300-operator_v1alpha1_dashboard_crd.yaml"

if [[ ! -f "${FILE}" ]]; then
  echo "error: missing ${FILE}" >&2
  exit 1
fi

line="$(grep -m1 -E '^[[:space:]]*operator\.tekton\.dev/release:' "${FILE}" || true)"
if [[ -z "${line}" ]]; then
  echo "error: no operator.tekton.dev/release in ${FILE}" >&2
  exit 1
fi

raw="$(
  printf '%s\n' "${line}" |
    sed -E 's/^[[:space:]]*operator\.tekton\.dev\/release:[[:space:]]+//'
)"
raw="${raw%$'\r'}"
if [[ "${raw}" == \"*\" ]]; then
  raw="${raw#\"}"
  raw="${raw%\"}"
elif [[ "${raw}" == \'*\' ]]; then
  raw="${raw#\'}"
  raw="${raw%\'}"
fi

if [[ "${raw}" == "devel" ]]; then
  printf '%s\n' "devel"
  exit 0
fi

if [[ "${raw}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  printf '%s\n' "${raw#v}"
  exit 0
fi

echo "error: unexpected operator.tekton.dev/release value '${raw}'" >&2
exit 1
