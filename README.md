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
  -c, --config string      Configuration file path (default "config.yaml")
      --log-level string   Log level (debug, info, warn, error, fatal) (default "info")
      --version            Prints current version and exits
```

## Configuration

```yaml
# Specify days before the actual birthday to send advance notifications
# i.e. to buy gifts or something. Default is to send only on the actual
# birthday itself.
notifyDaysInAdvance: [ 1 ]

# Configure how to notify you when there is a birthday pending / today.
# Each entry consists of a type and the settings for that kind of
# notifier. For settings and available types see below.
notifiers:
  - type: slack
    settings:
      webhook: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

# Specify your own template for the notification text. The default is
# shown below and whould yield something like this:
#
# Ava has their birthday on Wed, 13 Mar. They are turning 27.
template: >-
  {{ .contact | getName }} has their birthday
  {{ if .when | isToday -}} today {{- else -}}
  on {{ (.when | projectToNext).Format "Mon, 02 Jan" }} {{- end }}.
  {{ if gt .when.Year 1 -}}They are turning {{ .when | getAge }}.{{- end }}

# Configure how to connect to the CardDAV addressbooks inside the
# webdav server
webdav:
  # Base-URL for the webdav server (example for Nextcloud)
  baseURL: https://my-nextcloud.example.com/remote.php/dav/
  # How often to fetch new birthdays (default: 1h)
  fetchInterval: 1h
  # Password for the user
  pass: 'my super secret password'
  # Principal format for the webdav server (default as below is valid
  # for Nextcloud instances): `%s` will be replaced with the value of
  # the user field below.
  principal: 'principals/users/%s'
  # Username for the login to the webdav server
  user: 'my.username'
```

### Notifiers

#### `log`

Just sends the notification to the console logs

```yaml
notifiers:
  - type: log
    # No settings for this one
```

#### `pushover`

Send notification via [Pushover](https://pushover.net)

```yaml
notifiers:
  - type: pushover
    settings:
      # Token for the App you've created in the Pushover Dashboard
      apiToken: '...' 
      # Token for the User to send the notification to
      userKey: '...'
      # (Optional) Specify a sound to use
      sound: ''
```

#### `slack`

Send notification through Slack(-compatible) webhook

```yaml
notifiers:
  - type: slack
    settings:
      # Webhook URL (i.e. `https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX`
      # or `https://discord.com/api/webhooks/00000/XXXXX/slack`)
      webhook: 'https://...'
      # (Optional) Specify the channel to send to
      channel: ''
      # (Optional) Emoji to use as user icon
      iconEmoji: ''
      # (Optional) Overwrite the hooks username\
      username: ''
```
