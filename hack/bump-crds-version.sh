#!/usr/bin/env bash
# Bump operator.tekton.dev/release (and related) labels in config/ and cmd/
# via hack/release.sh set_version_label — same scope as release-setup-branch
# for manifests, but not the Helm chart (use hack/bump-chart-version.sh).
#
# Usage: ./hack/bump-crds-version.sh <current> <target>
#   current: devel | M.m.p (optional leading v)
#   target:  M.m.p (optional leading v)
#
# Expects repository root cwd or resolves via script path; verifies HEAD is
# release-vX.Y.x matching <target>'s major.minor.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
# shellcheck source=release.sh
source "${SCRIPT_DIR}/release.sh"

usage() {
  echo "Usage: $0 <current> <target>" >&2
  echo "  current: devel | M.m.p [| vM.m.p]" >&2
  echo "  target:  M.m.p [| vM.m.p]" >&2
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

release_mm_from_branch() {
  local b="${1}"
  if [[ "${b}" =~ ^release-v([0-9]+)\.([0-9]+)\.x$ ]]; then
    printf '%s.%s' "${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}"
  else
    return 1
  fi
}

current="${1}"
target="${2}"
target_base="$(strip_leading_v "${target}")"

if [[ ! "${target_base}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "::error::target must be semver M.m.p (got '${target}')" >&2
  exit 1
fi

cd "${REPO_ROOT}"

HEAD_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
EXPECTED_BRANCH="release-v$(formated_majorminorx "${target_base}")"
if [[ "${HEAD_BRANCH}" != "${EXPECTED_BRANCH}" ]]; then
  echo "::error::Must be on '${EXPECTED_BRANCH}' for target ${target_base} (currently on '${HEAD_BRANCH}')" >&2
  exit 1
fi

if [[ "${current}" == "devel" ]]; then
  cur_for_set="devel"
  target_mm="${target_base%.*}"
  head_mm="$(release_mm_from_branch "${HEAD_BRANCH}")" || {
    echo "::error::HEAD is not a release branch" >&2
    exit 1
  }
  if [[ "${head_mm}" != "${target_mm}" ]]; then
    echo "::error::Branch ${HEAD_BRANCH} is line ${head_mm} but target ${target_base} is line ${target_mm}" >&2
    exit 1
  fi
else
  current_base="$(strip_leading_v "${current}")"
  if [[ ! "${current_base}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "::error::current must be devel or semver M.m.p (got '${current}')" >&2
    exit 1
  fi
  target_mm="${target_base%.*}"
  current_mm="${current_base%.*}"
  if [[ "${current_mm}" != "${target_mm}" ]]; then
    echo "::error::current (${current_base}) and target (${target_base}) must share the same major.minor" >&2
    exit 1
  fi
  cur_for_set="${current_base}"
fi

tgt_for_set="${target_base}"

set_version_label "${cur_for_set}" "${tgt_for_set}"
echo "bump-crds-version: updated operator manifest labels (${cur_for_set} -> ${tgt_for_set})."
