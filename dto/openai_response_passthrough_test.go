package dto

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAITextResponsePreservesReasoningAndExtraFields(t *testing.T) {
	raw := []byte(`{
		"id":"chatcmpl-test",
		"object":"chat.completion",
		"created":123,
		"model":"gpt-5.4",
		"service_tier":"default",
		"choices":[{
			"index":0,
			"message":{
				"role":"assistant",
				"content":"answer",
				"reasoning_content":"private reasoning",
				"reasoning_details":[{"type":"summary","text":"brief"}]
			},
			"logprobs":{"content":[]},
			"finish_reason":"stop"
		}],
		"usage":{
			"prompt_tokens":1,
			"completion_tokens":2,
			"total_tokens":3,
			"completion_tokens_details":{"reasoning_tokens":1}
		}
	}`)

	var response OpenAITextResponse
	require.NoError(t, common.Unmarshal(raw, &response))

	response.Usage.PromptTokens = 4
	encoded, err := common.Marshal(response)
	require.NoError(t, err)

	require.Equal(t, "default", gjson.GetBytes(encoded, "service_tier").String())
	require.Equal(t, "private reasoning", gjson.GetBytes(encoded, "choices.0.message.reasoning_content").String())
	require.Equal(t, "summary", gjson.GetBytes(encoded, "choices.0.message.reasoning_details.0.type").String())
	require.True(t, gjson.GetBytes(encoded, "choices.0.logprobs").Exists())
	require.Equal(t, int64(1), gjson.GetBytes(encoded, "usage.completion_tokens_details.reasoning_tokens").Int())
}

func TestChatCompletionStreamDeltaPreservesExtraReasoningFields(t *testing.T) {
	raw := []byte(`{
		"id":"chatcmpl-test",
		"object":"chat.completion.chunk",
		"created":123,
		"model":"gpt-5.4",
		"choices":[{
			"index":0,
			"delta":{
				"reasoning_content":"step",
				"reasoning_details":[{"type":"summary","text":"brief"}]
			},
			"finish_reason":null
		}]
	}`)

	var response ChatCompletionsStreamResponse
	require.NoError(t, common.Unmarshal(raw, &response))

	encoded, err := common.Marshal(response)
	require.NoError(t, err)

	require.Equal(t, "step", gjson.GetBytes(encoded, "choices.0.delta.reasoning_content").String())
	require.Equal(t, "brief", gjson.GetBytes(encoded, "choices.0.delta.reasoning_details.0.text").String())
}
