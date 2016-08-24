package mathematics

func Test_Square_1(t *testing.T) {
	if Square(2) != 4 {
		t.Error("Square function failed one test")
	}
}