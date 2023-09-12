package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	appparams "github.com/jmesworld/core/v2/app/params"
	"github.com/jmesworld/core/v2/x/tokenfactory/types"
)

func TestDeconstructDenom(t *testing.T) {
	appparams.RegisterAddressesConfig()

	for _, tc := range []struct {
		desc             string
		denom            string
		expectedSubdenom string
		err              error
	}{
		{
			desc:  "empty is invalid",
			denom: "",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "normal",
			denom:            "factory/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
			expectedSubdenom: "bitcoin",
		},
		{
			desc:             "multiple slashes in subdenom",
			denom:            "factory/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin/1",
			expectedSubdenom: "bitcoin/1",
		},
		{
			desc:             "no subdenom",
			denom:            "factory/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/",
			expectedSubdenom: "",
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
			err:   types.ErrInvalidDenom,
		},
		{
			desc:             "subdenom of only slashes",
			denom:            "factory/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/////",
			expectedSubdenom: "////",
		},
		{
			desc:  "too long name",
			denom: "factory/jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			err:   types.ErrInvalidDenom,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			expectedCreator := "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g"
			creator, subdenom, err := types.DeconstructDenom(tc.denom)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedCreator, creator)
				require.Equal(t, tc.expectedSubdenom, subdenom)
			}
		})
	}
}

func TestGetTokenDenom(t *testing.T) {
	appparams.RegisterAddressesConfig()
	for _, tc := range []struct {
		desc     string
		creator  string
		subdenom string
		valid    bool
	}{
		{
			desc:     "normal",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "bitcoin",
			valid:    true,
		},
		{
			desc:     "multiple slashes in subdenom",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "bitcoin/1",
			valid:    true,
		},
		{
			desc:     "no subdenom",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "",
			valid:    true,
		},
		{
			desc:     "subdenom of only slashes",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "/////",
			valid:    true,
		},
		{
			desc:     "too long name",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			valid:    false,
		},
		{
			desc:     "subdenom is exactly max length",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
			subdenom: "bitcoinfsadfsdfeadfsafwefsefsefsdfsdafasefsf",
			valid:    true,
		},
		{
			desc:     "creator is exactly max length",
			creator:  "jmes19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2gjhgjhgkhjklhkjhkjhgjhgjgjghelug",
			subdenom: "bitcoin",
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := types.GetTokenDenom(tc.creator, tc.subdenom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
