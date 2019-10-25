package configuration

type registerWebhookResponse struct {
	ID string `json:"webhook_id"`
}

type unregisterWebhookRequest struct {
	ID string `json:"webhook_id"`
}
