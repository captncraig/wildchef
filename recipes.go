package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"

	. "github.com/captncraig/wildchef/constants"
)

type Ingredient struct {
	Name     string
	InvX     int
	InvY     int
	InvPage  int
	Category string
}

func (r Result) String() string {
	lines := []string{}
	eff := r.EffectString
	if eff == "" && r.Effect != effNone {
		eff = fmt.Sprintf("? (%d %x)", r.Effect, r.Effect)
	}
	if eff != "" {
		lines = append(lines, fmt.Sprintf("  %s Strength:%s (%x)", eff, r.StrengthString, r.Strength))
	}
	if r.Duration != 0 {
		lines = append(lines, fmt.Sprintf("  Duration: %s", time.Second*time.Duration(r.Duration)))
	}
	if len(lines) > 0 {
		lines = append(lines, "")
	}
	return fmt.Sprintf(`~~~~~~~~~~~~~~~~~~~~
Name: %s %s (%s)
  %d rupees   %.2g hearts
  %s
%s~~~~~~~~~~~~~~~~~~~~`, eff, r.Name, r.ID, r.Cost, float64(r.Hearts)/4, r.Ingredients, strings.Join(lines, "\n"))
}

type Recipe []*Ingredient

var allRecipes = make(chan Recipe)
var allIngs = map[string]*Ingredient{}

func loadRecipes() {
	var ings = []*Ingredient{}
	dat, err := ioutil.ReadFile("ingredients.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(dat, &ings); err != nil {
		log.Fatal(err)
	}
	//set x,y,page
	x, y, page := 0, 0, 0
	for _, ing := range ings {
		ing.InvX = x
		ing.InvY = y
		ing.InvPage = page
		x++
		if x == 5 {
			y++
			x = 0
			if y == 4 {
				page++
				y = 0
			}
		}
		allIngs[ing.Name] = ing
	}
	//All ingredients. Once.
	for _, i := range ings {
		r(i)
	}

	// for 2 ingredients, this tells if it is interesting
	shouldMake := func(a, b *Ingredient) bool {
		cats := []string{a.Category, b.Category}
		sort.Strings(cats)
		catsJ := strings.Join(cats, ",")
		if strings.Contains(a.Name, "Dinraal") || strings.Contains(b.Name, "Dinraal") {
			return false
		}
		if strings.Contains(a.Name, "Nyadra") || strings.Contains(b.Name, "Nyadra") {
			return false
		}
		if catsJ == "Insect,Monster" {
			return true
		}
		//monsters are dubious
		if strings.Contains(catsJ, "Monster") {
			// let's use bokoblins no matter what
			if a.Name == "Bokoblin Guts" || b.Name == "Bokoblin Guts" {
				return true
			}
			return false
		}
		// any insects left will give dubious food
		if strings.Contains(catsJ, "Insect") {
			// let's use Cold Darner no matter what
			if a.Name == "Cold Darner" || b.Name == "Cold Darner" {
				return true
			}
			return false
		}
		return true
	}

	// // everything together 1:1
	for i := 0; i < len(ings); i++ {
		for j := i + i; j < len(ings); j++ {
			a := ings[i]
			b := ings[j]
			if shouldMake(a, b) {
				r(a, b)
			}
		}
	}

	// //2,3,4,5 of each
	for _, i := range ings {
		r(i, i)
		r(i, i, i)
		r(i, i, i, i)
		r(i, i, i, i, i)
	}

	// all hearty
	hearty := []string{
		HeartyDurian,
		BigHeartyTruffle,
		HeartyTruffle,
		BigHeartyRadish,
		HeartyRadish,
		HeartySalmon,
		HeartyBlueshellSnail,
		HeartyBass,
		HeartyLizard,
		LynelGuts,
	}

	// energizing
	energizing := []string{
		StamellaMushroom,
		RestlessCricket,
		CourserBeeHoney,
		BrightEyedCrab,
		StaminokaBass,
		EnergeticRhinoBeetle,
		LynelGuts,
	}

	//enduring
	enduring := []string{
		EnduraShroom,
		TirelessFrog,
		EnduraCarrot,
		LynelGuts,
	}

	//Spicy
	spicy := []string{
		SpicyPepper,
		WarmSafflina,
		SummerwingButterfly,
		Sunshroom,
		WarmDarner,
		SizzlefinTrout,
		LynelGuts,
	}

	//chilly
	chilly := []string{
		Hydromelon,
		CoolSafflina,
		WinterwingButterfly,
		Chillshroom,
		ColdDarner,
		ChillfinTrout,
		LynelGuts,
	}

	//electro
	electro := []string{
		Voltfruit,
		ElectricSafflina,
		ThunderwingButterfly,
		Zapshroom,
		ElectricDarner,
		VoltfinTrout,
		LynelGuts,
	}

	//fireproof
	fireproof := []string{
		FireproofLizard,
		SmotherwingButterfly,
		LynelGuts,
	}

	//hasty
	hasty := []string{
		Rushroom,
		SwiftCarrot,
		HightailLizard,
		FleetLotusSeeds,
		SwiftViolet,
		HotFootedFrog,
		LynelGuts,
	}

	//sneaky
	sneaky := []string{
		BlueNightshade,
		SneakyRiverSnail,
		SunsetFirefly,
		SilentShroom,
		StealthfinTrout,
		SilentPrincess,
		LynelGuts,
	}

	//mighty
	mighty := []string{
		MightyThistle,
		BladedRhinoBeetle,
		MightyBananas,
		Razorshroom,
		MightyCarp,
		RazorclawCrab,
		MightyPorgy,
		LynelGuts,
	}

	//tough
	tough := []string{
		Armoranth,
		RuggedRhinoBeetle,
		FortifiedPumpkin,
		Ironshroom,
		ArmoredCarp,
		IronshellCrab,
		ArmoredPorgy,
		LynelGuts,
	}

	// a bunch of interesting things, including most ingredients
	interesting := []string{
		Apple,
		Voltfruit,
		HylianMushroom,
		Zapshroom,
		HyruleHerb,
		ElectricSafflina,
		RawPrimeMeat,
		CourserBeeHoney,
		HylianRice,
		BirdEgg,
		TabanthaWheat,
		FreshMilk,
		Acorn,
		CaneSugar,
		GoatButter,
		GoronSpice,
		RockSalt,
		MonsterExtract,
		StarFragment,
		BokoblinHorn,
		Amber,
		Fairy,
		ElectricDarner,
	}

	allCombos(hearty...)
	allCombos(energizing...)
	allCombos(enduring...)
	allCombos(energizing...)
	allCombos(spicy...)
	allCombos(chilly...)
	allCombos(electro...)
	allCombos(fireproof...)
	allCombos(hasty...)
	allCombos(sneaky...)
	allCombos(mighty...)
	allCombos(tough...)
	allCombos(interesting...)

	// and finish out with a brute force search.
	allNames := []string{}
	for _, i := range ings {
		allNames = append(allNames, i.Name)
	}
	allCombos(allNames...)
}

func allCombos(i ...string) {
	for _, s := range i {
		if allIngs[s] == nil {
			panic("YOU DON'T HAVE " + s)
		}
	}
	for a := 0; a < len(i); a++ {
		ia := allIngs[i[a]]
		r(ia)
		for b := a; b < len(i); b++ {
			ib := allIngs[i[b]]
			r(ia, ib)
			for c := b; c < len(i); c++ {
				ic := allIngs[i[c]]
				r(ia, ib, ic)
				for d := c; d < len(i); d++ {
					id := allIngs[i[d]]
					r(ia, ib, ic, id)
					for e := d; e < len(i); e++ {
						ie := allIngs[i[e]]
						r(ia, ib, ic, id, ie)
					}
				}
			}
		}
	}
}

func r(i ...*Ingredient) Recipe {
	r := Recipe(i)
	allRecipes <- r
	return r
}

func genAhkClear(n int) []byte {
	out := `#NoEnv  ; Recommended for performance and compatibility with future AutoHotkey releases.
; #Warn  ; Enable warnings to assist with detecting common errors.
SendMode Input  ; Recommended for new scripts due to its superior speed and reliability.
SetWorkingDir %A_ScriptDir%  ; Ensures a consistent starting directory.

`
	out += ahkKey("escape")
	for i := 0; i < 8; i++ {
		out += ahkKey("l")
	}
	for i := 0; i < n; i++ {
		out += ahkKey("right")
		out += ahkKey("w")
		out += ahkKey("right")
		out += ahkKey("w")
		out += ahkKey("right")
		out += ahkKey("up")
	}
	out += ahkKey("escape")
	return []byte(out)
}

func runAhk(d []byte) {
	e(ioutil.WriteFile("cook.ahk", d, 0777))
	cmd := exec.CommandContext(ctx, "AutoHotkey.exe", "cook.ahk")
	s := time.Now()
	err := cmd.Run()
	e(err)
	fmt.Println(time.Now().Sub(s))
}

func (r Recipe) genAhk() []byte {
	out := `#NoEnv  ; Recommended for performance and compatibility with future AutoHotkey releases.
; #Warn  ; Enable warnings to assist with detecting common errors.
SendMode Input  ; Recommended for new scripts due to its superior speed and reliability.
SetWorkingDir %A_ScriptDir%  ; Ensures a consistent starting directory.

`
	page := 0
	x := 0
	y := 0
	out += ahkKey("escape")

	out += ahkKey("left")
	for _, i := range r {
		out += ahkcomment(i.Name)
		//find it
		for page < i.InvPage {
			out += ahkKey("l")
			page++
		}
		for page > i.InvPage {
			out += ahkKey("j")
			page--
		}
		for x < i.InvX {
			out += ahkKey("d")
			x++
		}
		for x > i.InvX {
			out += ahkKey("a")
			x--
		}
		for y < i.InvY {
			out += ahkKey("s")
			y++
		}
		for y > i.InvY {
			out += ahkKey("w")
			y--
		}
		out += ahkcomment("add")
		out += ahkKey("right")
	}
	out += ahkcomment("exit inv")
	out += ahkKey("down")
	out += "Sleep 2000\n"
	out += ahkcomment("cook!")
	out += ahkKey("right")
	out += "Sleep 2000\n"
	out += ahkKey("left")
	out += "Sleep 4500\n"
	out += ahkKey("right")
	return []byte(out)
}

func ahkKey(s string) string {
	return fmt.Sprintf(`SendInput, {%s down}
Sleep 100
SendInput, {%s up}
Sleep 300
`, s, s)
}

func ahkcomment(s string) string {
	return "; " + s + "\n"
}
