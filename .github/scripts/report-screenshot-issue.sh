#!/usr/bin/env bash
# Post one CI summary to GitHub issue #14 (the screenshot board).
#
# PNG files under SCREENSHOT_DIR are embedded in the comment when
# uploads.github.com accepts this token. Actions GITHUB_TOKEN (ghs_) is
# rejected by that endpoint, so the usual result is a file list plus a
# link to the window-screenshots artifact. Nothing is committed.
#
# Re-runs of the same workflow run update the comment marked with that run id.
# Bash 3.2 (macOS /bin/bash) is enough.
set -euo pipefail

ISSUE=14
ARTIFACT_NAME=window-screenshots

dry_run=0
repo=""
run_id=""
server=""
status=""
event_name=""
branch=""
sha=""
run_url=""
attempt=""
pr_url=""
pr_number=""
test_outcome=""
screenshot_dir=""
token=""
artifact_url=""
png_names=()
png_paths=()
png_urls=()

usage() {
  echo "usage: report-screenshot-issue.sh [--dry-run]" >&2
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

collect_pngs() {
  png_names=()
  png_paths=()
  png_urls=()
  if [[ ! -d "$screenshot_dir" ]]; then
    return 0
  fi
  local path name
  while IFS= read -r path; do
    [[ -n "$path" ]] || continue
    name="$(basename "$path")"
    png_paths+=("$path")
    png_names+=("$name")
    png_urls+=("")
  done < <(find "$screenshot_dir" -maxdepth 1 -type f -name '*.png' | LC_ALL=C sort)
}

# Prints the asset URL on stdout. Writes the HTTP status to upload_http_file
# (the call runs in a command substitution, so a shell variable would not
# be visible to the caller). Logs go to stderr.
upload_http_file=""
upload_one() {
  local repo_id="$1" file="$2" name="$3"
  local out http url
  http="000"
  out="$(mktemp)"
  http="$(
    curl -sS --max-time 60 -o "$out" -w '%{http_code}' \
      -X POST \
      -H "Authorization: Bearer ${token}" \
      -H "Accept: application/json" \
      -H "Content-Type: image/png" \
      --data-binary @"${file}" \
      "https://uploads.github.com/user-attachments/assets?name=${name}&content_type=image%2Fpng&repository_id=${repo_id}"
  )" || http="000"
  if [[ -n "$upload_http_file" ]]; then
    printf '%s\n' "$http" > "$upload_http_file"
  fi
  if [[ "$http" != "201" ]]; then
    local detail
    detail="$(tr '\n' ' ' < "$out" | cut -c1-180)"
    echo "upload ${name} failed: HTTP ${http} ${detail}" >&2
    rm -f "$out"
    return 1
  fi
  if ! command -v jq >/dev/null 2>&1; then
    echo "jq not found; cannot read upload response" >&2
    rm -f "$out"
    return 1
  fi
  url="$(jq -r '.url // empty' "$out")"
  rm -f "$out"
  case "$url" in
    https://github.com/*) ;;
    *)
      echo "upload ${name} returned an unexpected url" >&2
      return 1
      ;;
  esac
  printf '%s\n' "$url"
}

try_embed_uploads() {
  local n i j url repo_id any status
  n=${#png_paths[@]}
  if [[ "$n" -eq 0 || "$dry_run" -eq 1 || "${REPORT_SKIP_IMAGE_UPLOAD:-}" == "1" ]]; then
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1 || ! command -v jq >/dev/null 2>&1; then
    echo "curl or jq missing; skipping inline images" >&2
    return 0
  fi
  if ! repo_id="$(gh api "repos/${repo}" --jq .id)"; then
    echo "could not read repository id; skipping inline images" >&2
    return 0
  fi
  case "$repo_id" in
    ''|*[!0-9]*)
      echo "unexpected repository id; skipping inline images" >&2
      return 0
      ;;
  esac
  upload_http_file="$(mktemp)"
  for ((i=0; i<n; i++)); do
    case "${png_names[$i]}" in
      ''|*[!A-Za-z0-9._-]*)
        echo "unusual filename, not uploading ${png_names[$i]}" >&2
        continue
        ;;
    esac
    if url="$(upload_one "$repo_id" "${png_paths[$i]}" "${png_names[$i]}")"; then
      png_urls[i]="$url"
      continue
    fi
    status="$(tr -d '[:space:]' < "$upload_http_file")"
    # Installation tokens are rejected on the first request. Stop before
    # walking the rest of the directory when nothing has uploaded.
    if [[ "$status" =~ ^(401|403|404|000)$ ]]; then
      any=0
      for ((j=0; j<i; j++)); do
        if [[ -n "${png_urls[$j]}" ]]; then
          any=1
        fi
      done
      if [[ "$any" -eq 0 ]]; then
        echo "inline image upload is unavailable to this token (HTTP ${status}); listing files instead" >&2
        rm -f "$upload_http_file"
        upload_http_file=""
        return 0
      fi
    fi
    echo "leaving ${png_names[$i]} as a filename only" >&2
  done
  rm -f "$upload_http_file"
  upload_http_file=""
}

resolve_artifact_url() {
  artifact_url=""
  if [[ -n "${REPORT_ARTIFACT_URL:-}" ]]; then
    artifact_url="$REPORT_ARTIFACT_URL"
    return 0
  fi
  if [[ "$dry_run" -eq 1 ]]; then
    return 0
  fi
  local id
  id="$(
    gh api "repos/${repo}/actions/runs/${run_id}/artifacts?per_page=100" \
      --jq "[.artifacts[] | select(.name==\"${ARTIFACT_NAME}\" and .expired==false)] | max_by(.id) | .id // empty"
  )" || {
    echo "artifact lookup failed; the comment will name ${ARTIFACT_NAME} without a direct link" >&2
    return 0
  }
  if [[ -z "$id" || "$id" == "null" ]]; then
    if [[ ${#png_names[@]} -eq 0 ]]; then
      return 0
    fi
    sleep "${REPORT_ARTIFACT_RETRY_SLEEP:-3}"
    id="$(
      gh api "repos/${repo}/actions/runs/${run_id}/artifacts?per_page=100" \
        --jq "[.artifacts[] | select(.name==\"${ARTIFACT_NAME}\" and .expired==false)] | max_by(.id) | .id // empty"
    )" || true
  fi
  case "$id" in
    ''|null) return 0 ;;
    *[!0-9]*)
      echo "unexpected artifact id" >&2
      return 0
      ;;
  esac
  artifact_url="${server}/${repo}/actions/runs/${run_id}/artifacts/${id}"
}

theme_of() {
  case "$1" in
    *-light.png) printf '%s\n' light ;;
    *-dark.png) printf '%s\n' dark ;;
    *) printf '%s\n' other ;;
  esac
}

embedded_count() {
  local i n count url
  n=${#png_names[@]}
  count=0
  for ((i=0; i<n; i++)); do
    url="${png_urls[$i]:-}"
    if [[ -n "$url" ]]; then
      count=$((count + 1))
    fi
  done
  printf '%s\n' "$count"
}

append_images_for_theme() {
  local theme="$1"
  local i n name url found
  n=${#png_names[@]}
  found=0
  for ((i=0; i<n; i++)); do
    name="${png_names[$i]}"
    url="${png_urls[$i]:-}"
    [[ -n "$url" ]] || continue
    [[ "$(theme_of "$name")" == "$theme" ]] || continue
    if [[ "$found" -eq 0 ]]; then
      echo
      echo "#### ${theme}"
      echo
    fi
    found=1
    echo "![${name}](${url})"
    echo
  done
}

append_names_for() {
  local mode="$1"
  local i n name url
  n=${#png_names[@]}
  for ((i=0; i<n; i++)); do
    name="${png_names[$i]}"
    url="${png_urls[$i]:-}"
    if [[ "$mode" == "all" || -z "$url" ]]; then
      echo "- \`${name}\`"
    fi
  done
}

build_body() {
  local short embedded n marker
  marker="<!-- clash-go-screenshot-report run=${run_id} -->"
  n=${#png_names[@]}
  embedded="$(embedded_count)"
  short="$sha"
  if [[ ${#sha} -ge 7 ]]; then
    short="${sha:0:7}"
  fi

  echo "${marker}"
  echo
  echo "## CI 窗口截图 · ${status}"
  echo
  echo "- **结果**：${status}"
  if [[ -n "$test_outcome" ]]; then
    echo "- **test 步骤**：${test_outcome}"
  fi
  echo "- **事件**：${event_name}"
  echo "- **分支**：\`${branch}\`"
  if [[ -n "$repo" && -n "$sha" && "$sha" != "unknown" ]]; then
    echo "- **SHA**：[\`${short}\`](${server}/${repo}/commit/${sha})"
  else
    echo "- **SHA**：\`${sha}\`"
  fi
  if [[ -n "$run_url" ]]; then
    echo "- **Actions**：[run](${run_url})（attempt ${attempt}）"
  fi
  if [[ -n "$pr_url" && -n "$pr_number" ]]; then
    echo "- **PR**：[#${pr_number}](${pr_url})"
  elif [[ -n "$pr_url" ]]; then
    echo "- **PR**：${pr_url}"
  fi
  echo
  echo "### 截图"
  echo
  if [[ "$n" -eq 0 ]]; then
    echo "这次没有 \`screenshots/*.png\`。有文件时在 artifact \`${ARTIFACT_NAME}\`。"
    return 0
  fi
  if [[ "$embedded" -eq 0 ]]; then
    echo "共 ${n} 张，在 artifact \`${ARTIFACT_NAME}\`。"
  elif [[ "$embedded" -eq "$n" ]]; then
    echo "共 ${n} 张，嵌在下面。原文件也在 artifact \`${ARTIFACT_NAME}\`。"
  else
    echo "共 ${n} 张，嵌进了 ${embedded} 张。其余文件名列在后面。原文件在 artifact \`${ARTIFACT_NAME}\`。"
  fi
  if [[ -n "$artifact_url" ]]; then
    echo
    echo "- [打开 ${ARTIFACT_NAME}](${artifact_url})"
  elif [[ -n "$run_url" ]]; then
    echo
    echo "- 在上面的 Actions run 页面打开 artifact \`${ARTIFACT_NAME}\`"
  fi
  if [[ "$embedded" -eq 0 ]]; then
    echo
    append_names_for all
    return 0
  fi
  append_images_for_theme light
  append_images_for_theme dark
  append_images_for_theme other
  if [[ "$embedded" -ne "$n" ]]; then
    echo "#### 未嵌入"
    echo
    append_names_for missing
  fi
}

post_or_update() {
  local body_file="$1"
  local marker existing_id
  marker="<!-- clash-go-screenshot-report run=${run_id} -->"
  existing_id="$(
    gh api --paginate \
      "repos/${repo}/issues/${ISSUE}/comments?per_page=100" \
      --jq ".[] | select(.body | contains(\"<!-- clash-go-screenshot-report run=${run_id} -->\")) | .id" \
      | awk 'NF { id=$0 } END { print id }'
  )"
  if [[ -n "$existing_id" ]]; then
    case "$existing_id" in
      *[!0-9]*)
        echo "unexpected comment id: ${existing_id}" >&2
        exit 1
        ;;
    esac
    jq -n --rawfile body "$body_file" '{body:$body}' \
      | gh api --method PATCH \
        "repos/${repo}/issues/comments/${existing_id}" \
        --input - \
        --jq .html_url
    echo "Updated issue #${ISSUE} comment ${existing_id} (${marker})" >&2
    return 0
  fi
  gh issue comment "$ISSUE" --repo "$repo" --body-file "$body_file"
}

main() {
  if [[ "${1:-}" == "--dry-run" ]]; then
    dry_run=1
  elif [[ $# -gt 0 ]]; then
    usage
    exit 2
  fi

  repo="${REPORT_REPOSITORY:-${GITHUB_REPOSITORY:-}}"
  run_id="${REPORT_RUN_ID:-${GITHUB_RUN_ID:-}}"
  server="${REPORT_SERVER_URL:-https://github.com}"
  server="${server%/}"
  status="${REPORT_JOB_STATUS:-unknown}"
  event_name="${REPORT_EVENT_NAME:-${GITHUB_EVENT_NAME:-unknown}}"
  branch="${REPORT_BRANCH:-unknown}"
  sha="${REPORT_SHA:-unknown}"
  run_url="${REPORT_RUN_URL:-}"
  attempt="${REPORT_RUN_ATTEMPT:-${GITHUB_RUN_ATTEMPT:-1}}"
  pr_url="${REPORT_PR_URL:-}"
  pr_number="${REPORT_PR_NUMBER:-}"
  test_outcome="${REPORT_TEST_OUTCOME:-}"
  screenshot_dir="${SCREENSHOT_DIR:-screenshots}"

  if [[ "$pr_url" == "null" ]]; then pr_url=""; fi
  if [[ "$pr_number" == "null" ]]; then pr_number=""; fi

  if [[ "$dry_run" -eq 0 ]]; then
    token="${GH_TOKEN:-${GITHUB_TOKEN:-}}"
    if [[ -z "$token" || -z "$repo" || -z "$run_id" ]]; then
      echo "GH_TOKEN, repository, and run id are required" >&2
      exit 1
    fi
    if [[ ! "$run_id" =~ ^[0-9]+$ ]]; then
      echo "run id must be numeric" >&2
      exit 1
    fi
    if [[ ! "$repo" =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]]; then
      echo "unexpected repository name" >&2
      exit 1
    fi
    export GH_TOKEN="$token"
    export GH_PROMPT_DISABLED=1
    require_cmd gh
    require_cmd jq
  fi

  collect_pngs
  try_embed_uploads
  resolve_artifact_url

  local body_file
  body_file="$(mktemp)"
  build_body > "$body_file"
  if [[ "$dry_run" -eq 1 ]]; then
    cat "$body_file"
    rm -f "$body_file"
    return 0
  fi
  post_or_update "$body_file"
  rm -f "$body_file"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  main "$@"
fi
