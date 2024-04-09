#!/usr/bin/env sh
set -Ex

function apply_path {

    echo "Check that we have NEXT_PUBLIC_REST_ENPOINT vars"
    test -n "$NEXT_PUBLIC_REST_ENPOINT"

    echo "Check that we have NEXT_PUBLIC_WS_ENPOINT vars"
    test -n "$NEXT_PUBLIC_WS_ENPOINT"

    find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#REST_API_URL#$NEXT_PUBLIC_REST_ENPOINT#g"
    find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#WS_URL#$NEXT_PUBLIC_WS_ENPOINT#g"
}

apply_path
echo "Starting Nextjs"
exec "$@"