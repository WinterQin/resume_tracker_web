package model

import "time"

type InterviewSchedule struct {
	ApplicationID      string    `json:"application_id"`
	WrittenTestDate    time.Time `json:"written_test_date"`
	AssessmentDate     time.Time `json:"assessment_date"`
	AIInterviewDate    time.Time `json:"ai_interview_date"`
	HRInterviewDate    time.Time `json:"hr_interview_date"`
	TechInterviewDate  time.Time `json:"tech_interview_date"`
	FinalInterviewDate time.Time `json:"final_interview_date"`
}
