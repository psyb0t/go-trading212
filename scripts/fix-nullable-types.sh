#!/bin/bash
# Rewrites OpenAPI 3.1 style nullable types in the downloaded spec into their
# OpenAPI 3.0 form, which oapi-codegen understands. The spec declares 3.0.1
# but some schemas use a type list such as `type: [number, 'null']`. Each one
# becomes `type: number` plus `nullable: true`. Reads api/openapi.yaml and
# writes api/openapi.fixed.yaml.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SPEC_IN="${SPEC_DIR}/openapi.yaml"
readonly SPEC_OUT="${SPEC_DIR}/openapi.fixed.yaml"

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

# Handles both YAML forms of a two-item type list:
#   type: [number, 'null']
# and
#   type:
#     - number
#     - 'null'
# Any other list is copied through untouched, so oapi-codegen fails loudly on
# it instead of this script guessing.
fix_types() {
	awk '
	function unquote(s) {
		gsub(/^[ \t]+|[ \t]+$/, "", s)
		gsub(/^["\047]|["\047]$/, "", s)
		return s
	}
	function emit(indent, a, b) {
		if (a == "null" && b != "null") { a = b }
		else if (b != "null") { return 0 }
		print indent "type: " a
		print indent "nullable: true"
		fixed++
		return 1
	}
	function flush_pending() {
		if (!pending) { return }
		if (!(n == 2 && emit(pindent, items[1], items[2]))) {
			print phead
			for (i = 1; i <= n; i++) { print raw[i] }
		}
		pending = 0
		n = 0
	}
	{
		if (pending) {
			if (match($0, /^[ \t]*- /)) {
				n++
				raw[n] = $0
				items[n] = unquote(substr($0, RLENGTH + 1))
				next
			}
			flush_pending()
		}
		if (match($0, /^[ \t]*type:[ \t]*$/)) {
			pending = 1
			phead = $0
			pindent = substr($0, 1, index($0, "type:") - 1)
			n = 0
			next
		}
		if (match($0, /^[ \t]*type:[ \t]*\[.*\][ \t]*$/)) {
			indent = substr($0, 1, index($0, "type:") - 1)
			list = $0
			sub(/^[^\[]*\[/, "", list)
			sub(/\][ \t]*$/, "", list)
			if (split(list, parts, ",") == 2 && emit(indent, unquote(parts[1]), unquote(parts[2]))) {
				next
			}
		}
		print
	}
	END {
		flush_pending()
		print fixed > "/dev/stderr"
	}
	' "$SPEC_IN"
}

main() {
	local tmp count
	tmp="$(mktemp "${SPEC_DIR}/.openapi.fixed.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}' '${tmp}.count'" EXIT

	fix_types >"$tmp" 2>"${tmp}.count"
	count="$(cat "${tmp}.count")"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$SPEC_OUT"
	log INFO "rewrote nullable type lists count=${count} out=${SPEC_OUT##*/}"
}

main "$@"
