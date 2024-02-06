package token

import "testing"

func TestXxx(t *testing.T) {
	tkMaker, err := NewJWTMaker("cae1X53au6agHqAOulzCRhgDr0BG52yv")
	if err != nil {
		t.Fatal(err)
	}

	payload, err := tkMaker.VerifyToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyYzdhYjQwLTI2ZjctNGM3Ny04NTlkLWI5MzcwN2RlNzQzMCIsInR5cGUiOiJhY2Nlc3MiLCJzdWIiOiJhMjcxMjk1Ni0yY2JkLTRjNzUtYTU3Zi05ZDFkMzNhN2ZkY2MiLCJpYXQiOiIyMDIzLTEyLTI4VDEwOjU4OjE5LjA0NjQ2MDIzNSswNzowMCIsImV4cCI6IjIwMjMtMTItMjhUMjI6NTg6MTkuMDQ2NDYwMzc2KzA3OjAwIn0.MpEQTBmvEgLbR5GhqaKlpK-cWKhLiGs-kaXY_erIVzY")
	if payload == nil {
		t.Fatal(err)
	}

	t.Log(*payload)
}
