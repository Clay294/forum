package threads

type MAINSECTION string

type SUBSECTION string

type THREADSTATUS byte

const (
	CreateThreadUrl        = "/thread"
	CreateThreadMethod     = "post"
	SearchByMainHomeUrl    = "/thread_mainhome"
	SearchByMainHomeMethod = "get"
)

const (
	UnitName = "threads"
)

var sectionTable = map[MAINSECTION]map[SUBSECTION]struct{}{
	Game: {
		Western:  struct{}{},
		Japanese: struct{}{},
		Domestic: struct{}{},
	},
	Video: {
		Western:  struct{}{},
		Japanese: struct{}{},
		Domestic: struct{}{},
		Movie:    struct{}{},
	},
	AC: {
		Anime: struct{}{},
		Comic: struct{}{},
	},
}

const (
	Game     MAINSECTION = "game"
	Video    MAINSECTION = "video"
	AC       MAINSECTION = "ac"
	Western  SUBSECTION  = "western"
	Japanese SUBSECTION  = "japanese"
	Movie    SUBSECTION  = "movie"
	Domestic SUBSECTION  = "domestic"
	Anime    SUBSECTION  = "anime"
	Comic    SUBSECTION  = "comic"
)

const (
	StatusDraft THREADSTATUS = iota + 1
	StatusPublished
)

const (
	DefaultPageSize   = 3
	DefaultPageNumber = 1
)
