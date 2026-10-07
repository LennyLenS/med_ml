package mappers

import (
	api "composition-api/internal/generated/http/api"
)

func Contor(contor []byte) (api.ExamContor, error) {
	out := &api.ExamContor{}

	err := out.UnmarshalJSON(contor)
	if err != nil {
		return api.ExamContor{}, err
	}

	return *out, nil
}
