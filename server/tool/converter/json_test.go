package converter

import (
	"testing"
)

func TestParseJson(t *testing.T) {
	data := []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev",
	  "upvote_num":123,
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

	res := NewJsonConverter().GenStruct(string(data))

	t.Log("\n")
	t.Log(res)
}

func TestParseNumber(t *testing.T) {
	res := NewJsonConverter().GenStruct(`{"followers": 3.14}`)
	t.Log(res)
}

func TestParseTime(t *testing.T) {
	res := NewJsonConverter().GenStruct(`{"time":"2021-07-10T22:58:21.502494+08:00"}`)
	t.Log(res)
}

func BenchmarkConvert(b *testing.B) {
	data := []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

	for i := 0; i < b.N; i++ {
		c := NewJsonConverter().GenStruct(string(data))
		_ = c
	}
}

func Test_snakesToCamel(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			args: args{
				str: "hell_world",
			},
			wantResult: "hellWorld",
		},
		{
			args: args{
				str: "Hell_World",
			},
			wantResult: "hellWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := snakesToCamel(tt.args.str); gotResult != tt.wantResult {
				t.Errorf("snakesToCamel() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
