package threads

type ReqCreateThreads struct {
	Title         string       `json:"title" validate:"required,min=1,max=100"`
	MainSection   MAINSECTION  `json:"main_section" validate:"required,validMainSection"`
	SubSection    SUBSECTION   `json:"sub_section" validate:"required,validSubSection"`
	Text          string       `json:"text" validate:"required"`
	Link          string       `json:"link" validate:"required,validLink"`
	LinkCode      string       `json:"link_code" validate:"required"`
	UnzipPassword string       `json:"unzip_password" validate:"required"`
	Price         int          `json:"price" validate:"required,number"`
	Tags          []string     `json:"tags" gorm:"serializer:json" validate:"required"`
	Status        THREADSTATUS `json:"stauts" validate:"omitempty,oneof=1 2"`
	*ReqCreateThreadsMeta
}

type ReqCreateThreadsMeta struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type Thread struct {
	*ThreadBase
	*ThreadMeta
}

type ThreadBase ReqCreateThreads

type ThreadMeta struct {
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
	PublishedAt int64 `json:"published_at"`
}

type ReqSearchByMainHome struct {
	Keywords string `json:"keywords" validate:"excluded_with=UserName"`
	UserName string `json:"user_name" validate:"excluded_with=Keywords"`
	*ReqSearchMeta
}

type ReqSearchBySection struct {
	MainSection MAINSECTION `json:"main_section" validate:"required,validMainSection"`
	SubSection  SUBSECTION  `json:"sub_section" validate:"required,validSubSection"`
}

type ReqSearchByManagement struct{}

type ReqSearchMeta struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

func NewReqSearchMeta() *ReqSearchMeta {
	return &ReqSearchMeta{
		PageSize:   3,
		PageNumber: 1,
	}
}

type ThreadsList struct {
	Total int64     `json:"total"`
	List  []*Thread `json:"list"`
}

func NewThreadsList() *ThreadsList {
	return &ThreadsList{
		List: make([]*Thread, 0, 32),
	}
}

func (*Thread) TableName() string {
	return UnitName
}
