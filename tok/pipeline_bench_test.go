package tok

import "testing"

const benchInput = `The quick brown fox jumps over the lazy dog near the riverbank.
Scientists have long debated whether the large hadron collider can produce results
that fundamentally change our understanding of particle physics. In particular,
the discovery of the Higgs boson in 2012 was a remarkable achievement for the
international team at CERN. Dr. Richard Feynman once said that the most important
thing is not to stop questioning. Despite the enormous complexity of modern
experimental apparatus, researchers continue to publish findings at an impressive rate.
However, critics argue that many results fail to replicate under controlled conditions.
The committee's report, released on Jan. 15, recommends sweeping reforms to peer review.
She won't accept the proposal unless the board revises its financial projections.
They're planning to finalize the contract by Friday, but it's unclear whether both
parties can agree on the indemnification clause. The high-performance computing
cluster at the university processed approximately 2.4 terabytes of data overnight.
Researchers said that they would release the full dataset in Q3 2026. The results
are striking: a 34% reduction in error rates and a significant improvement in throughput.
U.S. officials confirmed that talks with E.U. representatives are ongoing.
Mr. Smith, the director of operations, could not be reached for comment.`

func BenchmarkParse(b *testing.B) {
	b.SetBytes(int64(len(benchInput)))
	b.ReportAllocs()
	for b.Loop() {
		Parse(benchInput)
	}
}
