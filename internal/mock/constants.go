package mock

import "strings"

const (
	// MaxQuestions is the maximum number of questions in a session
	MaxQuestions = 10
	// MaxDurationMinutes is the maximum session duration in minutes
	MaxDurationMinutes = 15
	// ShadowPromptSurrender is the prompt injected when user surrenders
	ShadowPromptSurrender = "User surrenders. Give a snappy 1-2 sentence correction and move on."
	// InitialPromptTemplate is the template for the initial interview prompt
	// Placeholders: %s = resume content, %s = topic instructions
	InitialPromptTemplate = "Here is the user's resume:\n\n%s\n\n%s\n\nConduct a technical interview following these guidelines:\n- Ask questions that a real interviewer would ask for this role/experience level\n- Vary question types (conceptual, practical, problem-solving)\n- Avoid repeating similar questions. Reference conversation history to ensure variety\n- Ask follow-up questions based on the candidate's answers, not generic questions\n- Ask one question at a time\n- When you want to move to the next question, include the signal <NEXT> at the end of your response\n- When you've finished all questions, include the signal <ROAST> at the end of your final response."
)

// Grade thresholds based on surrender count
const (
	GradeASurrenders = 0
	GradeBSurrenders = 2
	GradeCSurrenders = 4
	GradeDSurrenders = 6
	// GradeF: 7+ surrenders
)

// ValidInterviewTopics contains the list of valid interview topics
var ValidInterviewTopics = []string{
	"Go",
	"System Design",
	"Algorithms",
	"Data Structures",
	"Databases",
	"Networking",
	"Concurrency",
	"Testing",
}

// PersonaLabels maps grades to persona labels
var PersonaLabels = map[string]string{
	"A": "ARCHITECT MATERIAL",
	"B": "SOLID CANDIDATE",
	"C": "NEEDS WORK",
	"D": "JUNIOR AT BEST",
	"F": "TERMINATED",
}

// BuildTopicInstructions builds topic instruction text for prompts
func BuildTopicInstructions(selectedTopics, excludedTopics []string) string {
	var parts []string

	if len(selectedTopics) > 0 {
		parts = append(parts, "Focus on these topics: "+strings.Join(selectedTopics, ", "))
	}

	if len(excludedTopics) > 0 {
		parts = append(parts, "Do not ask about: "+strings.Join(excludedTopics, ", "))
	}

	if len(parts) == 0 {
		return "You may ask about any technical topics relevant to the role."
	}

	return strings.Join(parts, "\n")
}

