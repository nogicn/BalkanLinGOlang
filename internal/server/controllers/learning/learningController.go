package learningcontroller

import (
	"BalkanLinGO/internal/db"
	dbr "BalkanLinGO/internal/db/repository"
	"database/sql"

	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type LearnController struct {
	repo *dbr.Queries
}

func New(dbService db.Service) *LearnController {
	return &LearnController{
		repo: dbService.GetRepositoryRW(),
	}
}

func (lc *LearnController) LearnSession(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	activequestion, activeerr := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
	if activequestion == (dbr.ActiveQuestion{}) {

		err := lc.createWords(c, idInt)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
		err = lc.setActiveQuestion(&activequestion, c, idInt, 1)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
	}

	if activeerr != nil {
		if !strings.Contains(activeerr.Error(), "no rows in result") {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi!", "link": "/dashboard"})
		}
	}
	if activeword.DictionaryID != int64(idInt) {
		err := lc.repo.DeleteActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod brisanja aktivne riječi!", "link": "/dashboard"})
		}
		err = lc.createWords(c, idInt)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
		err = lc.setActiveQuestion(&activequestion, c, idInt, 1)
		if err != nil {
			return c.Status(404).Render("forOfor", fiber.Map{"status": "500", "errorText": err, "link": "/dashboard"})
		}
	}

	switch activequestion.Type {
	case 1:
		lc.LearnSessionForeignNative(c)

	case 2:
		lc.LearnSessionNativeForeign(c)

	case 3:
		lc.LearnSessionWriting(c)

	case 4:
		lc.LearnSessionPronunciation(c)

	default:
		lc.LearnSessionForeignNative(c)

	}

	return nil
}

func (lc *LearnController) LearnSessionForeignNative(c *fiber.Ctx) error {
	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	rows, err := lc.repo.GetViableWordsForUserForDictionary(c.Context(), dbr.GetViableWordsForUserForDictionaryParams{
		UserID:       int64(c.Locals("user_id").(int64)),
		DictionaryID: int64(dictidInt),
	})
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}

	finalWords := []dbr.Word{}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(rows, finalWords, 3)

	for i := 0; i < len(finalWords); i++ {
		finalWords[i].Foreignword, finalWords[i].Nativeword = finalWords[i].Nativeword, finalWords[i].Foreignword
		finalWords[i].Foreigndescription, finalWords[i].Nativedescription = finalWords[i].Nativedescription, finalWords[i].Foreigndescription
	}

	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	activeword.Foreignword, activeword.Nativeword = activeword.Nativeword, activeword.Foreignword
	activeword.Foreigndescription, activeword.Nativedescription = activeword.Nativedescription, activeword.Foreigndescription
	return c.Render("learn/selectWord", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 2, "IsAdmin": c.Locals("is_admin")})

}

func (lc *LearnController) LearnSessionNativeForeign(c *fiber.Ctx) error {
	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}

	dictid := c.Params("id")
	dictidInt, err := strconv.Atoi(dictid)
	if err != nil {
		return c.Status(404).Status(400).SendString("Invalid ID")
	}

	rows, err := lc.repo.GetViableWordsForUserForDictionary(c.Context(), dbr.GetViableWordsForUserForDictionaryParams{
		UserID:       int64(c.Locals("user_id").(int64)),
		DictionaryID: int64(dictidInt),
	})
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja riječi! nema rijeci za dict", "link": "/dashboard"})
	}

	if len(rows) < 4 {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Nema više riječi za učenje!", "link": "/dashboard"})
	}
	finalWords := []dbr.Word{}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	finalWords = append(finalWords, activeword)
	finalWords = fillWordList(rows, finalWords, 3)

	rand.Shuffle(len(finalWords), func(i, j int) { finalWords[i], finalWords[j] = finalWords[j], finalWords[i] })
	return c.Render("learn/selectWord", fiber.Map{"words": finalWords, "dictionaryId": dictidInt, "currentWord": activeword, "next": 3, "IsAdmin": c.Locals("is_admin")})

}

func (lc *LearnController) LearnSessionWriting(c *fiber.Ctx) error {

	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}
	id := c.Params("id")
	dictidInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("learn/writeWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 4, "IsAdmin": c.Locals("is_admin")})

}

func (lc *LearnController) LearnSessionPronunciation(c *fiber.Ctx) error {
	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi activequestion", "link": "/dashboard"})
		}
	}
	id := c.Params("id")
	dictidInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		return c.Status(500).Render("forOfor", fiber.Map{"status": "500", "errorText": "Greška kod dohvaćanja aktivne riječi!", "link": "/dashboard"})
	}
	return c.Render("learn/sayWord", fiber.Map{"word": activeword, "dictionaryId": dictidInt, "next": 1, "IsAdmin": c.Locals("is_admin")})

}

func (lc *LearnController) CheckAnswer(c *fiber.Ctx) error {
	answer := c.Params("answer")
	answer = strings.ToLower(answer)

	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (dbr.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	answerInt, err := strconv.Atoi(answer)
	if err != nil {
		return c.Status(500).SendString("Invalid ID")
	}

	if activequestion.Type == 1 {
		activeword.Foreignword, activeword.Nativeword = activeword.Nativeword, activeword.Foreignword
		activeword.Foreigndescription, activeword.Nativedescription = activeword.Nativedescription, activeword.Foreigndescription
	}
	var correct bool

	if int64(answerInt) == activequestion.WordID.Int64 {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 1,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 0,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		activeword, err = lc.repo.GetWordByID(c.Context(), int64(answerInt))
		if err != nil {
			if !strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
			}
		}
		correct = false
	}
	lc.setActiveQuestion(&activequestion, c, int(activeword.DictionaryID), int(activequestion.Type+1))
	return c.Render("learn/partials/word", fiber.Map{"word": activeword, "correct": correct})

}

func (lc *LearnController) CheckWritingAnswer(c *fiber.Ctx) error {
	answer := c.FormValue("foreignWord")
	answer = strings.ToLower(answer)

	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (dbr.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	var correct bool
	activeword.Foreignword = strings.ToLower(activeword.Foreignword)
	fmt.Println(answer, activeword.Foreignword)
	if answer == activeword.Foreignword {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 1,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 0,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	lc.setActiveQuestion(&activequestion, c, int(activeword.DictionaryID), int(activequestion.Type+1))
	return c.Render("learn/partials/writeWordAnswer", fiber.Map{"word": activeword, "correct": correct})

}

func (lc *LearnController) CheckListeningAnswer(c *fiber.Ctx) error {
	activequestion, err := lc.repo.GetActiveQuestionByUserID(c.Context(), sql.NullInt64{Int64: int64(c.Locals("user_id").(int64)), Valid: true})
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result") {
			return c.Status(500).SendString("Greška kod dohvaćanja aktivne riječi activequestion")
		}
	}
	if activequestion == (dbr.ActiveQuestion{}) {
		return c.Redirect("/dashboard")
	}

	activeword, err := lc.repo.GetWordByID(c.Context(), activequestion.WordID.Int64)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(500).SendString("Greška kod dohvaćanja riječi!")
		}
	}

	random := rand.Intn(100)
	var correct bool
	if random > 50 {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 1,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = true
	} else {

		err = lc.repo.SetNewDelayForUser(c.Context(), dbr.SetNewDelayForUserParams{
			IsCorrect: 0,
			UserID:    int64(c.Locals("user_id").(int64)),
			WordID:    activequestion.WordID.Int64,
		})
		if err != nil {
			return c.Status(500).SendString("Greška kod postavljanja nove riječi!")
		}
		correct = false
	}
	lc.setActiveQuestion(&activequestion, c, int(activeword.DictionaryID), int(activequestion.Type+1))
	return c.Render("learn/partials/sayWordAnswer", fiber.Map{"word": activeword, "correct": correct})

}
