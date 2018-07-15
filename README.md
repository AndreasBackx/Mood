# Mood

Mood is a moody moderator bot for simplifying the moderation of Slack teams.

_Note: This is very much a work in progress, as this was written in an evening as the time of writing this._

## Features

- [x] Action menu to report spam/crossposts
    - Report is sent to a private moderation channel
    - The reported user is sent a message where the user is asked to remove the message or say whether the report was invalid.
    - The response from the reported user is updated in the report sent to the private moderation channel.
    - A message is also shown to the user who reported the post confirming that the report was (un)successful.
- [x] Welcome new team members with the rules
    - [ ] Provide a `/rules` command as well.
    - [ ] Properly update the message when clicking the "Yes!" button.

## Configuration

To get started, install the dependencies using dep:

```shell
dep ensure
```

A `config.json` file is required in the root of the project:

```javascript
{
    // Verification token
    // Found under the Apps page -> Settings -> Basic Information
    "verification_token": "",
    // Bot USer OAuth Access Token
    // Found under the Apps page -> Features -> OAuth & Permissions
    "bot_user_oauth_token": "",
    // Name of the group where the bot needs to post the reports.
    // The bot needs to be invited to the group.
    "report_group_name": "reports",

    // Port used to run the server on.
    "port" 8000
}
```

In order to run, simply run the `realize start` command which will run all of the necessary commands. See [Realize](https://github.com/oxequa/realize) and [`.realize.yaml`](.realize.yaml) for more info.
