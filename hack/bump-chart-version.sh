#!/usr/bin/env bash
# Bump Helm chart metadata under charts/tekton-operator for a release.
#
# Two arguments (think: FROM → TO):
#   1) Current version string already present in the chart YAML (what you are replacing).
#   2) Target release you want after the script runs (what it should say next).
#
# Usage:
#   ./hack/bump-chart-version.sh <current> <target>
#
# Examples:
#   ./hack/bump-chart-version.sh devel 0.79.1
#     # files say "devel" today; make them say 0.79.1 / v0.79.1 where appropriate
#   ./hack/bump-chart-version.sh 0.79.0 0.79.1
#     # files say v0.79.0 in labels; bump to v0.79.1 and Chart version 0.79.1
#   ./hack/bump-chart-version.sh v0.79.0 v0.79.1
#     # same as above; leading "v" is optional on the command line
#
# Same convention as hack/release.sh: numeric versions without a leading "v"
# is fine; "devel" is literal. CRD labels in YAML use "v" + semver.
#
# sed -i is not portable between GNU sed (Linux) and BSD sed (macOS). This script
# detects GNU sed and picks the right in-place flags. BusyBox sed on Alpine
# typically behaves like GNU for -i without a backup argument.

set -euo pipefail

usage() {
  echo "Usage: $0 <current> <target>" >&2
  echo "  current — version string already in the chart (e.g. 0.79.0 or devel)" >&2
  echo "  target  — release you want afterward (e.g. 0.79.1)" >&2
  echo "Example: $0 0.79.0 0.79.1" >&2
  exit 1
}

if [[ "${#}" -ne 2 ]]; then
  usage
fi

strip_leading_v() {
  local s="${1}"
  s="${s#v}"
  printf '%s' "${s}"
}

# "current" / "target" on the CLI; keep internal names short.
raw_current="${1}"
raw_target="${2}"

if [[ "${raw_current}" == "devel" ]]; then
  old_label="devel"
else
  old_base="$(strip_leading_v "${raw_current}")"
  old_label="v${old_base}"
fi

new_base="$(strip_leading_v "${raw_target}")"
new_label="v${new_base}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
CHART_DIR="${REPO_ROOT}/charts/tekton-operator"
CRD_GLOBS=(
  "${CHART_DIR}/templates/kubernetes-crds.yaml"
  "${CHART_DIR}/templates/openshift-crds.yaml"
)

for f in "${CRD_GLOBS[@]}"; do
  if [[ ! -f "${f}" ]]; then
    echo "error: expected file missing: ${f}" >&2
    exit 1
  fi
done

if [[ ! -f "${CHART_DIR}/Chart.yaml" ]]; then
  echo "error: expected file missing: ${CHART_DIR}/Chart.yaml" >&2
  exit 1
fi

# Escape for ERE (so dots in v0.79.0 are literal, unlike in hack/release.sh).
ere_escape() {
  printf '%s' "${1}" | sed 's/[[\.^$*+?{|()]/\\&/g'
}

old_ere="$(ere_escape "${old_label}")"

configure_sed_inplace() {
  local _v
  _v="$(sed --version 2>&1 || true)"
  if [[ "${_v}" == *"GNU sed"* ]]; then
    SED_INPLACE=(sed -i -E)
  else
    SED_INPLACE=(sed -i '' -E)
  fi
}

configure_sed_inplace

echo "Updating Helm chart at ${CHART_DIR}"
echo "  CRD labels: ${old_label} -> ${new_label}"
echo "  Chart.yaml: version -> ${new_base}, appVersion -> ${new_label}"

"${SED_INPLACE[@]}" \
  -e "s/(operator\\.tekton\\.dev\\/release: )\"?${old_ere}\"?/\\1${new_label}/g" \
  -e "s/(^[[:space:]]*version: )\"?${old_ere}\"?/\\1${new_label}/g" \
  "${CRD_GLOBS[@]}"

old_base_ere="$(ere_escape "$(strip_leading_v "${old_label}")")"
# Chart version is semver without "v"; appVersion includes "v".
if [[ "${old_label}" == "devel" ]]; then
  "${SED_INPLACE[@]}" \
    -e "s/^version:[[:space:]]*\"?devel\"?/version: ${new_base}/" \
    -e "s/^appVersion:[[:space:]]*\"?devel\"?/appVersion: ${new_label}/" \
    "${CHART_DIR}/Chart.yaml"
else
  "${SED_INPLACE[@]}" \
    -e "s/^version:[[:space:]]*\"?${old_base_ere}\"?/version: ${new_base}/" \
    -e "s/^appVersion:[[:space:]]*\"?${old_ere}\"?/appVersion: ${new_label}/" \
    "${CHART_DIR}/Chart.yaml"
fi

echo "Done."
