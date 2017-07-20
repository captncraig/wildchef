package constants

type Ingredient struct {
	Name     string
	Category string
	InvX     int
	InvY     int
	InvPage  int
}

const (
	Fruit    = "Fruit"
	Mushroom = "Mushroom"
	Plant    = "Plant"
	Meat     = "Meat"
	Other    = "Other"
	Dragon   = "Dragon"
	Nut      = "Nut"
	Fish     = "Fish"
	Insect   = "Insect"
	Monster  = "Monster"
	Ore      = "Ore"
)

var Ingredients = []*Ingredient{

	{Name: HeartyDurian, Category: Fruit},
	{Name: PalmFruit, Category: Fruit},
	{Name: Apple, Category: Fruit},
	{Name: Wildberry, Category: Fruit},
	{Name: Hydromelon, Category: Fruit},
	{Name: SpicyPepper, Category: Fruit},
	{Name: Voltfruit, Category: Fruit},
	{Name: FleetLotusSeeds, Category: Fruit},
	{Name: MightyBananas, Category: Fruit},

	{Name: BigHeartyTruffle, Category: Mushroom},
	{Name: HeartyTruffle, Category: Mushroom},
	{Name: EnduraShroom, Category: Mushroom},
	{Name: HylianMushroom, Category: Mushroom},
	{Name: StamellaMushroom, Category: Mushroom},
	{Name: Chillshroom, Category: Mushroom},
	{Name: Sunshroom, Category: Mushroom},
	{Name: Zapshroom, Category: Mushroom},
	{Name: Rushroom, Category: Mushroom},
	{Name: Razorshroom, Category: Mushroom},
	{Name: Ironshroom, Category: Mushroom},
	{Name: SilentShroom, Category: Mushroom},

	{Name: BigHeartyRadish, Category: Plant},
	{Name: HeartyRadish, Category: Plant},
	{Name: EnduraCarrot, Category: Plant},
	{Name: HyruleHerb, Category: Plant},
	{Name: SwiftCarrot, Category: Plant},
	{Name: FortifiedPumpkin, Category: Plant},
	{Name: CoolSafflina, Category: Plant},
	{Name: WarmSafflina, Category: Plant},
	{Name: ElectricSafflina, Category: Plant},
	{Name: SwiftViolet, Category: Plant},
	{Name: MightyThistle, Category: Plant},
	{Name: Armoranth, Category: Plant},
	{Name: BlueNightshade, Category: Plant},
	{Name: SilentPrincess, Category: Plant},

	{Name: RawGourmetMeat, Category: Meat},
	{Name: RawWholeBird, Category: Meat},
	{Name: RawPrimeMeat, Category: Meat},
	{Name: RawBirdThigh, Category: Meat},
	{Name: RawMeat, Category: Meat},
	{Name: RawBirdDrumstick, Category: Meat},

	{Name: CourserBeeHoney, Category: Other},
	{Name: HylianRice, Category: Other},
	{Name: BirdEgg, Category: Other},
	{Name: TabanthaWheat, Category: Other},
	{Name: FreshMilk, Category: Other},
	{Name: Acorn, Category: Nut},
	{Name: ChickalooTreeNut, Category: Nut},
	{Name: CaneSugar, Category: Other},
	{Name: GoatButter, Category: Other},
	{Name: GoronSpice, Category: Other},
	{Name: RockSalt, Category: Other},
	{Name: MonsterExtract, Category: Other},
	{Name: StarFragment, Category: Other},

	{Name: DinraalsScale, Category: Dragon},
	{Name: DinraalsClaw, Category: Dragon},
	{Name: ShardofDinraalsFang, Category: Dragon},
	{Name: ShardofDinraalsHorn, Category: Dragon},
	{Name: NyadrasScale, Category: Dragon},
	{Name: NyadrasClaw, Category: Dragon},
	{Name: ShardofNyadrasFang, Category: Dragon},
	{Name: ShardofNyadrasHorn, Category: Dragon},
	{Name: FaroshsScale, Category: Dragon},
	{Name: FaroshsClaw, Category: Dragon},
	{Name: ShardofFaroshsFang, Category: Dragon},
	{Name: ShardofFaroshsHorn, Category: Dragon},

	{Name: HeartySalmon, Category: Fish},
	{Name: HeartyBlueshellSnail, Category: Fish},
	{Name: HeartyBass, Category: Fish},
	{Name: HylianBass, Category: Fish},
	{Name: StaminokaBass, Category: Fish},
	{Name: ChillfinTrout, Category: Fish},
	{Name: SizzlefinTrout, Category: Fish},
	{Name: VoltfinTrout, Category: Fish},
	{Name: StealthfinTrout, Category: Fish},
	{Name: MightyCarp, Category: Fish},
	{Name: ArmoredCarp, Category: Fish},
	{Name: SankeCarp, Category: Fish},
	{Name: MightyPorgy, Category: Fish},
	{Name: ArmoredPorgy, Category: Fish},
	{Name: SneakyRiverSnail, Category: Fish},
	{Name: RazorclawCrab, Category: Fish},
	{Name: IronshellCrab, Category: Fish},
	{Name: BrightEyedCrab, Category: Fish},

	{Name: Fairy, Category: Other},

	{Name: WinterwingButterfly, Category: Insect},
	{Name: SummerwingButterfly, Category: Insect},
	{Name: ThunderwingButterfly, Category: Insect},
	{Name: SmotherwingButterfly, Category: Insect},
	{Name: ColdDarner, Category: Insect},
	{Name: WarmDarner, Category: Insect},
	{Name: ElectricDarner, Category: Insect},
	{Name: RestlessCricket, Category: Insect},
	{Name: BladedRhinoBeetle, Category: Insect},
	{Name: RuggedRhinoBeetle, Category: Insect},
	{Name: EnergeticRhinoBeetle, Category: Insect},
	{Name: SunsetFirefly, Category: Insect},
	{Name: HotFootedFrog, Category: Insect},
	{Name: TirelessFrog, Category: Insect},
	{Name: HightailLizard, Category: Insect},
	{Name: HeartyLizard, Category: Insect},
	{Name: FireproofLizard, Category: Insect},

	{Name: Flint, Category: Ore},
	{Name: Amber, Category: Ore},
	{Name: Opal, Category: Ore},
	{Name: LuminousStone, Category: Ore},
	{Name: Topaz, Category: Ore},
	{Name: Ruby, Category: Ore},
	{Name: Sapphire, Category: Ore},
	{Name: Diamond, Category: Ore},

	{Name: BokoblinHorn, Category: Monster},
	{Name: BokoblinFang, Category: Monster},
	{Name: BokoblinGuts, Category: Monster},
	{Name: MoblinHorn, Category: Monster},
	{Name: MoblinFang, Category: Monster},
	{Name: MoblinGuts, Category: Monster},
	{Name: LizalfosHorn, Category: Monster},
	{Name: LizalfosTalon, Category: Monster},
	{Name: LizalfosTail, Category: Monster},
	{Name: IcyLizalfosTail, Category: Monster},
	{Name: RedLizalfosTail, Category: Monster},
	{Name: YellowLizalfosTail, Category: Monster},
	{Name: LynelHorn, Category: Monster},
	{Name: LynelHoof, Category: Monster},
	{Name: LynelGuts, Category: Monster},
	{Name: ChuchuJelly, Category: Monster},
	{Name: WhiteChuchuJelly, Category: Monster},
	{Name: RedChuchuJelly, Category: Monster},
	{Name: YellowChuchuJelly, Category: Monster},
	{Name: KeeseWing, Category: Monster},
	{Name: IceKeeseWing, Category: Monster},
	{Name: FireKeeseWing, Category: Monster},
	{Name: ElectricKeeseWing, Category: Monster},
	{Name: KeeseEyeball, Category: Monster},
	{Name: OctorokTentacle, Category: Monster},
	{Name: OctorokEyeball, Category: Monster},
	{Name: OctoBalloon, Category: Monster},
	{Name: MoldugaFin, Category: Monster},
	{Name: MoldugaGuts, Category: Monster},
	{Name: HinoxToenail, Category: Monster},
	{Name: HinoxTooth, Category: Monster},
	{Name: HinoxGuts, Category: Monster},
	{Name: AncientScrew, Category: Monster},
	{Name: AncientSpring, Category: Monster},
	{Name: AncientGear, Category: Monster},
	{Name: AncientShaft, Category: Monster},
	{Name: AncientCore, Category: Monster},
	{Name: GiantAncientCore, Category: Monster},

	{Name: Wood, Category: Other},
}
