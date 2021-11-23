# Changelog

### [Current Changes]
* Add `active` parameter to `ResumeChat` and `StartChat`.
* Renamed `ActivateChat` to `ResumeChat`.
* Removed `BotStatus` parameter from `CreateBot` and `UpdateBot` methods.
* `BotAgent` and `BotAgentDetails` merged to `Bot` structure.
* Add `fields` parameter to `ListBots` method in order to get additional Bots information.
* Add `fields` parameter to `GetBot` method in order to get additional Bot information.
* Add `type` parameter to `RegisterWebhook` ("bot" or "license").
* Renamed `EnableWebhooks` to `EnableLicenseWebhooks`.
* Renamed `DisableWebhooks` to `DisableLicenseWebhooks`.
* Renamed `GetWebhooksState` to `GetLicenseWebhooksState`.
* Support LiveChat APIs v3.3.
* Add `default_value` parameter to `RegisterProperty`.
* Expanded `filters.surveys` parameter with `values`/`exclude_values`, `from`/`to`, `groups` and `exclude` in `ListArchives`.

### [v2.1.0]

* Added support for setting custom `AuthorID` (ie. for messages sent by bot).
* Added `Target` property to `RichMessageButton`.
* Added `Type` to `authorization.Token` in order to support different authentication schemes (`Bearer` and `Basic`).
* Added possibility for chat transfer within the current group.
* Fixed marshaling of `Avatar` in `CreateBot`.
* Deprecated setting of `Status` via `CreateBot` and `UpdateBot` - use `SetRoutingStatus` instead.

### [v2.0.0]

* Support LiveChat APIs v3.2.

### [v1.0.0]

* Support LiveChat APIs v3.1.
