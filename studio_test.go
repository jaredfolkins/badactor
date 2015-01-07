package badactor

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkStudioInfraction(b *testing.B) {
	st := NewStudio(65336)

	rn := "Login"
	r := &Rule{
		Name:        rn,
		Message:     "Failed to login too many times",
		StrikeLimit: 3,
		ExpireBase:  time.Second * 2,
		Sentence:    time.Second * 2,
	}

	st.AddRule(r)

	err := st.CreateDirectors(1024)
	if err != nil {
		log.Fatal(err)
	}

	st.Run()

	for i := 0; i < b.N; i++ {
		an := strconv.FormatInt(rand.Int63(), 10)
		st.Infraction(an, rn)
	}
}
