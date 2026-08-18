package volunteer_test

import (
	"testing"

	"github.com/janghanul090801/pico-backend/external/volunteer"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeCondition_MultipleAgeCandidatesGoToUnparsed(t *testing.T) {
	item := volunteer.DetailItem{
		ProgrmSj:      "영유아 특수학교 '화요일' 수업활동 보조 봉사(중식제공)",
		ProgrmCn:      "3세~6세까지의 발달지체 영유아... 20~40대 가능",
		YngbgsPosblAt: "N",
	}

	policy := volunteer.ToPolicy(item)
	// 나이 표현이 여러 개(3세~6세, 20~40대) 존재하므로 파서에서 정제된 min/max age 대신 UnparsedConditions로 문맥이 넘어감
	assert.Nil(t, policy.Condition.MaxAge)
	assert.NotNil(t, policy.Condition.MinAge)
	assert.Equal(t, 19, *policy.Condition.MinAge) // YngbgsPosblAt == "N" 성인 기준 최소 age 설정
	assert.NotEmpty(t, policy.Condition.UnparsedConditions)
}
