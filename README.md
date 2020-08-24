# LiveChat Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/livechat/lc-sdk-go/v2)](https://goreportcard.com/report/github.com/livechat/lc-sdk-go/v2)
[![GoDoc](https://godoc.org/github.com/livechat/lc-sdk-go/v2?status.svg)](https://godoc.org/github.com/livechat/lc-sdk-go/v2)

This Software Development Kit written in [Go](https://go.dev/) helps developers build external backend apps that extend LiveChat features. The SDK makes it easy to use [Agent Chat API](https://developers.livechatinc.com/docs/messaging/agent-chat-api/), [Customer Chat API](https://developers.livechatinc.com/docs/messaging/customer-chat-api/) and [Configuration API](https://developers.livechatinc.com/docs/management/configuration-api/).

### Technical docs

For technical documentation ([godoc](https://godoc.org/) format), please go to [LiveChat SDK Docs](https://godoc.org/github.com/livechat/lc-sdk-go/v2).

### API protocol docs

For protocol documentation of LiveChat APIs, please go to [Livechat Platform Docs](https://developers.livechatinc.com/docs/).

### Go modules vs API version

LiveChat Go SDK supports Go modules. Please note that minor LiveChat API versions can be incompatible. Here is the relation:
* lc-sdk-go 1.x.x -> LiveChat API 3.1
* lc-sdk-go 2.x.x -> LiveChat API 3.2
* ...

All versions of LiveChat API are available as git tags in lc-sdk-go. However, a developer-preview version (not completed yet, may introduce breaking changes in future) is avaiable in lc-sdk-go as a git branch.

### Usage guide and examples

* [Echo](examples/echo/README.md)

### Feedback

​If you find any bugs or have trouble implementing the code on your own, please create an issue or contact us [LiveChat for Developers](https://developers.livechatinc.com/).

### About LiveChat

LiveChat is an online customer service software with live support, help desk software, and web analytics capabilities. It's used by more than 30,000 companies all over the world. For more info, check out [LiveChat](https://livechat.com/).
