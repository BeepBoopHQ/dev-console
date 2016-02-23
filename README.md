# WIP - Not ready for use.

## [Beep Boop HQ](http://beepboophq.com) dev-console
Simulate running on the BeepBoop hosting platform. Useful for testing multiple team scenarios.

## Use

1. Download the appropriate executable from `/dist`

2. Run it: `./dev-console`. You should see output indicating the server is running.

3. Browser to `http://localhost:9000`

4. Each row in the form represents a bot instance. The "Config" key value pairs represent configuration that is injected into your bot process through the environment variables you specify in your bot.yml.

  In the case of a Slack bot, if need to connect to an actual team/bot, set a SLACK_TOKEN configuration value using a valid token (can be retrieved from the "Integrations" portion of the Slack website).

&lt;animated gif showing use&gt;
