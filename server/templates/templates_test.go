package templates_test

import (
	"fmt"
	"testing"

	"github.com/jitsucom/jitsu/server/events"
	"github.com/jitsucom/jitsu/server/logging"
	"github.com/jitsucom/jitsu/server/script/deno"
	"github.com/jitsucom/jitsu/server/script/node"
	"github.com/jitsucom/jitsu/server/templates"
)

const Expression = `return {...$, hello: $.id}`

func benchmarkExecutor(b *testing.B, executor templates.TemplateExecutor) {
	b.ResetTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		result, err := executor.ProcessEvent(events.Event{"id": i})
		if err != nil {
			b.Fatalf("error: %+v", err)
		}

		if fmt.Sprint(i) != fmt.Sprint(result.(map[string]interface{})["hello"]) {
			b.Fatalf("%v must be %d", result, i)
		}
	}
}

func BenchmarkNodeStdIO(b *testing.B) {
	factory, err := node.Factory()
	if err != nil {
		b.Fatal(err)
	}

	logging.LogLevel = logging.INFO
	templates.SetScriptFactory(factory)
	process, err := templates.NewScriptExecutor(templates.Expression(Expression), nil)
	if err != nil {
		b.Fatal(err)
	}

	defer process.Close()
	benchmarkExecutor(b, process)
}

func BenchmarkDenoStdIO(b *testing.B) {
	logging.LogLevel = logging.INFO
	templates.SetScriptFactory(deno.Factory())
	executor, err := templates.NewScriptExecutor(templates.Expression(Expression), nil)
	if err != nil {
		b.Fatal(err)
	}

	defer executor.Close()
	benchmarkExecutor(b, executor)
}