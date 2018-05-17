# food-price-search-engine

## Development

**only tested on Mac 10.13.2 + go 1.10.2**

### Requirement

1. go >= 1.9 (since golang/dep is used)

### (optional) Set up Slack bot to notify administrator the price alert

This is a workdaround way instead of having a admin site. Customizing a site or using GitHub issues is possible. 

1. follow
https://get.slack.help/hc/en-us/articles/115005265703-Create-a-bot-for-your-workspace
2. get `API token` from its configuration page
3. set up a slack channel and get its id to be notified, how to get a slack channel id: https://stackoverflow.com/a/44883343/7354486
4. create `.env` file including the following information
    ```
    SLACK_TOKEN=YOUR_SLACK_TOKEN
    SLACK_CHANNEL=SLACK_CHANNEL_ID
    ```
### (optional) Launch Dockerized MongoDB server as a cache storage

```
docker run --name some-mongo -p 27017:27017 -d mongo:3.7.9-jessie
```

This project also uses MongoDB's TTL (time to live) to expire data

### Install Go Dependencies

1. Install [`golang/dep`](https://github.com/golang/dep) if you do not have.
2. In the project root folder, `dep ensure` to get all packages put in `vendor`

### Launch

`go build`, then execute `search-food-price-server`

### Debug

Use VS Code to debug it.  

## Deployment

Not yet. The features are not too much. `Kubernetes` may be used.

## Issues & ToDo List
1. if MongoDB is not launched, sometimes it will wait for some secondes when tyrining to use it.
