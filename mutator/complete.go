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

package mutator

import (
	"github.com/rqme/neat"
)

type Complete struct {
	Phased
	Weight
	Trait
	Activation
}

func NewComplete(ps PhasedSettings, cs ComplexifySettings, ns PruningSettings, ws WeightSettings, ts TraitSettings, as ActivationSettings) *Complete {
	return &Complete{
		Phased:     *NewPhased(ps, cs, ns),
		Weight:     Weight{WeightSettings: ws},
		Trait:      Trait{TraitSettings: ts},
		Activation: Activation{ActivationSettings: as},
	}
}

func (m *Complete) SetContext(x neat.Context) error {
	m.ctx = x
	return m.Complexify.SetContext(x)
}

// Sets the population
func (m *Complete) SetPopulation(p neat.Population) error {
	return m.Phased.SetPopulation(p)
}

func (c Complete) Mutate(g *neat.Genome) error {
	old := g.Complexity()
	if err := c.Phased.Mutate(g); err != nil {
		return err
	}
	if g.Complexity() == old {
		if err := c.Weight.Mutate(g); err != nil {
			return err
		}
		if err := c.Trait.Mutate(g); err != nil {
			return err
		}
		if err := c.Activation.Mutate(g); err != nil {
			return err
		}
	}
	return nil
}