# Newsletter

Some skilled engineers even have a blog site where they push some gold content, but they doesn't have yet, a way for their fan base to have recurrent access to this content. Newsletter try to circumvent it, scraping their pages and triggering e-mails for the guys who have an interest in those hidden gems.

![newsletter](./assets/newsletter.png)



# Roadmap

This program aims to create the following features:

- Given a list of websites, that are located in a MongoDB collection, scrape the content of each website and save it in another MongoDB collection.
- After the scraping, calculate the similarity between the new content and the previous content of each website, and update the MongoDB collection
with this information.
- All the registered users will receive an email according to the URL that they have registered notifying them about news in their favorite engineers websites.

Obs: All these flows will be trigerred by a cron job.
