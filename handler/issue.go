package handler

import (
	"github.com/gofiber/fiber/v2"
	"internal-backend/database"
	"log"
)

type Issue struct {
}

// GetAll get all issues from db
func (i Issue) GetAll(ctx *fiber.Ctx) error {
	var (
		err   error
		mode  string
		issue []database.Issue
	)

	// Try to parse detail mode, default brief if null or malformed
	mode = ctx.Query("mode", "brief")

	// Retrieve all issues
	err = database.Issue{}.GetAll(mode, &issue)
	if err != nil {
		log.Printf(ErrStringMsg("issue/GetAll while retrieving aal issues", err))
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all issues",
		"retrieved": len(issue),
		"mode":      mode,
		"data":      issue,
	})
}

// PostIssue add new issue in db and return new record
func (i Issue) PostIssue(ctx *fiber.Ctx) error {
	var (
		err   error
		issue database.Issue
	)

	err = ctx.BodyParser(&issue)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = issue.PostIssue()
	if err != nil {
		log.Printf(ErrStringMsg("issue/PostIssue while adding record to db", err))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "issue added to db",
		"data":    issue,
	})
}

// CloseIssue close selected issue
func (i Issue) CloseIssue(ctx *fiber.Ctx) error {
	var (
		err   error
		issue database.Issue
	)

	if ctx.Params("id") == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	issue.Id = ctx.Params("id")
	err = issue.CloseIssue()
	if err != nil {
		log.Printf(ErrStringMsg("issue/CloseIssue while adding record to db", err))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Issue closed",
		"data":    issue.Id,
	})
}

// PostDetail add new detail do selected issue and return new record
func (i Issue) PostDetail(ctx *fiber.Ctx) error {
	var (
		err         error
		issueDetail database.IssueDetail
	)

	err = ctx.BodyParser(&issueDetail)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = issueDetail.PostIssueDetail(issueDetail.IssueID)
	if err != nil {
		log.Printf(ErrStringMsg("issue/PostDetail while adding record to db", err))
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "issue detail added to db",
		"data":    issueDetail,
	})
}
