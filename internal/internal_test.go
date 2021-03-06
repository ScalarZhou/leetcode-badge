package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"strconv"
	"time"
)

func initEnv(t *testing.T) *assert.Assertions {
	Conf = &config{Cache: true, CacheTTL: time.Minute, Release: false}
	InitConfig()
	return assert.New(t)
}

func TestParseTmpl(t *testing.T) {
	as := initEnv(t)
	var data = make(map[string]interface{})

	data["solved_question"] = 1
	data["all_question"] = 2
	tmpl := `Leetcode | Solved/Total-{{.solved_question}}/{{.all_question}}-{{ if le .solved_question_rate_float 0.3}}red.svg{{ else if le .solved_question_rate_float 0.6}}yellow.svg{{ else }}green.svg{{ end }}`

	{
		data["solved_question_rate_float"] = 0.3
		x, err := parseTmpl(tmpl, data)
		as.Nil(err)
		as.Equal("Leetcode | Solved/Total-1/2-red.svg", string(x))
	}
	{
		data["solved_question_rate_float"] = 0.6
		x, err := parseTmpl(tmpl, data)
		as.Nil(err)
		as.Equal("Leetcode | Solved/Total-1/2-yellow.svg", string(x))
	}
	{
		data["solved_question_rate_float"] = 0.9
		x, err := parseTmpl(tmpl, data)
		as.Nil(err)
		as.Equal("Leetcode | Solved/Total-1/2-green.svg", string(x))
	}
}

func TestCache(t *testing.T) {
	as := initEnv(t)

	{
		key := "test" + strconv.Itoa(int(time.Now().Unix()))
		cacheSetLeetcode(key, &LeetcodeData{AllQuestion: 1000})
		r, ok := cacheGetLeetcode(key)
		as.True(ok)
		as.NotNil(r)
		as.Equal(1000, r.AllQuestion)
	}

	{
		key := "test" + strconv.Itoa(int(time.Now().Unix()))
		cacheSetShields(key, "test 2")
		r, ok := cacheGetShields(key)
		as.True(ok)
		as.Equal("test 2", r)
	}
}

func TestLeetcode(t *testing.T) {
	as := initEnv(t)
	var r *LeetcodeData
	var err error

	{
		r, err = fetchLeetcodeData("chyroc")
		as.Nil(err)
		as.NotNil(r)

		// 下面的数据是已经验证的
		as.True(r.SolvedQuestion >= 80)
		as.True(r.AllQuestion >= 500)
		as.NotEmpty(r.SolvedQuestionRate)
		as.True(r.SolvedQuestionRateFloat > 0.1)

		as.True(r.AcceptedSubmission >= 100)
		as.True(r.AllSubmission >= 200)
		as.NotEmpty(r.AcceptedSubmissionRate)
		as.True(r.AcceptedSubmissionRateFloat >= 0.1)
	}
}
