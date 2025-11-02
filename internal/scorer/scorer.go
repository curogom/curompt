package scorer

// scoreCalculator implements Scorer interface
type scoreCalculator struct{}

// NewScorer creates a new score calculator
func NewScorer() Scorer {
	return &scoreCalculator{}
}

// Score calculates the overall score from metrics and weights
// Formula: score = Σ(weight[k] * metric[k])
func (s *scoreCalculator) Score(metrics Metrics, weights Weights) ScoreResult {
	overallScore :=
		metrics.Structure*weights.Structure +
			metrics.Conciseness*weights.Conciseness +
			metrics.InstructionFollowing*weights.InstructionFollowing +
			metrics.SelfConsistency*weights.SelfConsistency +
			metrics.LatencyCost*weights.LatencyCost +
			metrics.Risk*weights.Risk

	return ScoreResult{
		OverallScore: overallScore,
		Metrics:      metrics,
	}
}
