# Newsletter

Some skilled engineers even have a blog site where they push some gold content, but they doesn't have yet, a way for their fan base to have recurrent access to this content. Newsletter try to circumvent it, scraping their pages and triggering e-mails for the guys who have an interest in those hidden gems.

![newsletter](./assets/newsletter.png)



# Roadmap

This program aims to create the following features:

- Given a list of websites, that are located in a MongoDB collection, scrape the content of each website and save it in another MongoDB collection. ✅
- After the scraping, calculate the similarity between the new content and the previous content of each website, and update the MongoDB collection
with this information. ✅
- All the registered users will receive an email according to the URL that they have registered notifying them about news in their favorite engineers websites. ✅

Obs: All these flows will be trigerred by a cron job. ✅


## Environement Variables

The following environment variables are required to run the program:

- `LOG_LEVEL`: The level of the logs that will be printed. The values could be `DEBUG`, `INFO`, `WARNING` or `ERROR`.
- `LOG_TYPE`: The format of the logs that will be printed. The values could be `json` or `text`.
- `NL_MONGO_URI`: The URI of the MongoDB database that will be used to store the data.
- `NL_EMAIL_PASSWORD`: The password of the email that will be used to send the emails.
- `NL_EMAIL_USERNAME`: The user of the email that will be used to send the emails.

## Commands

`make help` - Show the available commands of this project. Using it, it's enoght to play around the project.


## Integration & Unit Tests

To run the integration tests, you need to have a MongoDB instance running in your machine. To do it, you can run the following command:

```bash
    make dev/start
```

Access the dev container and run the tests:

```bash
    make dev
    make test 
    make integration-test
```


