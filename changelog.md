# Changelog


### [v2.1.0]

* Added support for setting custom `AuthorID` (ie. for messages sent by bot).
* Added `Target` property to `RichMessageButton`.
* Added `Type` to `authorization.Token` in order to support different authentication schemes (`Bearer` and `Basic`).
* Added possibility for chat transfer within the current group.
* Fixed marshaling of `Avatar` in `CreateBot`.
* Deprecated setting of `Status` via `CreateBot` and `UpdateBot` - use `SetRoutingStatus` instead.\
