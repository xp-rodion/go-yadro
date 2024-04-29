package database

type Result struct {
	Weight   int
	Url      string
	Keywords []string
}

type Results []Result

func (results Results) Len() int           { return len(results) }
func (results Results) Swap(i, j int)      { results[i], results[j] = results[j], results[i] }
func (results Results) Less(i, j int) bool { return results[i].Weight > results[j].Weight }

type IndexResult struct {
	Weight int
	Id     int
}

type IndexResults []IndexResult

func (results IndexResults) Len() int           { return len(results) }
func (results IndexResults) Swap(i, j int)      { results[i], results[j] = results[j], results[i] }
func (results IndexResults) Less(i, j int) bool { return results[i].Weight > results[j].Weight }

func MapToIndexResult(m map[int]int) IndexResults {
	results := make(IndexResults, len(m))
	for k, v := range m {
		results = append(results, IndexResult{Id: k, Weight: v})
	}
	return results
}
