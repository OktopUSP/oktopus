#!/usr/bin/env sh
set -Ex

function apply_path {

    echo "Check that we have NEXT_PUBLIC_REST_ENDPOINT vars"
    test -n "$NEXT_PUBLIC_REST_ENDPOINT"

    echo "Check that we have NEXT_PUBLIC_WS_ENDPOINT vars"
    test -n "$NEXT_PUBLIC_WS_ENDPOINT"

    find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#REST_API_URL#$NEXT_PUBLIC_REST_ENDPOINT#g"
    find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#WS_URL#$NEXT_PUBLIC_WS_ENDPOINT#g"
}

apply_path
echo "Starting Nextjs"
exec "$@"
