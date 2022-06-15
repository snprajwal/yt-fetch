# Usage
1. Generate a YouTube API key by following the instructions [here](https://developers.google.com/youtube/v3/getting-started#before-you-start)
2. Export it in your shell with `export YT_API_KEY=your_key`
3. Run `docker compose up`
4. Access the server at `localhost:1234`

# API Reference
## GET `/ping`
- Response: `pong`

## GET `/`
- Query parameters:
    - `pageId`(*optional*): ID to fetch next page, returned in the response of the previous request.
- Response:

        {
            "items": [
                {
                    ...
                }
                {
                    ...
                }
            ],
            "nextPageId": "gY"
        }

## GET `/search`
- Query parameters:
    - `query`(*optional*): Search term to use. If empty, then generates same response as `/`.
    - `pageId`(*optional*): ID to fetch next page, returned in the response of the previous request.
- Response:

        {
            "items": [
                {
                    ...
                }
                {
                    ...
                }
            ],
            "nextPageId": "wR"
        }

# Features
- Keyset pagination for consistency and high performance through indexes
- Hashing pagination key from `int` to `string` to prevent business intelligence leaks (see [here](https://medium.com/lightrail/prevent-business-intelligence-leaks-by-using-uuids-instead-of-database-ids-on-urls-and-in-apis-17f15669fd2e) for more info)
- Pattern matching to get search results with partial match on search query
- [gorilla/mux](https://github.com/gorilla/mux) as the default router due to ease of configuration and high scalability
- Idiomatic functions which return `error`, with top-level error handling on the original calling function
- Diagnostics and logging to `stderr` by default
