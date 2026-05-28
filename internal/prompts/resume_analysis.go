package prompts

import "fmt"

const systemPrompt = `You are an expert technical recruiter. Analyze the provided resume against the job description and return a structured evaluation.`

func BuildPrompt(resume, jobDescription string) string {
	return fmt.Sprintf("Resume:\n%s\n\nJob Description:\n%s", resume, jobDescription)
}

func SystemPrompt() string {
	return systemPrompt
}
