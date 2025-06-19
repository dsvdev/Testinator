package llm

type LLMDriver interface {
	SendRequest(prompt string) (string, error)
}
