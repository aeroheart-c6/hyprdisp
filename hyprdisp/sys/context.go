package sys

type ContextKey string

const (
	ContextKeyLogger     ContextKey = "hyprdisp.logger"
	ContextKeyCLIActions ContextKey = "hyprdisp.actionRegistry"
)
