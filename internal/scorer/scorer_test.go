package scorer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScorer_CalculateOverallScore(t *testing.T) {
	// Red: 실패하는 테스트 작성
	metrics := Metrics{
		Structure:            80.0,
		Conciseness:          70.0,
		InstructionFollowing: 90.0,
		SelfConsistency:      85.0,
		LatencyCost:          75.0,
		Risk:                 95.0,
	}
	weights := DefaultWeights()
	scorer := NewScorer()

	result := scorer.Score(metrics, weights)

	require.NotNil(t, result)
	assert.Greater(t, result.OverallScore, 0.0)
	assert.LessOrEqual(t, result.OverallScore, 100.0)
}

func TestScorer_WeightedCalculation(t *testing.T) {
	// 모든 메트릭이 100일 때 총점도 100이어야 함
	metrics := Metrics{
		Structure:            100.0,
		Conciseness:          100.0,
		InstructionFollowing: 100.0,
		SelfConsistency:      100.0,
		LatencyCost:          100.0,
		Risk:                 100.0,
	}
	weights := DefaultWeights()
	scorer := NewScorer()

	result := scorer.Score(metrics, weights)

	assert.Equal(t, 100.0, result.OverallScore)
}

func TestScorer_CustomWeights(t *testing.T) {
	metrics := Metrics{
		Structure:            100.0,
		Conciseness:          0.0,
		InstructionFollowing: 0.0,
		SelfConsistency:      0.0,
		LatencyCost:          0.0,
		Risk:                 0.0,
	}
	weights := Weights{
		Structure:            1.0, // 100% 가중치
		Conciseness:          0.0,
		InstructionFollowing: 0.0,
		SelfConsistency:      0.0,
		LatencyCost:          0.0,
		Risk:                 0.0,
	}
	scorer := NewScorer()

	result := scorer.Score(metrics, weights)

	assert.Equal(t, 100.0, result.OverallScore)
}

func TestScorer_ZeroMetrics(t *testing.T) {
	metrics := Metrics{
		Structure:            0.0,
		Conciseness:          0.0,
		InstructionFollowing: 0.0,
		SelfConsistency:      0.0,
		LatencyCost:          0.0,
		Risk:                 0.0,
	}
	weights := DefaultWeights()
	scorer := NewScorer()

	result := scorer.Score(metrics, weights)

	assert.Equal(t, 0.0, result.OverallScore)
}

func TestScorer_PreservesMetrics(t *testing.T) {
	metrics := Metrics{
		Structure:            50.0,
		Conciseness:          60.0,
		InstructionFollowing: 70.0,
		SelfConsistency:      80.0,
		LatencyCost:          90.0,
		Risk:                 95.0,
	}
	weights := DefaultWeights()
	scorer := NewScorer()

	result := scorer.Score(metrics, weights)

	assert.Equal(t, metrics.Structure, result.Metrics.Structure)
	assert.Equal(t, metrics.Conciseness, result.Metrics.Conciseness)
	assert.Equal(t, metrics.InstructionFollowing, result.Metrics.InstructionFollowing)
}
