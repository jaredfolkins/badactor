package badactor

import (
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkStudioInfraction512(b *testing.B) {
	st := NewStudio(512)

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

func BenchmarkStudioInfraction1024(b *testing.B) {
	st := NewStudio(1024)

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

func BenchmarkStudioInfraction2048(b *testing.B) {
	st := NewStudio(2048)

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

func BenchmarkStudioInfraction4096(b *testing.B) {
	st := NewStudio(4096)

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

func BenchmarkStudioInfraction65536(b *testing.B) {
	st := NewStudio(65536)

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

func BenchmarkStudioInfraction262144(b *testing.B) {
	st := NewStudio(262144)

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
