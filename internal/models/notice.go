package models

import "time"

type NoticeType string

const (
	SchoolNotice          NoticeType = "SCHOOL"
	DeptGeneralNotice     NoticeType = "DEPT_GENERAL"
	DeptScholarshipNotice NoticeType = "DEPT_SCHOLARSHIP"
)

type Notice struct {
	ID        string     `json:"id"`
	Type      NoticeType `json:"type"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Category  string     `json:"category"`
	PostDate  time.Time  `json:"post_date"`
	CreatedAt time.Time  `json:"created_at"`
}
