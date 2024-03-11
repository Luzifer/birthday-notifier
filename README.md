# Luzifer / birthday-notifier

Previously I used an app on my phone to notify me on upcoming birthdays: One day in advance and on the same day. Then that app decided to add a skew of one day and notify me too late…

And that's why there is now a server based solution to notifying my of upcoming birthdays.

Features:

- Sync contacts using CardDAV
- Extract birthdays from the contacts
- Send notifications

Hosted somewhere it's always running and configured properly and birthday notifications are coming in properly.

## Usage

```console
# birthday-notifier --help
Usage of birthday-notifier:
      --fetch-interval duration       How often to fetch birthdays from CardDAV (default 1h0m0s)
      --log-level string              Log level (debug, info, warn, error, fatal) (default "info")
      --notify-days-in-advance ints   Send notification X days before birthday (default [1])
      --notify-via string             How to send the notification (one of: log, pushover) (default "log")
      --version                       Prints current version and exits
      --webdav-base-url string        Webdav server to connect to
      --webdav-pass string            Password for the Webdav user
      --webdav-principal string       Principal format to fetch the addressbooks for (%s will be replaced with the webdav-user) (default "principals/users/%s")
      --webdav-user string            Username for Webdav login
```

For Nextcloud leave the principal format the default, for other systems you might need to adjust it.

To adjust the notification text see the template in [`pkg/formatter/formatter.go`](./pkg/formatter/formatter.go) and provide your own as `NOTIFICATION_TEMPLATE` environment variable.

### Notifier configuration

- **`log`** - Just sends the notification to the console logs, no configuration available
- **`pushover`** - Send notification via [Pushover](https://pushover.net)
  - `PUSHOVER_API_TOKEN` - Token for the App you've created in the Pushover Dashboard
  - `PUSHOVER_USER_KEY` - Token for the User to send the notification to
  - `PUSHOVER_SOUND` - (Optional) Specify a sound to use
