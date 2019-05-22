package serviceprovider

import (
	"fmt"
	"testing"
)

func TestServiceProviderTypeString(t *testing.T) {
	testcases := map[string]struct {
		typ ServiceProviderType
		str string
	}{
		"PROVIDER": {
			typ: PROVIDER,
			str: "PROVIDER",
		},
		"SERVICE": {
			typ: SERVICE,
			str: "SERVICE",
		},
		"": {
			typ: -1,
			str: "",
		},
	}
	for n, tc := range testcases {
		t.Run(n, func(t *testing.T) {
			if fmt.Sprint(tc.typ) != tc.str {
				t.Errorf("wrong output: %v", tc.typ)
			}
		})
	}
}
