package crawler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"school-noti-telegram-go/internal/config"
	"school-noti-telegram-go/internal/models"

	"github.com/PuerkitoBio/goquery"
)

const (
	SchoolNoticePrefix = "school"
	DeptNoticePrefix   = "dept"
)

type Crawler struct {
	cfg *config.Config
}

func NewCrawler(cfg *config.Config) *Crawler {
	return &Crawler{
		cfg: cfg,
	}
}

func (c *Crawler) FetchAllNotices() ([]models.Notice, error) {
	var allNotices []models.Notice
	today := time.Now().In(time.FixedZone("KST", 9*60*60)).Format("2006.01.02")

	// 학교 공지사항 크롤링
	schoolNotices, err := c.fetchSchoolNotices(today)
	if err != nil {
		return nil, fmt.Errorf("학교 공지사항 크롤링 실패: %v", err)
	}
	allNotices = append(allNotices, schoolNotices...)

	// 학과 일반 공지사항 크롤링
	deptGeneralNotices, err := c.fetchDeptGeneralNotices(today)
	if err != nil {
		return nil, fmt.Errorf("학과 일반 공지사항 크롤링 실패: %v", err)
	}
	allNotices = append(allNotices, deptGeneralNotices...)

	// 학과 장학금 공지사항 크롤링
	deptScholarshipNotices, err := c.fetchDeptScholarshipNotices(today)
	if err != nil {
		return nil, fmt.Errorf("학과 장학금 공지사항 크롤링 실패: %v", err)
	}
	allNotices = append(allNotices, deptScholarshipNotices...)

	return allNotices, nil
}

func (c *Crawler) fetchSchoolNotices(today string) ([]models.Notice, error) {
	resp, err := http.Get(c.cfg.SchoolNoticeURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("잘못된 상태 코드: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("페이지 파싱 실패: %v", err)
	}

	var notices []models.Notice

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		// 날짜 확인
		dateStr := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		if dateStr != today {
			return
		}

		// 공지사항 상세 페이지 URL
		detailURL, exists := s.Find(".b-title-box a").Attr("href")
		if !exists {
			return
		}

		// 공지사항 ID 추출 및 프리픽스 추가
		rawID := strings.TrimSpace(s.Find(".b-num-box").Text())
		id := fmt.Sprintf("%s_%s", SchoolNoticePrefix, rawID)

		// 제목
		title := strings.TrimSpace(s.Find(".b-title").Text())

		postDate, _ := time.Parse("2006.01.02", dateStr)

		notice := models.Notice{
			ID:        id,
			Type:      models.SchoolNotice,
			Title:     title,
			URL:       c.cfg.SchoolNoticeURL + detailURL,
			PostDate:  postDate,
			CreatedAt: time.Now(),
		}

		// 키워드 필터링
		if c.containsKeywords(notice.Title, c.cfg.SchoolNoticeKeywords) {
			notices = append(notices, notice)
		}
	})

	return notices, nil
}

func (c *Crawler) fetchDeptGeneralNotices(today string) ([]models.Notice, error) {
	resp, err := http.Get(c.cfg.DeptGeneralURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("잘못된 상태 코드: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("페이지 파싱 실패: %v", err)
	}

	var notices []models.Notice

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		// 날짜 확인
		dateStr := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		if dateStr != today {
			return
		}

		detailURL, exists := s.Find(".b-title-box a").Attr("href")
		if !exists {
			return
		}

		rawID := strings.TrimSpace(s.Find(".b-num-box").Text())
		id := fmt.Sprintf("%s_general_%s", DeptNoticePrefix, rawID)

		title := strings.TrimSpace(s.Find(".b-title").Text())

		postDate, _ := time.Parse("2006.01.02", dateStr)

		notice := models.Notice{
			ID:        id,
			Type:      models.DeptGeneralNotice,
			Title:     title,
			URL:       c.cfg.DeptGeneralURL + detailURL,
			PostDate:  postDate,
			CreatedAt: time.Now(),
		}

		if c.containsKeywords(notice.Title, c.cfg.DeptNoticeKeywords) {
			notices = append(notices, notice)
		}
	})

	return notices, nil
}

func (c *Crawler) fetchDeptScholarshipNotices(today string) ([]models.Notice, error) {
	resp, err := http.Get(c.cfg.DeptScholarshipURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("잘못된 상태 코드: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("페이지 파싱 실패: %v", err)
	}

	var notices []models.Notice

	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		// 날짜 확인
		dateStr := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		if dateStr != today {
			return
		}

		detailURL, exists := s.Find(".b-title-box a").Attr("href")
		if !exists {
			return
		}

		rawID := strings.TrimSpace(s.Find(".b-num-box").Text())
		id := fmt.Sprintf("%s_scholarship_%s", DeptNoticePrefix, rawID)

		title := strings.TrimSpace(s.Find(".b-title").Text())

		postDate, _ := time.Parse("2006.01.02", dateStr)

		notice := models.Notice{
			ID:        id,
			Type:      models.DeptScholarshipNotice,
			Title:     title,
			URL:       c.cfg.DeptScholarshipURL + detailURL,
			PostDate:  postDate,
			CreatedAt: time.Now(),
		}

		if c.containsKeywords(notice.Title, c.cfg.DeptNoticeKeywords) {
			notices = append(notices, notice)
		}
	})

	return notices, nil
}

func (c *Crawler) containsKeywords(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
