# Changelog

### [Current Changes]
* Add `active` parameter to `ResumeChat` and `StartChat`.
* Renamed `ActivateChat` to `ResumeChat`.
* Support LiveChat APIs v3.3.

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
