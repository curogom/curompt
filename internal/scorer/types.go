package scorer

// Metrics contains individual metric scores
type Metrics struct {
	Structure            float64 // 0-100: 섹션 헤더, 중복 규칙, 금지어
	Conciseness          float64 // 0-100: 토큰 밀도, 압축 이득
	InstructionFollowing float64 // 0-100: JSONSchema/정규식 체크
	SelfConsistency      float64 // 0-100: 다중 샘플 간 일관성 (M2에서 구현)
	LatencyCost          float64 // 0-100: 토큰 수, 지연 시간 (M2에서 구현)
	Risk                 float64 // 0-100: 모호 표현, 민감 데이터 노출
}

// ScoreResult contains the final score and metrics
type ScoreResult struct {
	OverallScore float64 // 0-100: 종합 점수
	Metrics      Metrics
}

// MetricCalculator calculates individual metric scores
type MetricCalculator interface {
	Calculate() (float64, error)
}

// Scorer calculates the overall score from metrics
type Scorer interface {
	Score(metrics Metrics, weights Weights) ScoreResult
}

// Weights defines the weights for each metric
type Weights struct {
	Structure            float64
	Conciseness          float64
	InstructionFollowing float64
	SelfConsistency      float64
	LatencyCost          float64
	Risk                 float64
}

// DefaultWeights returns the default weight configuration
func DefaultWeights() Weights {
	return Weights{
		Structure:            0.2,
		Conciseness:          0.15,
		InstructionFollowing: 0.3,
		SelfConsistency:      0.15,
		LatencyCost:          0.1,
		Risk:                 0.1,
	}
}
