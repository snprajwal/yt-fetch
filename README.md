# Usage
1. Generate a YouTube API Key by follwing the instructions [here](https://developers.google.com/youtube/v3/getting-started#before-you-start)
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
        ]
        "nextPageId":123
    }

## GET `/search`
- Query parameters:
    - `query`(*optional*): Search term to use. If empty, then generates same response as `/`.
- Response:
    {
        "items": [
            {
                ...
            }
            {
                ...
            }
        ]
        "nextPageId":123
    }
