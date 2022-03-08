package message

/**
 * @Author  Flagship
 * @Date  2022/3/7 21:22
 * @Description
 */

type CommonResp struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type JobCityInfo struct {
	Code         string      `json:"code"`
	EnName       string      `json:"en_name"`
	I18nName     string      `json:"i18n_name"`
	LocationType interface{} `json:"location_type"`
	Name         string      `json:"name"`
	PyName       string      `json:"py_name"`
}

type JobCategory struct {
	Children *JobCategory `json:"children"`
	Depth    int          `json:"depth"`
	EnName   string       `json:"en_name"`
	I18nName string       `json:"i18n_name"`
	Id       string       `json:"id"`
	Name     string       `json:"name"`
}

type JobInfo struct {
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	SubTitle     string           `json:"sub_title"`
	Description  string           `json:"description"`
	Requirement  string           `json:"requirement"`
	JobCategory  *JobCategory     `json:"job_category"`
	CityInfo     *JobCityInfo     `json:"city_info"`
	RecruitType  *RecruitTypeInfo `json:"recruit_type"`
	PublishTime  int              `json:"publish_time"`
	JobHotFlag   interface{}      `json:"job_hot_flag"`
	JobSubject   *JobSubjectInfo  `json:"job_subject"`
	Code         string           `json:"code"`
	DepartmentId interface{}      `json:"department_id"`
	JobFunction  interface{}      `json:"job_function"`
	JobProcessId interface{}      `json:"job_process_id"`
	RecommendId  string           `json:"recommend_id"`
}

type JobSubjectName struct {
	EnUs interface{} `json:"en_us"`
	I18n string      `json:"i18n"`
	ZhCn string      `json:"zh_cn"`
}

type JobSubjectInfo struct {
	ActiveStatus int             `json:"active_status"`
	Id           string          `json:"id"`
	LimitCount   int             `json:"limit_count"`
	Name         *JobSubjectName `json:"name"`
}

type RecruitTypeInfo struct {
	ActiveStatus int              `json:"active_status"`
	Children     interface{}      `json:"children"`
	Depth        int              `json:"depth"`
	EnName       string           `json:"en_name"`
	I18nName     string           `json:"i18n_name"`
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Parent       *RecruitTypeInfo `json:"parent"`
}

type JobInfoData struct {
	Count       int        `json:"count"`
	Extra       string     `json:"extra"`
	JobPostList []*JobInfo `json:"job_post_list"`
}

type JobInfoResp struct {
	CommonResp
	Data *JobInfoData `json:"data"`
}

type JobInfoReq struct {
	JobCategoryIdList []interface{} `json:"job_category_id_list"`
	JobFunctionIdList []interface{} `json:"job_function_id_list"`
	Keyword           string        `json:"keyword"`
	Limit             int           `json:"limit"`
	LocationCodeList  []interface{} `json:"location_code_list"`
	Offset            int           `json:"offset"`
	PortalEntrance    int           `json:"portal_entrance"`
	PortalType        int           `json:"portal_type"`
	RecruitmentIdList []interface{} `json:"recruitment_id_list"`
	SubjectIdList     []interface{} `json:"subject_id_list"`
}
