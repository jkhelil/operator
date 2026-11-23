#!/usr/bin/env bash

# Copyright 2026 The Tekton Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Thin wrapper around tektoncd/plumbing's shared CRD compatibility check.
# Run the plumbing script (do not source it).
#
# Until plumbing#3820 is merged and vendored, scripts are copied into
# vendor/github.com/tektoncd/plumbing/scripts/ for fork testing.

set -o errexit
set -o nounset
set -o pipefail

ROOT="$(git rev-parse --show-toplevel)"
export CRD_DIRS="${CRD_DIRS:-config/base/generated-crds}"

SCRIPT="${ROOT}/vendor/github.com/tektoncd/plumbing/scripts/verify-crd-compat.sh"
if [[ ! -f "${SCRIPT}" ]]; then
  echo "error: missing ${SCRIPT}" >&2
  echo "Copy from plumbing PR #3820 or bump/vendor plumbing after merge." >&2
  exit 1
fi

exec bash "${SCRIPT}"
