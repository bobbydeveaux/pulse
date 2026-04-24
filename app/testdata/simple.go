package testdata

// SimpleFunc has low complexity.
func SimpleFunc(x int) int {
	return x + 1
}

// MediumFunc has moderate complexity.
func MediumFunc(x, y int) string {
	if x > 0 {
		if y > 0 {
			return "both positive"
		}
		return "x positive"
	} else if x == 0 {
		return "x zero"
	}

	for i := 0; i < y; i++ {
		if i%2 == 0 {
			continue
		}
	}

	return "negative"
}

// ComplexFunc has high complexity.
func ComplexFunc(items []string, flag bool) int {
	count := 0
	for _, item := range items {
		if item == "" {
			continue
		}

		switch {
		case item == "a" && flag:
			count++
		case item == "b" || !flag:
			count += 2
		case item == "c":
			if flag {
				count += 3
			} else {
				for i := 0; i < count; i++ {
					if i > 10 {
						break
					}
				}
			}
		default:
			count--
		}
	}

	if count > 100 || (count < 0 && flag) {
		return 0
	}

	return count
}
