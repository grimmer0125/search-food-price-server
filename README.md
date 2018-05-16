# food-price-search-engine

## set up Slack bot to notify administrator the price alert

1. follow
https://get.slack.help/hc/en-us/articles/115005265703-Create-a-bot-for-your-workspace
2. get `API token` from its configuration page
3. set up a slack channel and get its id to be notified, how to get a slack channel id: https://stackoverflow.com/a/44883343/7354486
4. create `.env` file including the following information
    ```
    SLACK_TOKEN=YOUR_SLACK_TOKEN
    SLACK_CHANNEL=SLACK_CHANNEL_ID
    ```
