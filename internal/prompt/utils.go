package prompt

import (
	internal_model "github.com/dsvdev/Testinator/internal/model"
	"github.com/dsvdev/Testinator/pkg/model"
	"os"
	"strings"
)

var basePrompt string

func init() {
	data, err := os.ReadFile("/Users/dsv/dsvdev/testinator/code/Testinator/internal/prompt/basePrompt")
	if err != nil {
		panic(err)
	}
	basePrompt = string(data)
}

func Generate(step model.TestStep, executionContext internal_model.TestExecutionContext, openApi string) string {
	sb := strings.Builder{}
	sb.WriteString(basePrompt)
	sb.WriteString("\n")
	if len(openApi) != 0 {
		sb.WriteString("This opeanapi.yaml file with info about http endpoint of application:\n")
		sb.WriteString(openApi)
		sb.WriteString("\n")
	}
	sb.WriteString("user: |\n")
	sb.WriteString("Previous Steps:\n")
	for _, s := range executionContext {
		sb.WriteString(s)
		sb.WriteString("\n")
	}
	sb.WriteString("\nNext step:\n")
	sb.WriteString(string(step))
	return sb.String()
}
