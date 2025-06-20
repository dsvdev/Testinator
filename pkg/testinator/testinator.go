package testinator

import (
	"encoding/json"
	"github.com/dsvdev/Testinator/internal/cfg"
	"github.com/dsvdev/Testinator/internal/engine"
	internal_model "github.com/dsvdev/Testinator/internal/model"
	llm_model "github.com/dsvdev/Testinator/internal/model/llm"
	"github.com/dsvdev/Testinator/internal/prompt"
	"github.com/dsvdev/Testinator/pkg/llm"
	"github.com/dsvdev/Testinator/pkg/model"
	"log"
	"os"
	"strings"
)

type Testinator struct {
	driver  llm.LLMDriver
	engine  *engine.TestEngine
	openApi string
}

func NewTestinator(driver llm.LLMDriver) *Testinator {
	c := cfg.TestinatorConfig
	return &Testinator{
		driver:  driver,
		engine:  engine.NewEngine(c.AppURL),
		openApi: "",
	}
}

func (t *Testinator) WithOpenApi(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	t.openApi = string(data)
}

func (t *Testinator) OpenApi() string {
	return t.openApi
}

func (t *Testinator) ExecuteTestCase(testCase *model.TestCase) (model.TestResult, error) {
	context := make(internal_model.TestExecutionContext, 0)
	steps := testCase.Steps

	for _, step := range steps {
		log.Printf("Executing step: %v", step)
		newContext, err := t.executeStep(step, context)
		if err != nil {
			return "", err
		}
		context = newContext
	}

	return "success", nil
}

func (t *Testinator) executeStep(step model.TestStep, executionContext internal_model.TestExecutionContext) (internal_model.TestExecutionContext, error) {
	generatedPrompt := prompt.Generate(step, executionContext, t.openApi)
	log.Printf("Generated prompt: %s", generatedPrompt)

	response, err := t.driver.SendRequest(generatedPrompt)
	log.Printf("LLM response: %v", response)
	if strings.HasPrefix(response, "```json") {
		strings.TrimPrefix(response, "```json`")
	}
	if strings.HasSuffix(response, "```") {
		strings.TrimSuffix(response, "```")
	}

	if err != nil {
		return nil, err
	}

	var llmResponse llm_model.LLMResponse
	err = json.Unmarshal([]byte(response), &llmResponse)
	if err != nil {
		return nil, err
	}
	executionContext = append(executionContext, "Step: "+string(step))
	executionContext = append(executionContext, "Action: "+response)

	return t.engine.ExecuteLlmResponse(llmResponse, executionContext)
}
