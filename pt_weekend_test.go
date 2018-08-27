package main

import "testing"

func BenchmarkRenderScene(b *testing.B) {
	renderScene("benching.png")
}
