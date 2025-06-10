package esub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/esub"
)

func Test_Eval_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tt := []struct {
		tmpl    string
		envList []string
		out     string
	}{
		{
			tmpl:    "foo:${foo} BAR:${BAR}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "foo:ZOO BAR:baz",
		},
		{
			tmpl:    "*** foo:${foo} *** BAR:${BAR} ***",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "*** foo:ZOO *** BAR:baz ***",
		},
		{
			tmpl:    "foo:$${foo} BAR:${BAR}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "foo:${foo} BAR:baz",
		},
		{
			tmpl:    "foo:${foo} BAR:$${BAR}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "foo:ZOO BAR:${BAR}",
		},
		{
			tmpl:    "foo:$${foo} BAR:$${BAR}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "foo:${foo} BAR:${BAR}",
		},
		{
			tmpl:    "foo:$$${foo} BAR:$${ ZOO:$$ baz:$} HOGE=$PIYO",
			envList: []string{"foo=ZOO", "BAR=baz"},
			out:     "foo:$${foo} BAR:${ ZOO:$$ baz:$} HOGE=$PIYO",
		},
	}

	for _, t := range tt {
		out, err := esub.Eval(t.tmpl, t.envList)
		require.NoError(err)
		assert.Equal(t.out, out)
	}
}

func Test_Fill_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Setenv("foo", "ooz")
	t.Setenv("BAR", "ZAB")
	out, err := esub.Fill("foo:${foo} BAR:${BAR}")
	require.NoError(err)
	assert.Equal("foo:ooz BAR:ZAB", out)
}

func Test_Eval_Err(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		tmpl    string
		envList []string
		err     string
	}{
		{
			tmpl:    "foo:${",
			envList: []string{"foo=ZOO", "BAR=baz"},
			err:     "syntax error: foo:${",
		},
		{
			tmpl:    "foo:${}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			err:     "syntax error: foo:${}",
		},
		{
			tmpl:    "foo:${foo",
			envList: []string{"foo=ZOO", "BAR=baz"},
			err:     "syntax error: foo:${foo",
		},
		{
			tmpl:    "foo:${FOO}",
			envList: []string{"foo=ZOO", "BAR=baz"},
			err:     "env 'FOO' not found: foo:${FOO}",
		},
	}

	for _, t := range tt {
		_, err := esub.Eval(t.tmpl, t.envList)
		assert.ErrorContains(err, t.err)
	}
}
