package main

var (
	alignment_patterns_table = [][][]int{
		{},                                   // version 1
		{{6, 18}, {18, 6}, {18, 18}, {6, 6}}, // version 2
		{{6, 22}, {22, 6}, {22, 22}, {6, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {6, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {6, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {6, 6}},
		{{6, 22}, {22, 6}, {22, 22}, {22, 38}, {38, 22}, {38, 38}, {6, 6}, {6, 38}, {38, 6}}, // version 7
		{{6, 24}, {24, 6}, {24, 24}, {24, 42}, {42, 24}, {42, 42}, {6, 6}, {6, 42}, {42, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 46}, {46, 26}, {46, 46}, {6, 6}, {6, 46}, {46, 6}},
		{{6, 28}, {28, 6}, {28, 28}, {28, 50}, {50, 28}, {50, 50}, {6, 6}, {6, 50}, {50, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}, {6, 6}, {6, 54}, {54, 6}},
		{{6, 32}, {32, 6}, {32, 32}, {32, 58}, {58, 32}, {58, 58}, {6, 6}, {6, 58}, {58, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {34, 62}, {62, 34}, {62, 62}, {6, 6}, {6, 62}, {62, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 46}, {46, 26}, {46, 46}, {26, 66}, {66, 26}, {66, 66}, {6, 6}, {6, 46}, {46, 6}, {6, 66}, {66, 6}}, // version 14
		{{6, 26}, {26, 6}, {26, 26}, {26, 48}, {48, 26}, {48, 48}, {26, 70}, {70, 26}, {70, 70}, {6, 6}, {6, 48}, {48, 6}, {6, 70}, {70, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 50}, {50, 26}, {50, 50}, {26, 74}, {74, 26}, {74, 74}, {6, 6}, {6, 50}, {50, 6}, {6, 74}, {74, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}, {30, 78}, {78, 30}, {78, 78}, {6, 6}, {6, 54}, {54, 6}, {6, 78}, {78, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 56}, {56, 30}, {56, 56}, {30, 82}, {82, 30}, {82, 82}, {6, 6}, {6, 56}, {56, 6}, {6, 82}, {82, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 58}, {58, 30}, {58, 58}, {30, 86}, {86, 30}, {86, 86}, {6, 6}, {6, 58}, {58, 6}, {6, 86}, {86, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {34, 62}, {62, 34}, {62, 62}, {34, 90}, {90, 34}, {90, 90}, {6, 6}, {6, 62}, {62, 6}, {6, 90}, {90, 6}},
		{{6, 28}, {28, 6}, {28, 28}, {28, 50}, {50, 28}, {50, 50}, {28, 72}, {72, 28}, {72, 72}, {28, 94}, {94, 28}, {94, 94}, {6, 6}, {6, 50}, {50, 6}, {6, 72}, {72, 6}, {6, 94}, {94, 6}}, // version 21
		{{6, 26}, {26, 6}, {26, 26}, {26, 50}, {50, 26}, {50, 50}, {26, 74}, {74, 26}, {74, 74}, {26, 98}, {98, 26}, {98, 98}, {6, 6}, {6, 50}, {50, 6}, {6, 74}, {74, 6}, {6, 98}, {98, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}, {30, 78}, {78, 30}, {78, 78}, {30, 102}, {102, 30}, {102, 102}, {6, 6}, {6, 54}, {54, 6}, {6, 78}, {78, 6}, {6, 102}, {102, 6}},
		{{6, 28}, {28, 6}, {28, 28}, {28, 54}, {54, 28}, {54, 54}, {28, 80}, {80, 28}, {80, 80}, {28, 106}, {106, 28}, {106, 106}, {6, 6}, {6, 54}, {54, 6}, {6, 80}, {80, 6}, {6, 106}, {106, 6}},
		{{6, 32}, {32, 6}, {32, 32}, {32, 58}, {58, 32}, {58, 58}, {32, 84}, {84, 32}, {84, 84}, {32, 110}, {110, 32}, {110, 110}, {6, 6}, {6, 58}, {58, 6}, {6, 84}, {84, 6}, {6, 110}, {110, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 58}, {58, 30}, {58, 58}, {30, 86}, {86, 30}, {86, 86}, {30, 114}, {114, 30}, {114, 114}, {6, 6}, {6, 58}, {58, 6}, {6, 86}, {86, 6}, {6, 114}, {114, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {34, 62}, {62, 34}, {62, 62}, {34, 90}, {90, 34}, {90, 90}, {34, 118}, {118, 34}, {118, 118}, {6, 6}, {6, 62}, {62, 6}, {6, 90}, {90, 6}, {6, 118}, {118, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 50}, {50, 26}, {50, 50}, {26, 74}, {74, 26}, {74, 74}, {26, 98}, {98, 26}, {98, 98}, {26, 122}, {122, 26}, {122, 122}, {6, 6}, {6, 50}, {50, 6}, {6, 74}, {74, 6}, {6, 98}, {98, 6}, {6, 122}, {122, 6}}, // version 28
		{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}, {30, 78}, {78, 30}, {78, 78}, {30, 102}, {102, 30}, {102, 102}, {30, 126}, {126, 30}, {126, 126}, {6, 6}, {6, 54}, {54, 6}, {6, 78}, {78, 6}, {6, 102}, {102, 6}, {6, 126}, {126, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 52}, {52, 26}, {52, 52}, {26, 78}, {78, 26}, {78, 78}, {26, 104}, {104, 26}, {104, 104}, {26, 130}, {130, 26}, {130, 130}, {6, 6}, {6, 52}, {52, 6}, {6, 78}, {78, 6}, {6, 104}, {104, 6}, {6, 130}, {130, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 56}, {56, 30}, {56, 56}, {30, 82}, {82, 30}, {82, 82}, {30, 108}, {108, 30}, {108, 108}, {30, 134}, {134, 30}, {134, 134}, {6, 6}, {6, 56}, {56, 6}, {6, 82}, {82, 6}, {6, 108}, {108, 6}, {6, 134}, {134, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {34, 60}, {60, 34}, {60, 60}, {34, 86}, {86, 34}, {86, 86}, {34, 112}, {112, 34}, {112, 112}, {34, 138}, {138, 34}, {138, 138}, {6, 6}, {6, 60}, {60, 6}, {6, 86}, {86, 6}, {6, 112}, {112, 6}, {6, 138}, {138, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 58}, {58, 30}, {58, 58}, {30, 86}, {86, 30}, {86, 86}, {30, 114}, {114, 30}, {114, 114}, {30, 142}, {142, 30}, {142, 142}, {6, 6}, {6, 58}, {58, 6}, {6, 86}, {86, 6}, {6, 114}, {114, 6}, {6, 142}, {142, 6}},
		{{6, 34}, {34, 6}, {34, 34}, {34, 62}, {62, 34}, {62, 62}, {34, 90}, {90, 34}, {90, 90}, {34, 118}, {118, 34}, {118, 118}, {34, 146}, {146, 34}, {146, 146}, {6, 6}, {6, 62}, {62, 6}, {6, 90}, {90, 6}, {6, 118}, {118, 6}, {6, 146}, {146, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 54}, {54, 30}, {54, 54}, {30, 78}, {78, 30}, {78, 78}, {30, 102}, {102, 30}, {102, 102}, {30, 126}, {126, 30}, {126, 126}, {30, 150}, {150, 30}, {150, 150}, {6, 6}, {6, 54}, {54, 6}, {6, 78}, {78, 6}, {6, 102}, {102, 6}, {6, 126}, {126, 6}, {6, 150}, {150, 6}}, // version 35
		{{6, 24}, {24, 6}, {24, 24}, {24, 50}, {50, 24}, {50, 50}, {24, 76}, {76, 24}, {76, 76}, {24, 102}, {102, 24}, {102, 102}, {24, 128}, {128, 24}, {128, 128}, {24, 154}, {154, 24}, {154, 154}, {6, 6}, {6, 50}, {50, 6}, {6, 76}, {76, 6}, {6, 102}, {102, 6}, {6, 128}, {128, 6}, {6, 154}, {154, 6}},
		{{6, 28}, {28, 6}, {28, 28}, {28, 54}, {54, 28}, {54, 54}, {28, 80}, {80, 28}, {80, 80}, {28, 106}, {106, 28}, {106, 106}, {28, 132}, {132, 28}, {132, 132}, {28, 158}, {158, 28}, {158, 158}, {6, 6}, {6, 54}, {54, 6}, {6, 80}, {80, 6}, {6, 106}, {106, 6}, {6, 132}, {132, 6}, {6, 158}, {158, 6}},
		{{6, 32}, {32, 6}, {32, 32}, {32, 58}, {58, 32}, {58, 58}, {32, 84}, {84, 32}, {84, 84}, {32, 110}, {110, 32}, {110, 110}, {32, 136}, {136, 32}, {136, 136}, {32, 162}, {162, 32}, {162, 162}, {6, 6}, {6, 58}, {58, 6}, {6, 84}, {84, 6}, {6, 110}, {110, 6}, {6, 136}, {136, 6}, {6, 162}, {162, 6}},
		{{6, 26}, {26, 6}, {26, 26}, {26, 54}, {54, 26}, {54, 54}, {26, 82}, {82, 26}, {82, 82}, {26, 110}, {110, 26}, {110, 110}, {26, 138}, {138, 26}, {138, 138}, {26, 166}, {166, 26}, {166, 166}, {6, 6}, {6, 54}, {54, 6}, {6, 82}, {82, 6}, {6, 110}, {110, 6}, {6, 138}, {138, 6}, {6, 166}, {166, 6}},
		{{6, 30}, {30, 6}, {30, 30}, {30, 58}, {58, 30}, {58, 58}, {30, 86}, {86, 30}, {86, 86}, {30, 114}, {114, 30}, {114, 114}, {30, 142}, {142, 30}, {142, 142}, {30, 170}, {170, 30}, {170, 170}, {6, 6}, {6, 58}, {58, 6}, {6, 86}, {86, 6}, {6, 114}, {114, 6}, {6, 142}, {142, 6}, {6, 170}, {170, 6}},
	}
)

func get_alignment_patterns_for_version(version int) [][]int {
	if version < 2 {
		return [][]int{}
	}

	return alignment_patterns_table[version-1]
}
