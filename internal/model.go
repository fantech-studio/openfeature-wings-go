package internal

type (
	EvalRequest struct {
		ID   string         `json:"id" validate:"required"`
		Meta map[string]any `json:"meta"`
	}

	EvalResponse struct {
		Variant string       `json:"variant"`
		Bool    *BoolValue   `json:"bool,omitempty"`
		Int     *IntValue    `json:"int,omitempty"`
		Float   *FloatValue  `json:"float,omitempty"`
		String  *StringValue `json:"string,omitempty"`
		Object  *ObjectValue `json:"object,omitempty"`
	}

	BoolValue struct {
		Value bool `json:"value"`
	}

	IntValue struct {
		Value int64 `json:"value"`
	}

	FloatValue struct {
		Value float64 `json:"value"`
	}

	StringValue struct {
		Value string `json:"value"`
	}

	ObjectValue struct {
		Value map[string]any `json:"value"`
	}
)
