package mock

const (
	// MaxQuestions is the maximum number of questions in a session
	MaxQuestions = 10
	// MaxDurationMinutes is the maximum session duration in minutes
	MaxDurationMinutes = 15
	// ShadowPromptSurrender is the prompt injected when user surrenders
	ShadowPromptSurrender = "User surrenders. Give a snappy 1-2 sentence correction and move on."
	// InitialPromptTemplate is the template for the initial interview prompt
	InitialPromptTemplate = "Here is the user's resume:\n\n%s\n\nConduct a technical interview. Ask one question at a time. When you want to move to the next question, include the signal <NEXT> at the end of your response. When you've finished all questions, include the signal <ROAST> at the end of your final response."
)

// Grade thresholds based on surrender count
const (
	GradeASurrenders = 0
	GradeBSurrenders = 2
	GradeCSurrenders = 4
	GradeDSurrenders = 6
	// GradeF: 7+ surrenders
)

// PersonaLabels maps grades to persona labels
var PersonaLabels = map[string]string{
	"A": "ARCHITECT MATERIAL",
	"B": "SOLID CANDIDATE",
	"C": "NEEDS WORK",
	"D": "JUNIOR AT BEST",
	"F": "TERMINATED",
}

