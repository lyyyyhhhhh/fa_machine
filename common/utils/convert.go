package utils

func ToInterfaceSlice[T any, I interface{}](in []T) []I {
	out := make([]I, 0, len(in))
	for _, v := range in {
		out = append(out, any(v).(I))
	}
	return out
}
