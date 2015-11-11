/*
Copyright (c) 2015, Brian Hummer (brian@redq.me)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"flag"
	"log"
	"fmt"
	"time"
	"math/rand"

	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
	"github.com/rqme/neat/x/starter"
	"github.com/rqme/neat/x/trials"
)

type Evaluator struct {
	show     bool
	useTrial bool
	trialNum int
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

// Evaluate computes the error for the XOR problem with the phenome
//
// To compute fitness, the distance of the output from the correct answer was summed for all four
// input patterns. The result of this error was subtracted from 4 so that higher fitness would mean
// better networks. The resulting number was squared to give proportionally more fitness the closer
// a network was to a solution. (Stanley, 43)
func (e Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {
	// Run experiment
	var err error
	stop := false

	f := Flappy{}
	f.Alive = true
	f.screen.x = 10
	f.screen.y = 10
	f.bird.velocity = 1
	f.bird.posX = 4
	f.bird.posY = 2
	f.obs.Y = 12
	f.obs.bottomX = 6
	f.obs.topX = 4

	in := make([]float64, f.screen.x * f.screen.y)

	for f.Alive && f.Fitness < 500 {
		in = f.Export()
		outputs, err := p.Activate(in)
		if err != nil {
			break
		}
		f.Next(outputs[0])
		var asd bool
		if outputs[0] >= 0.5 {
			asd = true
		} else {
			asd = false
		}
		if(e.show) {
			fmt.Println(f.obs.Y, " ", asd, "", f.bird.posX)
		}
	}

	// Calculate the result
	if f.Fitness > 499 {
		stop = true
	}
	r = result.New(p.ID(), f.Fitness, err, stop)
	return
}

func unhex(c byte) byte {
    switch {
    case '0' <= c && c <= '9':
        return c - '0'
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10
    }
    return 0
}

func (e *Evaluator) ShowWork(s bool) {
	e.show = s
}

func main() {
	flag.Parse()
	//defer profile.Start(profile.CPUProfile).Stop()
	if err := trials.Run(func(i int) (*neat.Experiment, error) {
		ctx := starter.NewContext(&Evaluator{})
		if exp, err := starter.NewExperiment(ctx, ctx, i); err != nil {
			return nil, err
		} else {
			return exp, nil
		}

	}); err != nil {
		log.Fatal("Could not run XOR: ", err)
	}

}

type Flappy struct {
	screen struct {
		x, y int
	}
	Alive bool
	Fitness float64
	bird struct {
		posX, posY int
		velocity int
	}
	obs struct {
		Y int
		topX int
		bottomX int
	}
}

func (f *Flappy) Export() (out []float64) {
	out = make([]float64, f.screen.x * f.screen.y)
	for x := 0; x < f.screen.x; x++ {
		for y := 0; y < f.screen.y; y++ {
			pos := x * f.screen.x + y
			if f.bird.posX == x && f.bird.posY == y {
				out[pos] = 1.0
			} else if f.obs.Y == y && (x <= f.obs.topX || x >= f.obs.bottomX) {
				out[pos] = 0.5
			}
		}
	}

	return
}

func (f *Flappy) Next(in float64) {
	if in >= 0.5 {
		f.bird.velocity = 3
	}

	f.BirdNext()
	f.ObstacleNext()
	f.CheckAlive()

	f.Fitness++
}

func (f *Flappy) BirdNext() {
	if f.bird.velocity > 0 {
		if f.bird.posX > 1 {
			f.bird.posX--
		}
		f.bird.velocity--
	} else {
		f.bird.posX++
	}
}

func (f *Flappy) ObstacleNext() {
	if f.obs.Y < 0 {
		f.obs.Y = f.screen.y
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		f.obs.topX = r.Intn(7)
		f.obs.bottomX = f.obs.bottomX + 2
	} else {
		f.obs.Y--
	}
}

func (f *Flappy) CheckAlive() {
	if f.bird.posX > f.screen.x {
		f.Alive = false
	}

	if f.bird.posY == f.obs.Y {
		if f.bird.posX > f.obs.bottomX || f.bird.posX < f.obs.topX {
			f.Alive = false
		}
	}
}
