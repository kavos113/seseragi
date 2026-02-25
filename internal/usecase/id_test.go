package usecase

type testIDGenerator struct {
	ids []string
}

func (g *testIDGenerator) GenerateID() string {
	if len(g.ids) == 0 {
		return "default-id"
	}
	id := g.ids[0]
	g.ids = g.ids[1:]
	return id
}

func newTestIDGenerator(ids ...string) *testIDGenerator {
	return &testIDGenerator{ids: ids}
}