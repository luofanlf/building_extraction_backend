package dto

import (
	"building-extraction/internal/model"
	"time"
)

type SaveProjectRequest struct {
	ProjectName string `json:"name"`
	InputImage  string `json:"input_image"`
	OutputImage string `json:"output_image"`
	ModelName   string `json:"model"`
}

type ProjectResponse struct {
	CreatedAt   time.Time `json:"created_at"`
	ProjectName string    `json:"name"`
	InputImage  string    `json:"input_image"`
	OutputImage string    `json:"output_image"`
	ModelName   string    `json:"model"`
}

func ProjectToResponse(project *model.Project) ProjectResponse {
	return ProjectResponse{
		CreatedAt:   project.CreatedAt,
		ProjectName: project.ProjectName,
		InputImage:  project.InputImage,
		OutputImage: project.OutputImage,
		ModelName:   project.ModelName,
	}
}

type ExtractionResponse struct {
	MaskUrl     string `json:"mask_url"`
	InputImage  string `json:"input_image"`
	OutputImage string `json:"output_image"`
}
